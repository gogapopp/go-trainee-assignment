// Package integration with integration tests for handlers
package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gogapopp/go-trainee-assignment/internal/http-server/handlers"
	"github.com/gogapopp/go-trainee-assignment/internal/libs/config"
	"github.com/gogapopp/go-trainee-assignment/internal/libs/logger"
	"github.com/gogapopp/go-trainee-assignment/internal/models"
	"github.com/gogapopp/go-trainee-assignment/internal/repository/postgres"
	"github.com/gogapopp/go-trainee-assignment/internal/service"
	"github.com/stretchr/testify/assert"
)

func setupAuthHandler(t *testing.T) (http.HandlerFunc, func()) {
	dsn := getDSN(t)
	cfg, err := config.New("../../../../.env")
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	repo, err := postgres.New(dsn, "secret")
	if err != nil {
		t.Fatalf("failed to create repository: %v", err)
	}

	svc := service.New(repo, cfg.JWTSecret)
	logger, _ := logger.New()
	handler := handlers.AuthHandler(logger, svc)

	cleanup := func() {
		repo.DB.Close(context.Background())
	}

	return handler, cleanup
}

func TestAuthHandler(t *testing.T) {
	handler, cleanup := setupAuthHandler(t)
	defer cleanup()

	tests := []struct {
		name           string
		requestBody    models.AuthRequest
		expectedStatus int
		expectedErrors string
	}{
		{
			name: "Successful authentication",
			requestBody: models.AuthRequest{
				Username: "testuser",
				Password: "password123",
			},
			expectedStatus: http.StatusOK,
			expectedErrors: "",
		},
		{
			name: "Invalid credentials",
			requestBody: models.AuthRequest{
				Username: "testuser",
				Password: "wrongpassword",
			},
			expectedStatus: http.StatusUnauthorized,
			expectedErrors: "invalid credentials",
		},
		{
			name: "Missing username",
			requestBody: models.AuthRequest{
				Username: "",
				Password: "password123",
			},
			expectedStatus: http.StatusBadRequest,
			expectedErrors: "validation error",
		},
		{
			name: "Missing password",
			requestBody: models.AuthRequest{
				Username: "testuser",
				Password: "",
			},
			expectedStatus: http.StatusBadRequest,
			expectedErrors: "validation error",
		},
		{
			name:           "Empty request body",
			requestBody:    models.AuthRequest{},
			expectedStatus: http.StatusBadRequest,
			expectedErrors: "validation error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, err := json.Marshal(tt.requestBody)
			if err != nil {
				t.Fatalf("failed to marshal request body: %v", err)
			}

			r := httptest.NewRequest(http.MethodPost, "/auth", bytes.NewBuffer(reqBody))
			r.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()

			handler.ServeHTTP(w, r)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.expectedStatus == http.StatusOK {
				bodyBytes, err := io.ReadAll(resp.Body)
				if err != nil {
					t.Fatalf("failed to read response body: %v", err)
				}
				var responseBody map[string]interface{}
				err = json.Unmarshal(bodyBytes, &responseBody)
				if err != nil {
					t.Fatalf("failed to unmarshal response body: %v", err)
				}
				token, ok := responseBody["token"].(string)
				assert.True(t, ok, "expected token in response")
				assert.NotEmpty(t, token, "token should not be empty")
			} else {
				bodyBytes, err := io.ReadAll(resp.Body)
				if err != nil {
					t.Fatalf("failed to read response body: %v", err)
				}
				var responseBody map[string]string
				err = json.Unmarshal(bodyBytes, &responseBody)
				if err != nil {
					t.Fatalf("failed to unmarshal response body: %v", err)
				}
				errorsMessage, ok := responseBody["errors"]
				assert.True(t, ok, "expected 'errors' key in response")
				assert.Equal(t, tt.expectedErrors, errorsMessage)
			}
		})
	}
}
