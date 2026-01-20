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
	Bot *telegramBot.Bot
}

func NewApp() (*App, error) {

	cfg, err := config.InitConfig()
	if err != nil {
		return nil, err
	}

	tgbot, err := tgbotapi.NewBotAPI(cfg.BotConfig.Token)

	if err != nil {
		return nil, err
	}

	presenter := out.NewTgPresenter(tgbot)

	uploader := out.NewUploader(tgbot)

	t := time.NewTicker(cfg.TickerConfig.TickTime)
	cron := cron.NewCron(uploader, t)
	_ = cron

	d := out.NewDownloader(*tgbot, cfg.DownloaderConfig)
	s := fileservice.NewFileService(d)
	h := dochandler.NewDocHandler(s, presenter)
	u := userhandler.NewUserHandler()

	router := telegramBot.NewRouter(u, h)
	bot, err := telegramBot.NewTgBot(tgbot, router)
	if err != nil {

		return nil, err
	}

	return &App{Bot: bot}, nil
}

func (a *App) RunApp(_ context.Context) {

	a.Bot.Start()
}
