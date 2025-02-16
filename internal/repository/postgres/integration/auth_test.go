// package with integration tests for postgres repository
package intergration

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/gogapopp/go-trainee-assignment/internal/models"
	"github.com/gogapopp/go-trainee-assignment/internal/repository"
	"github.com/gogapopp/go-trainee-assignment/internal/repository/postgres"
)

func TestAuthUser(t *testing.T) {
	testCases := []struct {
		name     string
		username string
		password string
		testFunc func(t *testing.T, s *postgres.Storage, username, password string)
	}{
		{
			name:     "Registration",
			username: fmt.Sprintf("test_auth_reg_%d", time.Now().UnixNano()),
			password: "secret",
			testFunc: func(t *testing.T, s *postgres.Storage, username, password string) {
				ctx := context.Background()
				userID, err := s.AuthUser(ctx, models.AuthRequest{
					Username: username,
					Password: password,
				})
				if err != nil {
					t.Fatalf("Registration failed: %v", err)
				}
				if userID <= 0 {
					t.Fatalf("Invalid user id: %d", userID)
				}
			},
		},
		{
			name:     "LoginCorrect",
			username: fmt.Sprintf("test_auth_login_%d", time.Now().UnixNano()),
			password: "secret",
			testFunc: func(t *testing.T, s *postgres.Storage, username, password string) {
				ctx := context.Background()
				userID1, err := s.AuthUser(ctx, models.AuthRequest{
					Username: username,
					Password: password,
				})
				if err != nil {
					t.Fatalf("Registration failed: %v", err)
				}
				if userID1 <= 0 {
					t.Fatalf("Invalid user id: %d", userID1)
				}
				userID2, err := s.AuthUser(ctx, models.AuthRequest{
					Username: username,
					Password: password,
				})
				if err != nil {
					t.Fatalf("Login failed: %v", err)
				}
				if userID1 != userID2 {
					t.Fatalf("Expected same user id on login, got %d and %d", userID1, userID2)
				}
			},
		},
		{
			name:     "LoginWrongPassword",
			username: fmt.Sprintf("test_auth_wrong_%d", time.Now().UnixNano()),
			password: "secret",
			testFunc: func(t *testing.T, s *postgres.Storage, username, password string) {
				ctx := context.Background()
				userID, err := s.AuthUser(ctx, models.AuthRequest{
					Username: username,
					Password: password,
				})
				if err != nil {
					t.Fatalf("Registration failed: %v", err)
				}
				if userID <= 0 {
					t.Fatalf("Invalid user id: %d", userID)
				}
				_, err = s.AuthUser(ctx, models.AuthRequest{
					Username: username,
					Password: "wrong",
				})
				if err == nil {
					t.Fatal("Expected error for wrong password, got nil")
				}
				if !errors.Is(err, repository.ErrInvalidCredentials) {
					t.Fatalf("Expected ErrInvalidCredentials, got: %v", err)
				}
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dsn := getDSN(t)
			s, err := postgres.New(dsn, "secret")
			if err != nil {
				t.Fatalf("failed to create storage: %v", err)
			}
			defer s.DB.Close(context.Background())
			tc.testFunc(t, s, tc.username, tc.password)
		})
	}
}
