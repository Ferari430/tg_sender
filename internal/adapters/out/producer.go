package out

import (
	"github.com/IBM/sarama"
)

type Producer interface {
	PublishTaskCreated(msg *sarama.ProducerMessage) error
}

type producer struct {
	prod  sarama.SyncProducer
	topic string
}

func NewProducer(client sarama.Client, t string) (Producer, error) {

	prod, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		return nil, err
	}
	return &producer{prod: prod,
		topic: t}, err
}

func (p *producer) PublishTaskCreated(msg *sarama.ProducerMessage) error {
	_, _, err := p.prod.SendMessage(msg)
	if err != nil {
		return err
	}

	return nil
}
