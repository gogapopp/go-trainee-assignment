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
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv(dbHost),
		os.Getenv(dbPort),
		os.Getenv(dbUser),
		os.Getenv(dbPassword),
		os.Getenv(dbName),
	)

	return &pgConfig{
		DSN: dsn,
	}
}
