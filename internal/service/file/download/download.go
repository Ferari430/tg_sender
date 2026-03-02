package download

import (
	"errors"
	"log"
	"slices"
	"strings"
	"time"

	"github.com/Ferari430/tg_sender/internal/domain/events"
	"github.com/Ferari430/tg_sender/internal/domain/models"
	"github.com/Ferari430/tg_sender/internal/infra/kafka"
	"github.com/IBM/sarama"
	"github.com/google/uuid"
)

type FileDownloader interface {
	DownloadZip(fileName, fileID string) (string, error)
}

type Reposiroty interface {
	SaveFile(f models.File) error
	GetFileByID(fileID string) (*models.File, bool)
	GetFilesByUser(chatID int64) []models.File
	DeleteFile(fileID string)
	GetFileByName(name string) bool
}

type Producer interface {
	PublishTaskCreated(msg *sarama.ProducerMessage) error
}

type FileService struct {
	downloader    FileDownloader
	db            Reposiroty
	kafkaProducer Producer
}

func NewFileService(downloader FileDownloader, Database Reposiroty, Prod Producer) *FileService {
	return &FileService{
		downloader:    downloader,
		db:            Database,
		kafkaProducer: Prod,
	}
}

func (fs *FileService) DownloadZip(dto *DocDTO) error {
	if dto == nil {
		return errors.New("nil pointer")
	}

	path, err := fs.downloader.DownloadZip(dto.FileName, dto.FileID)
	if err != nil {
		return err
	}

	log.Println("File downloaded:", path)
	dto.Path = path
	file := DtoToFileModel(dto)

	err = fs.db.SaveFile(file)
	if err != nil {
		return err
	}

	ev := event.TaskCreated{
		EventID:   uuid.New().String(),
		EventType: "File_Created",
		TaskID:    uuid.New().String(),
		ChatID:    file.OwnerID,
		FilePath:  file.Path,
		FileName:  file.Name,
		CreatedAt: file.CreatedAt,
	}

	msg, err := kafka.TaskCreatedToMessage("mytopic", ev)
	if err != nil {
		return err
	}

	err = fs.kafkaProducer.PublishTaskCreated(msg)

	if err != nil {
		return err
	}

	return nil
}

func (fs *FileService) ValidateArchive(dto *DocDTO) error {
	archiveExtensions := []string{"zip", "tar.gz", "tgz", "7z", "rar"}
	log.Println("Validating archive:", dto.FileName)

	parts := strings.Split(dto.FileName, ".")
	ext := parts[len(parts)-1]
	if !slices.Contains(archiveExtensions, ext) {
		return errors.New("invalid file extension")
	}
	dto.Extension = ext

	return nil
}

func (fs *FileService) AlreadyExisted(dto *DocDTO) bool {
	if ok := fs.db.GetFileByName(dto.FileName); ok {
		return true
	}
	return false
}

func DtoToFileModel(d *DocDTO) models.File {
	return models.File{
		ID:        d.FileID,
		OwnerID:   d.OwnerID,
		Name:      d.FileName,
		Path:      d.Path,
		Size:      d.Size,
		Extension: d.Extension,
		CreatedAt: time.Now(),
	}
}
