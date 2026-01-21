package models

import "time"

type File struct {
	ID        string    // уникальный идентификатор файла (можно FileID от Telegram)
	OwnerID   int64     // ChatID пользователя, который загрузил файл
	Name      string    // имя файла
	Path      string    // путь на диске, куда сохранён
	Size      int       // размер в байтах
	Extension string    // расширение файла
	CreatedAt time.Time // время загрузки
}
