package kafkasarama

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"github.com/probuborka/messaggio/internal/domain"
)

type Producer struct {
	producer sarama.SyncProducer
}

func NewProducer(kafkaConfig domain.KafkaConfig) (*Producer, error) {

	producer, err := sarama.NewSyncProducer([]string{fmt.Sprintf("%s:%s", kafkaConfig.Host, kafkaConfig.Port)}, nil)
	if err != nil {
		return nil, err
	}

	return &Producer{
		producer: producer,
	}, err
}

func (p Producer) SendMessage(ctx context.Context, msg []byte, topic string) error {

	requestID := uuid.New().String()

	proMsg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(requestID),
		Value: sarama.ByteEncoder(msg),
	}

	_, _, err := p.producer.SendMessage(proMsg)
	if err != nil {
		return err
	}
	return nil
}

func (p Producer) Close() {

	p.producer.Close()

}
