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
	fileservice "github.com/Ferari430/tg_sender/internal/service/file"
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

	tgBot, err := tgbotapi.NewBotAPI(cfg.BotConfig.Token)

	if err != nil {
		return nil, err
	}

	presenter := out.NewTgPresenter(tgBot)

	uploader := out.NewUploader(tgBot)

	t := time.NewTicker(cfg.TickerConfig.TickTime)

	send := sender.NewSender(uploader, t)

	d := out.NewDownloader(*tgBot, cfg.DownloaderConfig)
	s := fileservice.NewFileService(d)
	h := dochandler.NewDocHandler(s, presenter)
	u := userhandler.NewUserHandler()

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
	go a.Sender.Start()
	select {}
}

func initSender() {

}

func initTgBot() {

}
