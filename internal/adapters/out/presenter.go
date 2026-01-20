package out

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TgPresenter struct {
	bot *tgbotapi.BotAPI
}

func NewTgPresenter(b *tgbotapi.BotAPI) *TgPresenter {
	return &TgPresenter{bot: b}
}

func (tg *TgPresenter) Successes(id int64, text string) error {
	msg := tgbotapi.NewMessage(id, text)
	_, err := tg.bot.Send(msg)
	return err
}

func (tg *TgPresenter) Error(id int64, text string) error {
	msg := tgbotapi.NewMessage(id, text)
	_, err := tg.bot.Send(msg)
	return err
}
