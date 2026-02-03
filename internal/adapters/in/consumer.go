package in

import "github.com/IBM/sarama"

type Consumer interface {
}

type consumer struct {
	consumerGroup sarama.ConsumerGroup
}

func NewConsumer(client sarama.Client, groupID string) (Consumer, error) {
	cons, err := sarama.NewConsumerGroupFromClient(groupID, client)
	if err != nil {
		return nil, err
	}
	return &consumer{consumerGroup: cons}, err
}
