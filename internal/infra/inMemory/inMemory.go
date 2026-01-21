package inMemory

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/Ferari430/tg_sender/internal/models"
)

type InMemory struct {
	mu    sync.Mutex
	users map[int64]*models.User
	files map[string]*models.File
}

func NewInMemory() *InMemory {
	return &InMemory{
		users: make(map[int64]*models.User),
		files: make(map[string]*models.File),
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

func (im *InMemory) FileNames(id int64) ([]string, error) {
	result := make([]string, len(im.files))
	for _, file := range im.files {
		if file.OwnerID == id {
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

// files

func (im *InMemory) SaveFile(f models.File) error {
	im.mu.Lock()
	defer im.mu.Unlock()

	if _, exists := im.files[f.ID]; exists {
		return nil
	}

	f.CreatedAt = time.Now()
	im.files[f.ID] = &f
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

	f, exists := im.files[fileID]
	if !exists {
		return nil, false
	}

	return f, true
}

func (im *InMemory) GetFileByName(name string) bool {
	for _, file := range im.files {
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
	for _, f := range im.files {
		if f.OwnerID == chatID {
			res = append(res, *f)
		}
	}
	return res
}

func (im *InMemory) DeleteFile(fileID string) {
	im.mu.Lock()
	defer im.mu.Unlock()

	delete(im.files, fileID)
	log.Println("Deleted file:", fileID)
}
