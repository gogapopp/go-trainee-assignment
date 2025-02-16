package service

import (
	"context"
	"errors"
	"testing"

	"github.com/gogapopp/go-trainee-assignment/internal/models"
	"github.com/gogapopp/go-trainee-assignment/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetUserInfo(t *testing.T) {
	mockInfo := models.InfoResponse{
		Coins: 1000,
		Inventory: []models.InventoryItem{
			{Type: "t-shirt", Quantity: 2},
		},
		CoinHistory: models.CoinHistory{
			Received: []models.ReceivedTransaction{
				{FromUser: "admin", Amount: 1000},
			},
		},
	}

	tests := []struct {
		name        string
		userID      int
		mockSetup   func(*mocks.Repo)
		want        models.InfoResponse
		wantErr     bool
		errContains string
	}{
		{
			name:   "successful info retrieval",
			userID: 1,
			mockSetup: func(m *mocks.Repo) {
				m.On("GetUserInfo", mock.Anything, 1).
					Return(mockInfo, nil)
			},
			want: mockInfo,
		},
		{
			name:   "user not found",
			userID: 999,
			mockSetup: func(m *mocks.Repo) {
				m.On("GetUserInfo", mock.Anything, 999).
					Return(models.InfoResponse{}, errors.New("user not found"))
			},
			wantErr:     true,
			errContains: "user not found",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockRepo := new(mocks.Repo)
			srv := New(mockRepo, "")

			tt.mockSetup(mockRepo)

			result, err := srv.GetUserInfo(context.Background(), tt.userID)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, result)
			mockRepo.AssertExpectations(t)
		})
	}
}
