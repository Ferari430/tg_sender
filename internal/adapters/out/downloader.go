package out

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Ferari430/tg_sender/internal/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Downloader struct {
	bot *tgbotapi.BotAPI
	cfg *config.DownloaderConfig
}

func (d *Downloader) DownloadZip(fileName, fileID string) error {
	log.Printf("начинаю скачивание файла %s c FileId= %s", fileName, fileID)

	file, err := d.bot.GetFile(tgbotapi.FileConfig{FileID: fileID})
	if err != nil {
		return err
	}
	url := file.Link(d.bot.Token)

	res, err := http.Get(url)
	defer func() error {
		err := res.Body.Close()
		return err
	}()

	if err != nil {
		return err
	}

	out, err := os.Create(filepath.Join(d.cfg.RootDir, fileName))
	if err != nil {
		return err
	}

	defer out.Close()

	_, err = io.Copy(out, res.Body)
	if err != nil {
		return err
	}

	return nil
}

func NewDownloader(b tgbotapi.BotAPI, c *config.DownloaderConfig) *Downloader {
	return &Downloader{bot: &b,
		cfg: c,
	}
}
