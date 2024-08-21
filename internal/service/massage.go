package service

import (
	"context"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/probuborka/messaggio/internal/domain"
	"github.com/probuborka/messaggio/internal/infrastructure/producer"
	"github.com/probuborka/messaggio/internal/infrastructure/repository"
)

type MessageService struct {
	repo    repository.Message
	produce producer.Message
}

func NewMessageService(repo repository.Message, produce producer.Message) *MessageService {
	return &MessageService{
		repo:    repo,
		produce: produce,
	}
}

func (m MessageService) Create(ctx context.Context, message domain.Message) error {

	//data
	id := uuid.New().String()

	message.Id = domain.Id(id)

	//db
	if err := m.repo.Create(ctx, message); err != nil {
		return err
	}

	//broker
	if err := m.produce.SendMessage(ctx, message); err != nil {
		return err
	}

	return nil
}

func (m MessageService) Process(ctx context.Context, message domain.Message) error {
	//
	dateProcessedStart := time.Now()

	//Искусственная задержка от 1 до 10 секунд
	src := rand.NewSource(time.Now().Unix())
	randomNumber := src.Int63()
	s := randomNumber%10 + 1
	time.Sleep(time.Duration(s) * time.Second)

	//
	err := m.repo.Process(ctx, message.Id, dateProcessedStart)
	if err != nil {
		return err
	}

	return nil
}

func (m MessageService) Statistics(ctx context.Context) ([]domain.Message, error) {

	messages, err := m.repo.ReadAll(ctx)
	if err != nil {
		return nil, err
	}

	return messages, nil
}
