package fileservice

type Repository interface {
	GetRandomFilePath() (string, error)
}
type Uploader interface {
	UploadDocument(path string) error
}

type Consumer interface {
}

type RandomFileService struct {
	db   Repository
	u    Uploader
	cons Consumer
}

func (s *RandomFileService) GetRandomFilePath() (string, error) {
	return s.db.GetRandomFilePath()
}

func NewRandomFileService(database Repository, consumer Consumer) *RandomFileService {
	return &RandomFileService{
		db:   database,
		cons: consumer,
	}
}
