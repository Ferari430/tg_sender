package models

import (
	"time"
)

type User struct {
	ID         string
	ChatID     int64 // Telegram chat id
	Username   string
	FirstName  string
	LastName   string
	CreatedAt  time.Time
	LastSeen   time.Time
	FilesCount int64
}
