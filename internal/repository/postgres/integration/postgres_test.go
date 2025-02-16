// package with integration tests for postgres repository
package intergration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gogapopp/go-trainee-assignment/internal/libs/config"
)

func getDSN(t *testing.T) string {
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
