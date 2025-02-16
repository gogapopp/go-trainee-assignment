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

func TestAuthUser(t *testing.T) {
	tests := []struct {
		name        string
		input       models.AuthRequest
		mockSetup   func(*mocks.Repo)
		wantErr     bool
		errContains string
	}{
		{
			name: "successful authentication",
			input: models.AuthRequest{
				Username: "validuser",
				Password: "validpass",
			},
			mockSetup: func(m *mocks.Repo) {
				m.On("AuthUser", mock.Anything, mock.MatchedBy(func(req models.AuthRequest) bool {
					return req.Username == "validuser" && req.Password == "validpass"
				})).Return(1, nil)
			},
		},
		{
			name: "missing username",
			input: models.AuthRequest{
				Password: "validpass",
			},
			mockSetup:   func(m *mocks.Repo) {},
			wantErr:     true,
			errContains: "Field validation for 'Username' failed on the 'required' tag",
		},
		{
			name: "invalid credentials",
			input: models.AuthRequest{
				Username: "invalid",
				Password: "invalid",
			},
			mockSetup: func(m *mocks.Repo) {
				m.On("AuthUser", mock.Anything, mock.Anything).
					Return(0, errors.New("invalid credentials"))
			},
			wantErr:     true,
			errContains: "invalid credentials",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockRepo := new(mocks.Repo)
			srv := New(mockRepo, "test-secret")

			tt.mockSetup(mockRepo)

			token, err := srv.AuthUser(context.Background(), tt.input)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
				return
			}

			require.NoError(t, err)
			assert.NotEmpty(t, token)
			mockRepo.AssertExpectations(t)
		})
	}
}
