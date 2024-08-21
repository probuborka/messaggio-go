package producer

import (
	"context"
	"encoding/json"

	"github.com/probuborka/messaggio/internal/domain"
)

type MessageProducer struct {
	producer producer
}

func newMessageProducer(producer producer) *MessageProducer {
	return &MessageProducer{
		producer: producer,
	}
}

func (m *MessageProducer) SendMessage(ctx context.Context, message domain.Message) error {
	//query
	bytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = m.producer.SendMessage(ctx, bytes, "message")
	if err != nil {
		return err
	}

	return nil
}
