package fileservice

type Repository interface {
	GetRandomFilePath() (string, error)
}

type RandomFileService struct {
	db Repository
}

func (s *RandomFileService) GetRandomFilePath() (string, error) {
	return s.db.GetRandomFilePath()
}

func NewRandomFileService(database Repository) *RandomFileService {
	return &RandomFileService{
		db: database,
	}
}
