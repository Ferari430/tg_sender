package send

import (
	"log"

	"github.com/Ferari430/tg_sender/internal/domain/models"
)

type Repository interface {
	GetRandomPDFPathForEachUser(chatID []int64) (map[int64]*models.File, error)
	ChatId() ([]int64, error)
}

type Uploader interface {
	UploadDocument(chatID int64, path string) error
}

type Presenter interface {
	Successes(id int64, text string) error
	Error(id int64, text string) error
	Message(id int64, text string) error
}

type RandomFileService struct {
	db Repository
	u  Uploader
	p  Presenter
}

func NewRandomFileService(
	database Repository,
	uploader Uploader,
	presenter Presenter,
) *RandomFileService {
	return &RandomFileService{
		db: database,
		u:  uploader,
		p:  presenter,
	}
}

func (s *RandomFileService) UploadDocument() error {
	chatIDs, err := s.db.ChatId()
	if err != nil {
		return err
	}

	result, err := s.db.GetRandomPDFPathForEachUser(chatIDs)
	if err != nil {
		return err
	}

	for chatID, file := range result {
		log.Printf("найден файл %v для пользователя %d", file.Path, chatID)

		err = s.u.UploadDocument(chatID, file.Path)
		if err != nil {
			_ = s.p.Error(chatID, "Ошибка отправки файла")
			continue
		}

		_ = s.p.Successes(chatID, "Файл успешно отправлен")
	}

	return nil
}
