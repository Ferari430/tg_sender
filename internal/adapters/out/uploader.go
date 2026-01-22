package out

import (
	"bytes"
	"io"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramUploader struct {
	bot *tgbotapi.BotAPI
}

func NewTelegramUploader(b *tgbotapi.BotAPI) *TelegramUploader {
	return &TelegramUploader{bot: b}
}

// отправка архива пользователю
func (u *TelegramUploader) UploadDocument(path string) error {

	//path := `B:\data\curl.txt`

	f, err := os.Open(path)
	if err != nil {
		return err
	}

	defer func() error {
		if err := f.Close(); err != nil {
			log.Printf("error closing file: %v", err)
			return err
		}
		return nil
	}()

	var buf bytes.Buffer

	_, err = io.Copy(&buf, f)
	if err != nil {
		return err
	}

	log.Println(buf.Len())

	file := tgbotapi.FileBytes{"name.txt", buf.Bytes()}
	doc := tgbotapi.NewDocument(449237834, file)

	_, err = u.bot.Send(doc)
	if err != nil {
		return err
	}

	return nil
}
