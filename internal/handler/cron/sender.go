package sender

import (
	"log"
	"time"
)

type Uploader interface {
	UploadArchive() error
}

type Sender struct {
	t *time.Ticker
	u Uploader
}

func NewSender(uploader Uploader, ticker *time.Ticker) *Sender {
	return &Sender{u: uploader,
		t: ticker,
	}
}

func (c *Sender) Start() {
	for range c.t.C {
		log.Println("отправка файла пользователю...")
		err := c.u.UploadArchive()
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("файл успешно отправлен")
	}

}
