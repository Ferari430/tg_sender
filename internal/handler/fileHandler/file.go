package dochandler

import (
	"fmt"
	"log"
	"time"

	fileservice "github.com/Ferari430/tg_sender/internal/service/file"
	"github.com/Ferari430/tg_sender/pkg"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Presenter interface {
	Successes(id int64, text string) error
	Error(id int64, text string) error
	Message(id int64, text string) error
}

type DocHandler struct {
	docService *fileservice.FileService
	P          Presenter
}

func NewDocHandler(s *fileservice.FileService, presenter Presenter) *DocHandler {
	return &DocHandler{
		docService: s,
		P:          presenter,
	}
}

func (d *DocHandler) HandleMessage(msg string) {

}

func (d *DocHandler) HandleDoc(msg *tgbotapi.Message) error {

	doc := msg.Document
	if doc == nil {
		log.Println("there are no doc")
		return tgbotapi.Error{}

	}

	dto := &fileservice.DocDTO{
		OwnerID:  msg.Chat.ID,
		FileName: doc.FileName,
		FileID:   doc.FileID,
		Size:     doc.FileSize,
	}

	err := d.docService.ValidateArchive(dto)
	if err != nil {
		log.Println(err)
		return err
	}

	err = d.P.Message(msg.Chat.ID, "Архив прошел проверки! Начинаю скачивание файла!")
	if err != nil {
		return err
	}

	now := time.Now()
	err = d.docService.DownloadZip(dto)
	if err != nil {
		return err
	}
	t := pkg.TimeSpent(now)

	m := fmt.Sprintf("скачивание %s успешно завершено. Потребовалось %v", dto.FileName, t)
	log.Printf(m)
	err = d.P.Successes(msg.Chat.ID, m)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
