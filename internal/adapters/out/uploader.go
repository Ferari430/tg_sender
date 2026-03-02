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

func (u *TelegramUploader) UploadDocument(chatID int64, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("error closing file: %v", err)
		}
	}()

	var buf bytes.Buffer
	if _, err = io.Copy(&buf, f); err != nil {
		return err
	}

	file := tgbotapi.FileBytes{
		Name:  "file.txt",
		Bytes: buf.Bytes(),
	}

	doc := tgbotapi.NewDocument(chatID, file)
	_, err = u.bot.Send(doc)
	return err
}
