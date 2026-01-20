package cron

import (
	"time"
)

type Uploader interface {
	UploadArchive() error
}

type Cron struct {
	t *time.Ticker
	u Uploader
}

func NewCron(uploader Uploader, ticker *time.Ticker) *Cron {
	return &Cron{u: uploader,
		t: ticker,
	}
}
