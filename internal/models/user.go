package models

import "time"

type User struct {
	TelegramID int64 // внутренний ID (DB)
	ChatID     int64 // Telegram chat id
	Username   string
	FirstName  string
	LastName   string
	CreatedAt  time.Time
	LastSeen   time.Time
	FilesCount int64
}
