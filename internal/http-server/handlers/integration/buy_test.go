// Package integration with integration tests for handlers
package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/gogapopp/go-trainee-assignment/internal/http-server/handlers"
	"github.com/gogapopp/go-trainee-assignment/internal/http-server/middlewares"
	"github.com/gogapopp/go-trainee-assignment/internal/libs/config"
	"github.com/gogapopp/go-trainee-assignment/internal/models"
	"github.com/gogapopp/go-trainee-assignment/internal/repository/postgres"
	"github.com/gogapopp/go-trainee-assignment/internal/service"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func setupBuyHandler(t *testing.T) (http.HandlerFunc, func()) {
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
	logger := zap.NewNop().Sugar()

	r := chi.NewRouter()
	r.Use(middlewares.AuthMiddleware(cfg.JWTSecret))
	r.Get("/api/buy/{item}", handlers.BuyItemHandler(logger, svc))

	cleanup := func() {
		repo.DB.Close(context.Background())
	}
	return r.ServeHTTP, cleanup
}

func TestBuyItemHandler(t *testing.T) {
	buyHandler, buyCleanup := setupBuyHandler(t)
	authHandler, authCleanup := setupAuthHandler(t)
	defer buyCleanup()
	defer authCleanup()

	tests := []struct {
		name         string
		itemName     string
		initialCoins int
		tokenValid   bool
		wantStatus   int
		wantError    string
	}{
		{
			name:         "successful purchase",
			itemName:     "pen",
			initialCoins: 20,
			tokenValid:   true,
			wantStatus:   http.StatusOK,
		},
		{
			name:         "item not found",
			itemName:     "invalid-item",
			initialCoins: 100,
			tokenValid:   true,
			wantStatus:   http.StatusNotFound,
			wantError:    "item not found",
		},
		{
			name:         "insufficient funds",
			itemName:     "hoody",
			initialCoins: 200,
			tokenValid:   true,
			wantStatus:   http.StatusBadRequest,
			wantError:    "insufficient funds",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userReq := models.AuthRequest{
				Username: "testuser_" + tt.name,
				Password: "password123",
			}
			token := createTestUser(authHandler, userReq)

			setInitialCoins(t, userReq.Username, tt.initialCoins)

			url := fmt.Sprintf("/api/buy/%s", tt.itemName)
			req := httptest.NewRequest("GET", url, nil)

			if tt.tokenValid {
				req.Header.Set("Authorization", "Bearer "+token)
			} else {
				req.Header.Set("Authorization", "Bearer invalid_token")
			}

			w := httptest.NewRecorder()
			buyHandler.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.wantError != "" {
				var response map[string]string
				err := json.NewDecoder(w.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.wantError, response["errors"])
			}

			if w.Code == http.StatusOK {
				coins, err := getCurrentCoins(t, userReq.Username)
				assert.NoError(t, err)
				itemPrice := getItemPrice(tt.itemName)
				assert.Equal(t, tt.initialCoins-itemPrice, coins)

				hasItem, err := checkInventory(t, userReq.Username, tt.itemName)
				assert.NoError(t, err)
				assert.True(t, hasItem)
			}
		})
	}
}

func setInitialCoins(t *testing.T, username string, coins int) {
	dsn := getDSN(t)
	repo, err := postgres.New(dsn, "secret")
	if err != nil {
		t.Fatalf("failed to create repo: %v", err)
	}
	defer repo.DB.Close(context.Background())

	_, err = repo.DB.Exec(context.Background(),
		"UPDATE users SET coins = $1 WHERE username = $2",
		coins, username,
	)
	if err != nil {
		t.Fatalf("failed to set coins: %v", err)
	}
}

func getCurrentCoins(t *testing.T, username string) (int, error) {
	dsn := getDSN(t)
	repo, err := postgres.New(dsn, "secret")
	if err != nil {
		t.Fatalf("failed to create repo: %v", err)
	}
	defer repo.DB.Close(context.Background())

	var coins int
	err = repo.DB.QueryRow(context.Background(),
		"SELECT coins FROM users WHERE username = $1",
		username,
	).Scan(&coins)

	return coins, err
}

func checkInventory(t *testing.T, username, itemName string) (bool, error) {
	dsn := getDSN(t)
	repo, err := postgres.New(dsn, "secret")
	if err != nil {
		t.Fatalf("failed to create repo: %v", err)
	}
	defer repo.DB.Close(context.Background())

	var count int
	err = repo.DB.QueryRow(context.Background(),
		`SELECT COUNT(*) 
		FROM user_inventory ui
		JOIN users u ON ui.user_id = u.id
		JOIN items i ON ui.item_id = i.id
		WHERE u.username = $1 AND i.name = $2`,
		username, itemName,
	).Scan(&count)

	return count > 0, err
}

func createTestUser(handler http.Handler, user models.AuthRequest) string {
	body, _ := json.Marshal(user)
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(body))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	var response map[string]string
	json.NewDecoder(w.Body).Decode(&response)
	return response["token"]
}

func getCurrentBalance(t *testing.T, username string) (int, error) {
	dsn := getDSN(t)
	repo, err := postgres.New(dsn, "secret")
	if err != nil {
		t.Fatalf("failed to create repo: %v", err)
	}
	defer repo.DB.Close(context.Background())

	var balance int
	err = repo.DB.QueryRow(context.Background(),
		"SELECT balance FROM users WHERE username = $1",
		username,
	).Scan(&balance)

	return balance, err
}

func getItemPrice(itemName string) int {
	prices := map[string]int{
		"t-shirt":    80,
		"cup":        20,
		"book":       50,
		"pen":        10,
		"powerbank":  200,
		"hoody":      300,
		"umbrella":   200,
		"socks":      10,
		"wallet":     50,
		"pink-hoody": 500,
	}
	return prices[itemName]
}
