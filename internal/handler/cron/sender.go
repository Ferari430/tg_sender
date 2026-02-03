package sender

import (
	"log"
	"time"

	fileservice "github.com/Ferari430/tg_sender/internal/service/file"
)

type Uploader interface {
	UploadDocument(path string) error
}

type Getter interface {
	GetRandomFilePath() (string, error)
}

type Sender struct {
	t *time.Ticker
	u Uploader
	g Getter
	s fileservice.RandomFileService
}

func NewSender(uploader Uploader, ticker *time.Ticker, fileService Getter) *Sender {
	return &Sender{u: uploader,
		t: ticker,
		g: fileService,
	}
}

func (c *Sender) Start() {
	for range c.t.C {
		path, err := c.g.GetRandomFilePath()
		if err != nil {
			log.Println(err)
			continue
		}

		log.Println("найден файл:", path)
		log.Println("отправка файла пользователю...")

		err = c.u.UploadDocument(path)
		if err != nil {
			log.Println(err)
		}

		log.Println("файл успешно отправлен")
	}

}
