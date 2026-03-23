package app

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/Ferari430/tg_sender/internal/adapters/in"
	"github.com/Ferari430/tg_sender/internal/adapters/out"
	telegramBot "github.com/Ferari430/tg_sender/internal/bot"
	"github.com/Ferari430/tg_sender/internal/config"
	sender "github.com/Ferari430/tg_sender/internal/handler/cron"
	dochandler "github.com/Ferari430/tg_sender/internal/handler/fileHandler"
	userhandler "github.com/Ferari430/tg_sender/internal/handler/userHandler"
	"github.com/Ferari430/tg_sender/internal/infra/inMemory"
	"github.com/Ferari430/tg_sender/internal/infra/kafka"
	"github.com/Ferari430/tg_sender/internal/service/file/download"
	"github.com/Ferari430/tg_sender/internal/service/file/saveConverted"
	"github.com/Ferari430/tg_sender/internal/service/file/send"
	userservice "github.com/Ferari430/tg_sender/internal/service/userService"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Consumer interface {
	Consume(ctx context.Context) error
}

type App struct {
	Bot    *telegramBot.Bot
	Sender *sender.Sender
	cons   Consumer
}

func NewApp() (*App, error) {
	cfg, err := config.InitConfig()
	if err != nil {
		return nil, err
	}

	tgBot, err := initTgBot(cfg)
	if tgBot == nil {
		return nil, err
	}

	presenter := out.NewTgPresenter(tgBot)
	uploader := out.NewTelegramUploader(tgBot)
	t := time.NewTicker(cfg.TickerConfig.TickTime)

	// Инициализация Kafka клиента
	c, err := kafka.NewClient(cfg.KafkaConfig)
	if err != nil {
		return nil, err
	}

	prod, err := out.NewProducer(c, cfg.KafkaConfig.Topic)
	if err != nil {
		return nil, err
	}

	db := inMemory.NewInMemory()

	// Создаем UseCase (сервис бизнес-логики)
	scs := saveConverted.NewSaveConvertedService(db, uploader)

	// Создаем Handler, передавая UseCase как EventHandler
	// Теперь scs должен реализовывать domain.EventHandler интерфейс
	consumerH := in.NewHandler(scs)

	// Убираем вызов SetConsumer, так как теперь scs сам реализует EventHandler
	// и не нуждается в отдельном консьюмере

	// Создаем Consumer с Handler
	cons, err := in.NewConsumer(c, cfg.KafkaConfig.ConsumerGroupID, consumerH)
	if err != nil {
		return nil, err
	}

	// Остальные компоненты остаются без изменений
	ss := send.NewRandomFileService(db, uploader, presenter)
	sendScheduler := sender.NewSender(t, ss)

	d := out.NewDownloader(*tgBot, cfg.DownloaderConfig)
	s := download.NewFileService(d, db, prod)
	h := dochandler.NewDocHandler(s, presenter)
	us := userservice.NewUserService(db)

	u := userhandler.NewUserHandler(us, presenter)

	router := telegramBot.NewRouter(u, h)
	bot, err := telegramBot.NewTgBot(tgBot, router)
	if err != nil {
		return nil, err
	}

	return &App{
		Bot:    bot,
		Sender: sendScheduler,
		cons:   cons,
	}, nil
}

func (a *App) RunApp(ctx context.Context) {
	// Запускаем бота
	go a.Bot.Start()

	// Запускаем отправитель файлов
	go a.Sender.Start()

	// Запускаем Kafka consumer
	go func() {
		if err := a.cons.Consume(ctx); err != nil {
			// Логируем ошибку
			// В зависимости от требований можно добавить retry логику
		}
	}()

	// Ожидаем завершения контекста
	<-ctx.Done()
}

func initTgBot(cfg *config.Config) (*tgbotapi.BotAPI, error) {

	commands := []tgbotapi.BotCommand{
		{Command: "start", Description: "Запустить бота"},
		{Command: "help", Description: "Помощь"},
		{Command: "files", Description: "Показать мои файлы"},
	}

	var tgBot *tgbotapi.BotAPI
	c := configurateClient()

	if !cfg.BotConfig.WithClient {
		log.Println("configurating standart telegram_bot_api")
		bot, err := tgbotapi.NewBotAPI(cfg.BotConfig.Token)
		if err != nil {
			return nil, err
		}
		tgBot = bot

	} else {
		tgEndpoint := tgbotapi.APIEndpoint
		botWithClient, err := tgbotapi.NewBotAPIWithClient(cfg.BotConfig.Token, tgEndpoint, c)
		if err != nil {
			return nil, err
		}
		tgBot = botWithClient
	}

	command := tgbotapi.NewSetMyCommands(commands...)

	_, err := tgBot.Request(command)
	if err != nil {
		return nil, err
	}

	return tgBot, nil
}

func configurateClient() *http.Client {

	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout:   30 * time.Second,
		ResponseHeaderTimeout: 30 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   10,
		IdleConnTimeout:       90 * time.Second,
	}

	client := &http.Client{
		Timeout:   60 * time.Second,
		Transport: transport,
	}

	return client
}
