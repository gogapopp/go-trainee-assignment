// package with integration tests for postgres repository
package intergration

import (
	"testing"

	"github.com/gogapopp/go-trainee-assignment/internal/libs/config"
)

func getDSN(t *testing.T) string {
	cfg, err := config.New("../../../../.env")
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}
	return cfg.PGConfig.DSN
}
