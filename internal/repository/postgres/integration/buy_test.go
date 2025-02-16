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

func TestBuyItem(t *testing.T) {
	type testCase struct {
		name          string
		setupItem     bool
		price         int
		buyerBalance  int
		callItemName  string
		expectedErr   error
		postCheckFunc func(t *testing.T, s *postgres.Storage, buyerID int, itemName string)
	}

	testCases := []testCase{
		{
			name:         "Successful Purchase",
			setupItem:    true,
			price:        100,
			buyerBalance: 200,
			expectedErr:  nil,
			postCheckFunc: func(t *testing.T, s *postgres.Storage, buyerID int, itemName string) {
				ctx := context.Background()
				var quantity int
				err := s.DB.QueryRow(ctx, `
					SELECT ui.quantity
					FROM user_inventory ui
					JOIN items i ON ui.item_id = i.id
					WHERE ui.user_id = $1 AND i.name = $2`, buyerID, itemName).Scan(&quantity)
				if err != nil {
					t.Fatalf("failed to query user inventory: %v", err)
				}
				if quantity != 1 {
					t.Fatalf("expected quantity 1, got %d", quantity)
				}
			},
		},
		{
			name:         "Insufficient Funds",
			setupItem:    true,
			price:        100,
			buyerBalance: 50,
			expectedErr:  repository.ErrInsufficientFunds,
		},
		{
			name:         "Non-existent Item",
			setupItem:    false,
			buyerBalance: 200,
			callItemName: "NonExistentItem",
			expectedErr:  repository.ErrItemNotFound,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			dsn := getDSN(t)
			s, err := postgres.New(dsn)
			if err != nil {
				t.Fatalf("failed to create storage: %v", err)
			}
			defer s.DB.Close(ctx)

			var itemName string
			if tc.setupItem {

				itemName = fmt.Sprintf("TestItem_%s_%d", tc.name, time.Now().UnixNano())

				if tc.callItemName == "" {
					tc.callItemName = itemName
				}
				_, err = s.DB.Exec(ctx, "INSERT INTO items(name, price) VALUES($1, $2)", itemName, tc.price)
				if err != nil {
					t.Fatalf("failed to insert test item: %v", err)
				}
			} else {
				itemName = tc.callItemName
			}

			username := fmt.Sprintf("buyer_%d", time.Now().UnixNano())
			password := "buyerpass"
			buyerID, err := s.AuthUser(ctx, models.AuthRequest{Username: username, Password: password})
			if err != nil {
				t.Fatalf("failed to create test user: %v", err)
			}
			_, err = s.DB.Exec(ctx, "UPDATE users SET coins = $1 WHERE id = $2", tc.buyerBalance, buyerID)
			if err != nil {
				t.Fatalf("failed to update user coins: %v", err)
			}

			err = s.BuyItem(ctx, buyerID, tc.callItemName)
			if tc.expectedErr != nil {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if !errors.Is(err, tc.expectedErr) {
					t.Fatalf("expected error %v, got %v", tc.expectedErr, err)
				}
			} else {
				if err != nil {
					t.Fatalf("BuyItem failed: %v", err)
				}
				if tc.postCheckFunc != nil {
					tc.postCheckFunc(t, s, buyerID, tc.callItemName)
				}
			}
		})
	}
}
