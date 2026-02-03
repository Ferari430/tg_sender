package out

import (
	"log"

	event "github.com/Ferari430/tg_sender/internal/domain/events"
	"github.com/IBM/sarama"
)

type Producer interface {
	SendMessage(message event.TaskCreated) error
}

type producer struct {
	prod sarama.SyncProducer
}

func NewProducer(client sarama.Client) (Producer, error) {
	
	prod, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		return nil, err
	}
	return &producer{prod: prod}, err
}

func (p *producer) SendMessage(message event.TaskCreated) error {
	log.Println("sending message...:", message)
	return nil
}
