package out

import (
	"bytes"
	"io"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Uploader struct {
	bot *tgbotapi.BotAPI
}

func NewUploader(b *tgbotapi.BotAPI) *Uploader {
	return &Uploader{bot: b}
}

func (u *Uploader) UploadArchive() error {

	path := `B:\data\curl.txt`

	f, err := os.Open(path)
	if err != nil {
		return err
	}

	defer f.Close()
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
