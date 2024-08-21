package postgres

import (
	"context"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/probuborka/messaggio/internal/domain"
)

func New(ctx context.Context, dbCfg domain.DBConfig) (*pgxpool.Pool, error) {
	var (
		db     *pgxpool.Pool
		pgOnce sync.Once
		err    error
	)

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		dbCfg.User,
		dbCfg.Password,
		dbCfg.Host,
		dbCfg.Port,
		dbCfg.DB,
	)

	pgOnce.Do(func() {
		db, err = pgxpool.New(ctx, dbURL)
		if err != nil {
			return //fmt.Errorf("unable to create connection pool: %w", err)
		}
	})

	return db, err
}
