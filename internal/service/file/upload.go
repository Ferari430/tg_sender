package fileservice

type FileDownloader interface {
	DownloadZip(fileName, fileID string) error
}

type FileService struct {
	downloader FileDownloader
}

func NewFileService(downloader FileDownloader) *FileService {
	return &FileService{
		downloader: downloader,
	}

}

func (fs *FileService) UploadZip(dto DocDTO) error {
	return fs.downloader.DownloadZip(dto.FileName, dto.FileID)
}
