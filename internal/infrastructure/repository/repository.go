package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/probuborka/messaggio/internal/domain"
	pgsql "github.com/probuborka/messaggio/internal/infrastructure/repository/postgresql"
)

type Message interface {
	Create(ctx context.Context, message domain.Message) error
	ReadAll(ctx context.Context) ([]domain.Message, error)
	Process(ctx context.Context, id domain.Id, dateProcessedStart time.Time) error
}

type Repositories struct {
	Message Message
}

func New(db *pgxpool.Pool) *Repositories {
	return &Repositories{
		Message: pgsql.NewMessage(db),
	}
}
