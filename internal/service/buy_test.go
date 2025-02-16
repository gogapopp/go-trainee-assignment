package service

import (
	"context"
	"errors"
	"testing"

	"github.com/gogapopp/go-trainee-assignment/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestBuyItem(t *testing.T) {
	tests := []struct {
		name        string
		userID      int
		itemName    string
		mockSetup   func(*mocks.Repo)
		wantErr     bool
		errContains string
	}{
		{
			name:     "successful purchase",
			userID:   1,
			itemName: "t-shirt",
			mockSetup: func(m *mocks.Repo) {
				m.On("BuyItem", mock.Anything, 1, "t-shirt").
					Return(nil)
			},
		},
		{
			name:     "item not found",
			userID:   1,
			itemName: "invalid-item",
			mockSetup: func(m *mocks.Repo) {
				m.On("BuyItem", mock.Anything, 1, "invalid-item").
					Return(errors.New("item not found"))
			},
			wantErr:     true,
			errContains: "item not found",
		},
		{
			name:     "insufficient funds",
			userID:   2,
			itemName: "cup",
			mockSetup: func(m *mocks.Repo) {
				m.On("BuyItem", mock.Anything, 2, "cup").
					Return(errors.New("insufficient funds"))
			},
			wantErr:     true,
			errContains: "insufficient funds",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockRepo := new(mocks.Repo)
			srv := New(mockRepo, "")

			tt.mockSetup(mockRepo)

			err := srv.BuyItem(context.Background(), tt.userID, tt.itemName)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
				return
			}

			require.NoError(t, err)
			mockRepo.AssertExpectations(t)
		})
	}
}
