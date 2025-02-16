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

func TestSendCoins(t *testing.T) {
	validRequest := models.SendCoinRequest{
		ToUser: "recipient",
		Amount: 100,
	}

	tests := []struct {
		name        string
		senderID    int
		request     models.SendCoinRequest
		mockSetup   func(*mocks.Repo)
		wantErr     bool
		errContains string
	}{
		{
			name:     "successful transfer",
			senderID: 1,
			request:  validRequest,
			mockSetup: func(m *mocks.Repo) {
				m.On("SendCoins", mock.Anything, 1, validRequest).
					Return(nil)
			},
		},
		{
			name:     "missing to_user",
			senderID: 1,
			request: models.SendCoinRequest{
				Amount: 100,
			},
			mockSetup:   func(m *mocks.Repo) {},
			wantErr:     true,
			errContains: "Field validation for 'ToUser' failed on the 'required' tag",
		},
		{
			name:     "insufficient funds",
			senderID: 2,
			request:  validRequest,
			mockSetup: func(m *mocks.Repo) {
				m.On("SendCoins", mock.Anything, 2, validRequest).
					Return(errors.New("insufficient funds"))
			},
			wantErr:     true,
			errContains: "insufficient funds",
		},
		{
			name:     "send to self",
			senderID: 1,
			request: models.SendCoinRequest{
				ToUser: "sender",
				Amount: 100,
			},
			mockSetup: func(m *mocks.Repo) {
				m.On("SendCoins", mock.Anything, 1, mock.Anything).
					Return(errors.New("cant send to yourself"))
			},
			wantErr:     true,
			errContains: "cant send to yourself",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockRepo := new(mocks.Repo)
			srv := New(mockRepo, "")

			tt.mockSetup(mockRepo)

			err := srv.SendCoins(context.Background(), tt.senderID, tt.request)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
			} else {
				require.NoError(t, err)
				mockRepo.AssertExpectations(t)
			}
		})
	}
}
