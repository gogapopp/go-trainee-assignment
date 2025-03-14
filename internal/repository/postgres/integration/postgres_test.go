// Package integration for postgres repository
package intergration

import (
	"fmt"
	"os"
	"testing"

	"github.com/gogapopp/go-trainee-assignment/internal/libs/config"
)

const envPath = "../../../../.env"

func getDSN(t *testing.T) string {
	user := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASSWORD")
	dbName := os.Getenv("DATABASE_NAME")
	host := os.Getenv("DATABASE_HOST")
	port := os.Getenv("DATABASE_PORT")

	if user != "" && password != "" && dbName != "" && host != "" && port != "" {
		return fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			user, password, host, port, dbName,
		)
	}

	cfg, err := config.New(envPath)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	return cfg.PGConfig.DSN
}
