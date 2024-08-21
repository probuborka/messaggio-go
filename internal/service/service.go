package service

import (
	"context"

	"github.com/probuborka/messaggio/internal/domain"
	"github.com/probuborka/messaggio/internal/infrastructure/producer"
	"github.com/probuborka/messaggio/internal/infrastructure/repository"
)

type Message interface {
	Create(ctx context.Context, message domain.Message) error
	Process(ctx context.Context, message domain.Message) error
	Statistics(ctx context.Context) ([]domain.Message, error)
}

type Services struct {
	Message Message
}

func New(repo *repository.Repositories, producer *producer.Producer) *Services {
	return &Services{
		Message: NewMessageService(repo.Message, producer.Message),
	}
}
