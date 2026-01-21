package fileservice

import (
	"errors"
	"log"
	"time"

	"github.com/Ferari430/tg_sender/internal/models"
)

type FileDownloader interface {
	DownloadZip(fileName, fileID string) (string, error)
}

type Reposiroty interface {
	SaveFile(f models.File) error
	GetFileByID(fileID string) (*models.File, bool)
	GetFilesByUser(chatID int64) []models.File
	DeleteFile(fileID string)
}

type FileService struct {
	downloader FileDownloader
	db         Reposiroty
}

func NewFileService(downloader FileDownloader, Database Reposiroty) *FileService {
	return &FileService{
		downloader: downloader,
		db:         Database,
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

	return nil
}

func DtoToFileModel(d *DocDTO) models.File {
	return models.File{
		ID:        d.FileID,
		OwnerID:   d.OwnerID,
		Name:      d.FileName,
		Path:      d.Path,
		Size:      d.Size,
		Extension: d.Extension,
		CreatedAt: time.Now(), // время создания файла
	}
}
