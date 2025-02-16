package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestLogger(t *testing.T) {
	logger, err := New()
	assert.NoError(t, err)
	assert.NotNil(t, logger)
	assert.IsType(t, &zap.SugaredLogger{}, logger)
}
