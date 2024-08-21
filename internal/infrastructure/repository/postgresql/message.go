package pgsql

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/probuborka/messaggio/internal/domain"
)

type Message struct {
	db *pgxpool.Pool
}

func NewMessage(db *pgxpool.Pool) *Message {
	return &Message{
		db: db,
	}
}

func (m *Message) Create(ctx context.Context, message domain.Message) error {
	//query
	query := `INSERT INTO messages (id, message) VALUES (@id, @messages)`

	//data
	args := pgx.NamedArgs{
		"id":       message.Id,
		"messages": message.Message,
	}

	//db
	_, err := m.db.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("unable to insert row: %w", err)
	}

	return nil
}

func (m *Message) ReadAll(ctx context.Context) ([]domain.Message, error) {
	query := `SELECT id, message, processed, date_create, date_processed_start, date_processed_end FROM messages ORDER BY date_create`

	rows, err := m.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("unable to query messages: %w", err)
	}
	defer rows.Close()

	// messages := []domain.Message{}
	// for rows.Next() {
	// 	message := domain.Message{}
	// 	err := rows.Scan(&message.Id, &message.Message, &message.Processed, &message.DateCreate, &message.DateProcessed)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("unable to scan row: %w", err)
	// 	}
	// 	messages = append(messages, message)
	// }
	// return messages, nil

	return pgx.CollectRows(rows, pgx.RowToStructByName[domain.Message])
}

func (m *Message) Process(ctx context.Context, id domain.Id, dateProcessedStart time.Time) error {
	//query
	query := `UPDATE messages SET processed = true, date_processed_start = @dateProcessedStart, date_processed_end = now() WHERE id = @id`

	//data
	args := pgx.NamedArgs{
		"id":                 id,
		"dateProcessedStart": dateProcessedStart,
	}

	//db
	_, err := m.db.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("unable to update row: %w", err)
	}

	return nil
}
