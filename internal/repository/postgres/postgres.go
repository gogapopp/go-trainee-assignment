package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type Storage struct {
	DB         *pgx.Conn
	passSecret string
}

func New(dsn string, passSecret string) (*Storage, error) {
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

	return &Storage{
		DB:         conn,
		passSecret: passSecret,
	}, nil
}
