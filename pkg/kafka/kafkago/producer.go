package kafkago

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/probuborka/messaggio/internal/domain"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(kafkaConfig domain.KafkaConfig) (*Producer, error) {
	brokers := []string{fmt.Sprintf("%s:%s", kafkaConfig.Host, kafkaConfig.Port)}

	dialer := &kafka.Dialer{
		Timeout:  10 * time.Second,
		ClientID: "registration-service-client",
	}

	c := kafka.WriterConfig{
		Brokers: brokers,
		//Topic:        "message",
		Balancer:     &kafka.LeastBytes{},
		Dialer:       dialer,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
		BatchSize:    1,
	}

	return &Producer{
		writer: kafka.NewWriter(c),
	}, nil
}

func (p Producer) SendMessage(ctx context.Context, msg []byte, topic string) error {

	key := uuid.New().String()

	message := kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: msg,
	}

	return p.writer.WriteMessages(ctx, message)

}

func (p Producer) Close() {

	p.writer.Close()

}
