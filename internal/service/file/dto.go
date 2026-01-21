package fileservice

type DocDTO struct {
	OwnerID   int64  // ChatID пользователя
	FileName  string // имя файла
	FileID    string // Telegram FileID
	Extension string // расширение файла
	Size      int    // размер файла в байтах (опционально)
	Path      string // путь к файлу на диске (опционально)
}
