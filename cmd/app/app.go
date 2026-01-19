package app

import (
	"context"
	"log"

	telegramBot "github.com/Ferari430/tg_sender/internal/bot"
)

type App struct {
	Bot *telegramBot.Bot
}

func NewApp() (*App, error) {
	router := telegramBot.NewRouter()
	bot, err := telegramBot.NewTgBot(router)

	if err != nil {

		return nil, err
	}

	return &App{Bot: bot}, nil
}

func (a *App) RunApp(_ context.Context) {

	log.Println(a.Bot.Bot.Self.UserName)
	a.Bot.Start()
}
