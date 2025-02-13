package config

import (
	"fmt"
	"os"
)

const (
	dbHost     = "DATABASE_HOST"
	dbPort     = "DATABASE_PORT"
	dbUser     = "DATABASE_USER"
	dbPassword = "DATABASE_PASSWORD"
	dbName     = "DATABASE_NAME"
)

type pgConfig struct {
	DSN string
}

func newPGConfig() *pgConfig {
	var dsn string

	dsn = fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv(dbUser),
		os.Getenv(dbPassword),
		os.Getenv(dbHost),
		os.Getenv(dbPort),
		os.Getenv(dbName),
	)

	return &pgConfig{
		DSN: dsn,
	}
}
