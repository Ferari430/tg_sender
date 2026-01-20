package dochandler

import (
	"log"
	"strings"

	fileservice "github.com/Ferari430/tg_sender/internal/service/file"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Presenter interface {
	Successes(id int64, text string) error
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

	parts := strings.Split(doc.FileName, ".")
	ext := parts[len(parts)-1]
	dto := fileservice.DocDTO{
		UserID:    msg.Chat.ID,
		FileName:  doc.FileName,
		Extension: ext,
		FileID:    doc.FileID,
	}

	err := d.docService.UploadZip(dto)
	if err != nil {
		return err
	}

	log.Printf("скачивание %s успешно завершено", dto.FileName)

	err = d.P.Successes(msg.Chat.ID, "скачивание успешно завершено")
	if err != nil {
		log.Println(err)
		return err
	}

	return nil

}
