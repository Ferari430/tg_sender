package inMemory

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/Ferari430/tg_sender/internal/domain/models"
)

type InMemory struct {
	mu       sync.Mutex
	users    map[int64]*models.User //chatID:User
	archives map[string]*models.File
	files    map[int64][]*models.File //chatID:PDF
}

func NewInMemory() *InMemory {
	return &InMemory{
		users:    make(map[int64]*models.User),
		archives: make(map[string]*models.File),
	}
}

func (im *InMemory) SaveUser(u models.User) error {
	im.mu.Lock()
	defer im.mu.Unlock()

	if _, exists := im.users[u.ChatID]; !exists {
		u.CreatedAt = time.Now()
	}
	u.LastSeen = time.Now()
	im.users[u.ChatID] = &u

	log.Println("Saved user", u.ChatID)
	return nil
}

func (im *InMemory) GetUserById(chatID int64) (*models.User, bool) {
	im.mu.Lock()
	defer im.mu.Unlock()

	u, exists := im.users[chatID]
	if !exists {
		return nil, false
	}
	return u, true
}

func (im *InMemory) FileNames(ownerId int64) ([]string, error) {
	result := make([]string, len(im.archives))
	for _, file := range im.archives {
		if file.OwnerID == ownerId {
			result = append(result, file.Name)
		}
	}
	if len(result) == 0 {
		return result, errors.New("not found")
	}

	return result, nil
}

func (im *InMemory) Exists(chatID int64) bool {
	im.mu.Lock()
	defer im.mu.Unlock()

	_, exists := im.users[chatID]
	return exists
}

// Save archive
func (im *InMemory) SaveFile(f models.File) error {
	im.mu.Lock()
	defer im.mu.Unlock()

	// подумать надо ли проверять
	if _, exists := im.archives[f.ID]; exists {
		return nil
	}
	f.CreatedAt = time.Now()
	im.archives[f.ID] = &f
	v, ok := im.users[f.OwnerID]
	if ok {
		v.FilesCount++
		log.Println("количество файлов у пользователя", v.Username, "=", v.FilesCount)
	}

	log.Println("Saved file:", f.Name, "owner:", f.OwnerID)
	return nil
}

func (im *InMemory) GetFileByID(fileID string) (*models.File, bool) {
	im.mu.Lock()
	defer im.mu.Unlock()

	f, exists := im.archives[fileID]
	if !exists {
		return nil, false
	}

	return f, true
}

func (im *InMemory) GetFileByName(name string) bool {
	for _, file := range im.archives {
		if file.Name == name {
			return true
		}
	}

	return false
}

func (im *InMemory) GetFilesByUser(chatID int64) []models.File {
	im.mu.Lock()
	defer im.mu.Unlock()

	var res []models.File
	for _, f := range im.archives {
		if f.OwnerID == chatID {
			res = append(res, *f)
		}
	}
	return res
}

func (im *InMemory) DeleteFile(fileID string) {
	im.mu.Lock()
	defer im.mu.Unlock()

	delete(im.archives, fileID)
	log.Println("Deleted file:", fileID)
}

func (im *InMemory) GetRandomFilePath() (string, error) {
	if len(im.archives) == 0 {
		return "", fmt.Errorf("база данных пуста")
	}

	randomIndex := rand.Intn(len(im.archives))
	i := 0
	for _, file := range im.archives {
		if i == randomIndex {
			return file.Path, nil
		}
		i++
	}

	return "", errors.New("not found")
}

func (im InMemory) ChatId() ([]int64, error) {
	if !(len(im.users) > 0) {
		return nil, errors.New("there are no users")
	}

	chatIds := make([]int64, len(im.users))
	for _, user := range im.users {
		chatIds = append(chatIds, user.ChatID)
	}
	return chatIds, nil
}

func (im *InMemory) GetRandomPDFPathForEachUser(chatID []int64) (map[int64]*models.File, error) {
	result := make(map[int64]*models.File) //ChatId ->file

	if len(im.files) == 0 {
		return nil, fmt.Errorf("база данных пуста")
	}

	for _, id := range chatID {
		userFiles, ok := im.files[id]
		if !ok {
			log.Println("not found", id)
		}
		if len(userFiles) < 1 {
			log.Println("len arr < 1")
		}

		r := rand.Intn(len(userFiles))
		randomFile := userFiles[r]
		result[id] = randomFile
	}

	return result, errors.New("not found")
}

func (im *InMemory) SavePDF(f models.File) error {
	im.mu.Lock()
	defer im.mu.Unlock()

	// подумать надо ли проверять
	if _, exists := im.archives[f.ID]; exists {
		return nil
	}
	f.CreatedAt = time.Now()
	im.archives[f.ID] = &f
	v, ok := im.users[f.OwnerID]
	if ok {
		v.FilesCount++
		log.Println("количество файлов у пользователя", v.Username, "=", v.FilesCount)
	}

	log.Println("Saved file:", f.Name, "owner:", f.OwnerID)
	return nil
}

func (im *InMemory) Save() error {
	return nil
}
