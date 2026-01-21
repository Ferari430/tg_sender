package out

import (
	"fmt"

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

func (tg *TgPresenter) Welcome(id int64, name string) error {
	text := fmt.Sprintf("Привет %s! 👋  \nЯ бот для работы с файлами 📦\n\nЧто я умею:\n• принимаю ZIP-архивы  \n• сохраняю их на сервере  \n• запоминаю, кто и что загрузил  \n\nПросто отправь ZIP-файл в этот чат.\n\nДоступные команды:\n/help — помощь  \n/myfiles — посмотреть загруженные файлы", name)
	msg := tgbotapi.NewMessage(id, text)
	_, err := tg.bot.Send(msg)

	return err
}

func (tg *TgPresenter) Files(id int64, text string) error {

	msg := tgbotapi.NewMessage(id, text)
	_, err := tg.bot.Send(msg)

	return err
}

func (tg *TgPresenter) Message(id int64, text string) error {

	msg := tgbotapi.NewMessage(id, text)
	_, err := tg.bot.Send(msg)

	return err
}
