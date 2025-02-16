// package with integration tests for postgres repository
package intergration

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/gogapopp/go-trainee-assignment/internal/libs/config"
)

func getDSN(t *testing.T) string {
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")

	if user == "" || password == "" || dbName == "" || host == "" || port == "" {
		envPath := "../../../../.env"

		if ws := os.Getenv("GITHUB_WORKSPACE"); ws != "" {
			envPath = filepath.Join(ws, ".env")
		}

		cfg, err := config.New(envPath)
		if err != nil {
			t.Fatalf("failed to load config: %v", err)
		}

		return cfg.PGConfig.DSN
	}
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv(user),
		os.Getenv(password),
		os.Getenv(host),
		os.Getenv(port),
		os.Getenv(dbName),
	)
}
