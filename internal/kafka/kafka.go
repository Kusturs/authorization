package kafka

import (
	"log"

	"github.com/Shopify/sarama"
)

type Producer struct {
	producer sarama.SyncProducer
}

func NewKafkaProducer(brokers []string) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &Producer{producer: producer}, nil
}

func (kp *Producer) SendMessage(topic, message string) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	_, _, err := kp.producer.SendMessage(msg)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
		return err
	}

	return nil
}

func (kp *Producer) Close() error {
	return kp.producer.Close()
}
