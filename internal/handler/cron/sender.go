package sender

import (
	"log"
	"time"
)

type Getter interface {
	UploadDocument() error
}

type Sender struct {
	t *time.Ticker
	g Getter
}

func NewSender(ticker *time.Ticker, fileService Getter) *Sender {
	return &Sender{
		t: ticker,
		g: fileService,
	}
}

func (c *Sender) Start() {

	for range c.t.C {
		log.Println("tick")
		err := c.g.UploadDocument()
		if err != nil {
			log.Println("START Sender error:", err)
			return
		}
	}
}
