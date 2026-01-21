package app

import (
	"context"
	"time"

	"github.com/Ferari430/tg_sender/internal/adapters/out"
	telegramBot "github.com/Ferari430/tg_sender/internal/bot"
	"github.com/Ferari430/tg_sender/internal/config"
	"github.com/Ferari430/tg_sender/internal/handler/cron"
	dochandler "github.com/Ferari430/tg_sender/internal/handler/fileHandler"
	userhandler "github.com/Ferari430/tg_sender/internal/handler/userHandler"
	"github.com/Ferari430/tg_sender/internal/infra/inMemory"
	fileservice "github.com/Ferari430/tg_sender/internal/service/file"
	userservice "github.com/Ferari430/tg_sender/internal/service/userService"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type App struct {
	Bot    *telegramBot.Bot
	Sender *sender.Sender
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

	uploader := out.NewUploader(tgBot)
	t := time.NewTicker(cfg.TickerConfig.TickTime)

	send := sender.NewSender(uploader, t)

	db := inMemory.NewInMemory()

	d := out.NewDownloader(*tgBot, cfg.DownloaderConfig)
	s := fileservice.NewFileService(d, db)
	h := dochandler.NewDocHandler(s, presenter)
	us := userservice.NewUserService(db)

	u := userhandler.NewUserHandler(us, presenter)

	router := telegramBot.NewRouter(u, h)
	bot, err := telegramBot.NewTgBot(tgBot, router)
	if err != nil {

		return nil, err
	}

	return &App{Bot: bot,
		Sender: send,
	}, nil
}

func (a *App) RunApp(_ context.Context) {

	go a.Bot.Start()
	//go a.Sender.Start()
	select {}
}

func initSender() {

}

func initTgBot(cfg *config.Config) (*tgbotapi.BotAPI, error) {
	commands := []tgbotapi.BotCommand{
		{Command: "start", Description: "Запустить бота"},
		{Command: "help", Description: "Помощь"},
		{Command: "files", Description: "Показать мои файлы"},
	}

	tgBot, err := tgbotapi.NewBotAPI(cfg.BotConfig.Token)
	if err != nil {
		return nil, err
	}

	command := tgbotapi.NewSetMyCommands(commands...)

	_, err = tgBot.Request(command)
	if err != nil {
		return nil, err
	}

	return tgBot, nil
}
