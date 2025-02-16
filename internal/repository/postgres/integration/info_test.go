// Package integration tests for postgres repository
package intergration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gogapopp/go-trainee-assignment/internal/models"
	"github.com/gogapopp/go-trainee-assignment/internal/repository/postgres"
)

func TestGetUserInfo(t *testing.T) {
	type testCase struct {
		name                string
		userAInitialBalance int
		userBInitialBalance int
		purchasePrice       int
		transferAmount      int
		expectedUserACoins  int
		expectedUserBCoins  int
		postCheckFunc       func(t *testing.T, s *postgres.Storage, userAID int, userBName, itemName string, transferAmount int)
	}

	testCases := []testCase{
		{
			name:                "Info After Purchase and Transfer",
			userAInitialBalance: 300,
			userBInitialBalance: 100,
			purchasePrice:       50,
			transferAmount:      70,
			expectedUserACoins:  180,
			expectedUserBCoins:  170,
			postCheckFunc: func(t *testing.T, s *postgres.Storage, userAID int, userBName, itemName string, transferAmount int) {
				ctx := context.Background()

				var quantity int
				err := s.DB.QueryRow(ctx, `
					SELECT ui.quantity
					FROM user_inventory ui
					JOIN items i ON ui.item_id = i.id
					WHERE ui.user_id = $1 AND i.name = $2`, userAID, itemName).Scan(&quantity)
				if err != nil {
					t.Fatalf("failed to query userA inventory: %v", err)
				}
				if quantity != 1 {
					t.Fatalf("expected quantity 1 for item %s, got %d", itemName, quantity)
				}
				rows, err := s.DB.Query(ctx, `
					SELECT uch.amount FROM user_coin_history uch
					WHERE uch.from_user_id = $1 AND uch.amount = $2`, userAID, transferAmount)
				if err != nil {
					t.Fatalf("failed to query userA coin history: %v", err)
				}
				defer rows.Close()
				if !rows.Next() {
					t.Fatalf("expected a sent transaction of amount %d for userA", transferAmount)
				}
			},
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

			userAName := fmt.Sprintf("userA_%d", time.Now().UnixNano())
			userBName := fmt.Sprintf("userB_%d", time.Now().UnixNano())
			password := "password"
			userAID, err := s.AuthUser(ctx, models.AuthRequest{Username: userAName, Password: password})
			if err != nil {
				t.Fatalf("failed to create userA: %v", err)
			}
			userBID, err := s.AuthUser(ctx, models.AuthRequest{Username: userBName, Password: password})
			if err != nil {
				t.Fatalf("failed to create userB: %v", err)
			}

			_, err = s.DB.Exec(ctx, "UPDATE users SET coins = $1 WHERE id = $2", tc.userAInitialBalance, userAID)
			if err != nil {
				t.Fatalf("failed to update userA coins: %v", err)
			}
			_, err = s.DB.Exec(ctx, "UPDATE users SET coins = $1 WHERE id = $2", tc.userBInitialBalance, userBID)
			if err != nil {
				t.Fatalf("failed to update userB coins: %v", err)
			}

			itemName := fmt.Sprintf("Item_%s_%d", tc.name, time.Now().UnixNano())

			_, err = s.DB.Exec(ctx, "INSERT INTO items(name, price) VALUES($1, $2)", itemName, tc.purchasePrice)
			if err != nil {
				t.Fatalf("failed to insert item: %v", err)
			}
			err = s.BuyItem(ctx, userAID, itemName)
			if err != nil {
				t.Fatalf("userA failed to buy item: %v", err)
			}

			err = s.SendCoins(ctx, userAID, models.SendCoinRequest{
				ToUser: userBName,
				Amount: tc.transferAmount,
			})
			if err != nil {
				t.Fatalf("userA failed to send coins to userB: %v", err)
			}

			infoA, err := s.GetUserInfo(ctx, userAID)
			if err != nil {
				t.Fatalf("GetUserInfo for userA failed: %v", err)
			}
			if infoA.Coins != tc.expectedUserACoins {
				t.Fatalf("expected userA coins %d, got %d", tc.expectedUserACoins, infoA.Coins)
			}
			infoB, err := s.GetUserInfo(ctx, userBID)
			if err != nil {
				t.Fatalf("GetUserInfo for userB failed: %v", err)
			}
			if infoB.Coins != tc.expectedUserBCoins {
				t.Fatalf("expected userB coins %d, got %d", tc.expectedUserBCoins, infoB.Coins)
			}

			if tc.postCheckFunc != nil {
				tc.postCheckFunc(t, s, userAID, userBName, itemName, tc.transferAmount)
			}
		})
	}
}
