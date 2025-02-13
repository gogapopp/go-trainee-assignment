package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type storage struct {
	db *pgx.Conn
}

func New(dsn string) (*storage, error) {
	const op = "internal.repository.postgres.postgres.New"

	ctx := context.Background()

	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = conn.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &storage{
		db: conn,
	}, nil
}

func (s *storage) Close(ctx context.Context) {
	s.db.Close(ctx)
}
