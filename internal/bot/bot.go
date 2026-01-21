package telegramBot

import (
	"log"

	docHandler "github.com/Ferari430/tg_sender/internal/handler/fileHandler"
	userhandler "github.com/Ferari430/tg_sender/internal/handler/userHandler"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	tgBot  *tgbotapi.BotAPI
	router *Router
}

type Router struct {
	userHandler *userhandler.UserHandler
	docHandler  *docHandler.DocHandler
}

func NewRouter(u *userhandler.UserHandler, d *docHandler.DocHandler) *Router {
	return &Router{
		userHandler: u,
		docHandler:  d,
	}
}

func NewTgBot(tgBot *tgbotapi.BotAPI, r *Router) (*Bot, error) {
	return &Bot{
		tgBot:  tgBot,
		router: r,
	}, nil
}

func (b *Bot) Start() {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	updates := b.tgBot.GetUpdatesChan(updateConfig)

	log.Println("Bot started, waiting for updates...")

	// Обрабатываем обновления
	for update := range updates {
		b.HandleMessage(update)
	}

}

func (b *Bot) HandleMessage(u tgbotapi.Update) {
	if u.Message == nil {
		log.Println("message is nil")
		return
	}

	id := u.Message.Chat.ID

	if u.Message.Document != nil {
		fName := u.Message.Document.FileName
		log.Println("filename:", fName)
		err := b.router.docHandler.HandleDoc(u.Message)
		if err != nil {
			log.Println(err)
		}
		return
	}

	if u.Message != nil {
		log.Println("recieved message:", u.Message.Text)
		msg := tgbotapi.NewMessage(id, "info")

		b.router.userHandler.HandleMessage(u.Message)
		_, err := b.tgBot.Send(msg)
		if err != nil {
			log.Println(err)
		}
	}

}
