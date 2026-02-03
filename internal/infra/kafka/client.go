package kafka

import (
	"github.com/Ferari430/tg_sender/internal/config"
	"github.com/IBM/sarama"
)

func NewClient(kafkaCfg config.KafkaConfig) (sarama.Client, error) {
	saramaCfg := sarama.NewConfig()
	addr := []string{kafkaCfg.KafkaPort}

	client, err := sarama.NewClient(addr, saramaCfg)

	if err != nil {
		return nil, err
	}

	return client, nil
}
