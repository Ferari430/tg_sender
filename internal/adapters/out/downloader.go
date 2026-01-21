package out

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Ferari430/tg_sender/internal/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const MaxFileSize = 10 * 1024 * 1024

type Downloader struct {
	bot *tgbotapi.BotAPI
	cfg *config.DownloaderConfig
}

func (d *Downloader) DownloadZip(fileName, fileID string) (string, error) {
	log.Printf("начинаю скачивание файла %s", fileName)

	file, err := d.bot.GetFile(tgbotapi.FileConfig{FileID: fileID})
	if err != nil {
		return "", err
	}
	url := file.Link(d.bot.Token)

	if file.FileSize > MaxFileSize {
		log.Println("file size too big")
		return "", errors.New("file size too big")
	}

	res, err := http.Get(url)
	defer func() error {
		err := res.Body.Close()
		return err
	}()

	if err != nil {
		return "", err
	}
	path := filepath.Join(d.cfg.RootDir, fileName)
	out, err := os.Create(path)
	if err != nil {
		return "", err
	}

	defer out.Close()

	_, err = io.Copy(out, res.Body)
	if err != nil {
		return "", err
	}

	return path, nil
}

func NewDownloader(b tgbotapi.BotAPI, c *config.DownloaderConfig) *Downloader {
	return &Downloader{bot: &b,
		cfg: c,
	}
}
