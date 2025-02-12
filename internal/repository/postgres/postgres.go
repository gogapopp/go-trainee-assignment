package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func New(dsn string) (*pgx.Conn, error) {
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

	return conn, nil
}
