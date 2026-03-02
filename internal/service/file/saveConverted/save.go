package saveConverted

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	event "github.com/Ferari430/tg_sender/internal/domain/events"
)

type Repository interface {
	Save() error
	// Добавляем другие необходимые методы
}

type Uploader interface {
	// Добавляем методы uploader'а
	UploadDocument(chatID int64, path string) error
}

type SaveConverted struct {
	db Repository
	u  Uploader
	// Убираем зависимость от FileConsumer, так как теперь сервис сам обрабатывает события
}

func NewSaveConvertedService(db Repository, uploader Uploader) *SaveConverted {
	return &SaveConverted{
		db: db,
		u:  uploader,
	}
}

func (sc *SaveConverted) Handle(msg *event.KafkaMessage) error {
	// 1. Десериализуем бизнес-ивент из сообщения Kafka
	var taskEvent event.TaskConverted
	if err := json.Unmarshal(msg.Value, &taskEvent); err != nil {
		return fmt.Errorf("failed to unmarshal task event: %w", err)
	}

	// 2. Валидируем ивент
	if err := sc.validateTaskEvent(&taskEvent); err != nil {
		return fmt.Errorf("invalid task event: %w", err)
	}

	log.Printf("Processing task event: %s for chat %d", taskEvent.TaskID, taskEvent.ChatID)

	// 3. Выполняем бизнес-логику
	return sc.processTaskConverted(&taskEvent)
}

func (sc *SaveConverted) validateTaskEvent(event *event.TaskConverted) error {
	// Проверяем обязательные поля
	if event.TaskID == "" {
		return errors.New("task_id is required")
	}
	if event.ChatID == 0 {
		return errors.New("chat_id is required")
	}
	if event.FilePath == "" {
		return errors.New("file_path is required")
	}
	return nil
}

func (sc *SaveConverted) processTaskConverted(event *event.TaskConverted) error {
	// 1. Сохраняем информацию о задаче в БД
	if err := sc.db.Save(); err != nil {
		return fmt.Errorf("failed to save task: %w", err)
	}

	// 2. Если есть uploader, используем его для загрузки файла
	if sc.u != nil {
		if err := sc.u.UploadDocument(event.ChatID, event.FilePath); err != nil {
			return fmt.Errorf("failed to upload file: %w", err)
		}
	}

	// 3. Дополнительная бизнес-логика (конвертация, обработка и т.д.)
	// ...

	log.Printf("Successfully processed task %s for chat %d", event.TaskID, event.ChatID)
	return nil
}

// Устаревшие методы - можно удалить или оставить для обратной совместимости
func (sc *SaveConverted) RecieveFiles() error {
	// Этот метод больше не используется в новой архитектуре
	// Оставляем для обратной совместимости или удаляем
	return errors.New("method RecieveFiles is deprecated, use Handle instead")
}
