package telegramBot

import (
	"log"

	"github.com/Ferari430/tg_sender/internal/config"
	userhandler "github.com/Ferari430/tg_sender/internal/handler/userHandler"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	Bot    *tgbotapi.BotAPI
	Router *Router
}

type Router struct {
	userHandler *userhandler.UserHandler
}

func NewRouter() *Router {
	return &Router{}
}

func NewTgBot(r *Router) (*Bot, error) {
	cfg, err := config.InitConfig()
	if err != nil {
		return nil, err
	}

	log.Println(cfg.BotConfig.Token)
	bot, err := tgbotapi.NewBotAPI(cfg.BotConfig.Token)

	if err != nil {
		return nil, err
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)
	return &Bot{Bot: bot,
		Router: r,
	}, nil
}

func (b *Bot) Start() error {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	updates := b.Bot.GetUpdatesChan(updateConfig)

	log.Println("Bot started, waiting for updates...")

	// Обрабатываем обновления
	for update := range updates {
		b.Router.HandleMessage(update)
	}

	return nil

}

func (r *Router) HandleMessage(u tgbotapi.Update) {
	if u.Message != nil {
		log.Println("recieved message:", u.Message.Text)
	}

	if u.Message.Document != nil {
		fname := u.Message.Document.FileName
		log.Println("filename:", fname)
	}
}
