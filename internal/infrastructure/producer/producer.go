package producer

import (
	"context"

	"github.com/probuborka/messaggio/internal/domain"
)

type producer interface {
	SendMessage(ctx context.Context, msg []byte, topic string) error
}

type Message interface {
	SendMessage(ctx context.Context, message domain.Message) error
}

type Producer struct {
	Message Message
}

func New(producer producer) *Producer {
	return &Producer{
		Message: newMessageProducer(producer),
	}
}
