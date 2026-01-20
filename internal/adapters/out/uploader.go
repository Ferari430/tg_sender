package out

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Uploader struct {
	bot *tgbotapi.BotAPI
}

func NewUploader(b *tgbotapi.BotAPI) *Uploader {
	return &Uploader{bot: b}
}

func (u *Uploader) UploadArchive() error {
	return nil
}
