package out

import (
	"fmt"
	"strings"

	"github.com/Ferari430/tg_sender/pkg"
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

func (tg *TgPresenter) Message(id int64, text string) error {
	msg := tgbotapi.NewMessage(id, text)
	_, err := tg.bot.Send(msg)
	return err
}

func (tg *TgPresenter) Welcome(id int64) error {
	msg := tgbotapi.NewMessage(id, StartMessage)
	_, err := tg.bot.Send(msg)
	return err
}

func (tg *TgPresenter) Files(id int64, fileNames []string) error {
	arr := pkg.RemoveDuplicates(fileNames)
	m := fmt.Sprintf("Твои файлы:\n%s", strings.Join(arr, "\n"))
	msg := tgbotapi.NewMessage(id, m)
	_, err := tg.bot.Send(msg)
	return err
}

func (tg *TgPresenter) Help(id int64) error {
	msg := tgbotapi.NewMessage(id, helpMessage)
	_, err := tg.bot.Send(msg)
	return err
}
