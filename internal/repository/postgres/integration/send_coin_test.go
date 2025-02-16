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

func TestSendCoins(t *testing.T) {
	type testCase struct {
		name             string
		senderBalance    int
		recipientBalance int
		transferAmount   int
		selfTransfer     bool
		setupRecipient   bool
		expectedErr      error
		postCheckFunc    func(t *testing.T, s *postgres.Storage, senderID int, recipientName string)
	}

	testCases := []testCase{
		{
			name:             "Successful Transfer",
			senderBalance:    200,
			recipientBalance: 50,
			transferAmount:   100,
			selfTransfer:     false,
			setupRecipient:   true,
			expectedErr:      nil,
			postCheckFunc: func(t *testing.T, s *postgres.Storage, senderID int, recipientName string) {
				ctx := context.Background()
				var senderCoins int
				err := s.DB.QueryRow(ctx, "SELECT coins FROM users WHERE id = $1", senderID).Scan(&senderCoins)
				if err != nil {
					t.Fatalf("failed to query sender coins: %v", err)
				}
				if senderCoins != 100 {
					t.Fatalf("expected sender coins to be 100, got %d", senderCoins)
				}
				var recipientCoins int
				err = s.DB.QueryRow(ctx, "SELECT coins FROM users WHERE username = $1", recipientName).Scan(&recipientCoins)
				if err != nil {
					t.Fatalf("failed to query recipient coins: %v", err)
				}
				if recipientCoins != 150 {
					t.Fatalf("expected recipient coins to be 150, got %d", recipientCoins)
				}
			},
		},
		{
			name:           "Self Transfer",
			senderBalance:  200,
			transferAmount: 10,
			selfTransfer:   true,
			setupRecipient: true,
			expectedErr:    repository.ErrSameUser,
		},
		{
			name:             "Insufficient Funds",
			senderBalance:    100,
			recipientBalance: 50,
			transferAmount:   200,
			selfTransfer:     false,
			setupRecipient:   true,
			expectedErr:      repository.ErrInsufficientFunds,
		},
		{
			name:           "Non-existent Recipient",
			senderBalance:  200,
			transferAmount: 10,
			selfTransfer:   false,
			setupRecipient: false,
			expectedErr:    repository.ErrUserNotFound,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			dsn := getDSN(t)
			s, err := postgres.New(dsn, "secret")
			if err != nil {
				t.Fatalf("failed to create storage: %v", err)
			}
			defer s.DB.Close(ctx)

			senderName := fmt.Sprintf("sender_%d", time.Now().UnixNano())
			password := "pass123"
			senderID, err := s.AuthUser(ctx, models.AuthRequest{Username: senderName, Password: password})
			if err != nil {
				t.Fatalf("failed to create sender: %v", err)
			}
			_, err = s.DB.Exec(ctx, "UPDATE users SET coins = $1 WHERE id = $2", tc.senderBalance, senderID)
			if err != nil {
				t.Fatalf("failed to update sender coins: %v", err)
			}

			var recipientName string
			if tc.selfTransfer {
				recipientName = senderName
			} else if tc.setupRecipient {
				recipientName = fmt.Sprintf("recipient_%d", time.Now().UnixNano())
				_, err = s.AuthUser(ctx, models.AuthRequest{Username: recipientName, Password: password})
				if err != nil {
					t.Fatalf("failed to create recipient: %v", err)
				}
				_, err = s.DB.Exec(ctx, "UPDATE users SET coins = $1 WHERE username = $2", tc.recipientBalance, recipientName)
				if err != nil {
					t.Fatalf("failed to update recipient coins: %v", err)
				}
			} else {
				recipientName = fmt.Sprintf("nonexistent_%d", time.Now().UnixNano())
			}

			req := models.SendCoinRequest{
				ToUser: recipientName,
				Amount: tc.transferAmount,
			}
			err = s.SendCoins(ctx, senderID, req)
			if tc.expectedErr != nil {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if !errors.Is(err, tc.expectedErr) {
					t.Fatalf("expected error %v, got %v", tc.expectedErr, err)
				}
			} else {
				if err != nil {
					t.Fatalf("SendCoins failed: %v", err)
				}
				if tc.postCheckFunc != nil {
					tc.postCheckFunc(t, s, senderID, recipientName)
				}
			}
		})
	}
}
