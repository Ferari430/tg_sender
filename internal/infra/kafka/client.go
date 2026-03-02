package kafka

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Ferari430/tg_sender/internal/config"
	"github.com/IBM/sarama"
)

func NewClient(kafkaCfg config.KafkaConfig) (sarama.Client, error) {
	log.Printf("Connecting to Kafka at: %s", kafkaCfg.BrokersAddr)

	// Проверяем что адрес не пустой
	if kafkaCfg.BrokersAddr == "" {
		return nil, fmt.Errorf("Kafka broker address is empty")
	}

	// Проверяем формат адреса
	if !strings.Contains(kafkaCfg.BrokersAddr, ":") {
		return nil, fmt.Errorf("invalid Kafka address format: %s, expected host:port", kafkaCfg.BrokersAddr)
	}

	saramaCfg := sarama.NewConfig()
	saramaCfg.Version = sarama.V3_6_0_0
	saramaCfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	saramaCfg.Producer.Return.Successes = true
	sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)

	saramaCfg.Net.DialTimeout = 5 * time.Second
	saramaCfg.Net.ReadTimeout = 5 * time.Second
	saramaCfg.Net.WriteTimeout = 5 * time.Second
	saramaCfg.Net.KeepAlive = 5 * time.Second

	addr := []string{kafkaCfg.BrokersAddr}

	log.Printf("Trying to connect to Kafka brokers: %v", addr)

	client, err := sarama.NewClient(addr, saramaCfg)
	if err != nil {
		log.Printf("Failed to connect to Kafka: %v", err)
		return nil, err
	}

	log.Println("Successfully connected to Kafka")

	// Проверяем доступность брокеров
	brokers := client.Brokers()
	log.Printf("Available brokers: %d", len(brokers))
	for _, broker := range brokers {
		log.Printf("Broker: %s", broker.Addr())
	}

	if err := createTopicIfNotExists(client, kafkaCfg.Topic); err != nil {
		log.Printf("Failed to create topic: %v", err)
		return nil, err
	}

	// Проверяем доступность топика
	topics, err := client.Topics()
	if err != nil {
		log.Printf("Warning: cannot list topics: %v", err)
		return nil, err
	}

	for _, topic := range topics {
		log.Printf("Topic: %s", topic)
	}

	return client, nil
}

func createTopicIfNotExists(client sarama.Client, topic string) error {
	// 1. Создаем ClusterAdmin
	admin, err := sarama.NewClusterAdminFromClient(client)
	if err != nil {
		return fmt.Errorf("failed to create admin: %w", err)
	}

	// 2. Проверяем существует ли топик
	topics, err := admin.ListTopics()
	if err != nil {
		return fmt.Errorf("failed to list topics: %w", err)
	}

	if _, exists := topics[topic]; exists {
		log.Printf("Topic %s already exists", topic)
		return nil
	}

	// 3. Настраиваем топик
	topicDetail := &sarama.TopicDetail{
		NumPartitions:     1,
		ReplicationFactor: 1,
		ConfigEntries: map[string]*string{
			"retention.ms":   stringPtr("604800000"), // Retention 7 дней
			"cleanup.policy": stringPtr("delete"),
		},
	}

	// 4. Создаем топик
	err = admin.CreateTopic(topic, topicDetail, false)
	if err != nil {
		return fmt.Errorf("failed to create topic %s: %w", topic, err)
	}

	time.Sleep(time.Second * 5)

	log.Printf("Topic %s created successfully", topic)
	return nil
}

func stringPtr(s string) *string {
	return &s
}
