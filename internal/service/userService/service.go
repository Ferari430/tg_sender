package userservice

import (
	"time"

	"github.com/Ferari430/tg_sender/internal/models"
)

type Reposiroty interface {
	SaveUser(u models.User) error
	GetUserById(chatID int64) (*models.User, bool)
	Exists(chatID int64) bool
	FileNames(id int64) ([]string, error)
}

type UserService struct {
	db Reposiroty
}

func NewUserService(database Reposiroty) *UserService {
	return &UserService{
		db: database,
	}
}

func (us *UserService) Start(dto UserDTO) error {
	user := DtoToUserModel(dto)
	err := us.db.SaveUser(user)
	if err != nil {
		return err
	}

	return nil
}

func (us *UserService) Help() {

}

// запрос готовых файлов
func (us *UserService) Files(id int64) ([]string, error) {
	files, err := us.db.FileNames(id)
	if err != nil {
		return nil, err
	}
	return files, nil
}

func DtoToUserModel(d UserDTO) models.User {
	return models.User{
		ChatID:    d.ChatID,
		Username:  d.Username,
		FirstName: d.FirstName,
		LastName:  d.LastName,
		CreatedAt: time.Now(),
		LastSeen:  time.Now(),
	}
}
