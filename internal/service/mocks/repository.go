package mocks

import (
	"context"

	"github.com/gogapopp/go-trainee-assignment/internal/models"
	"github.com/stretchr/testify/mock"
)

type Repo struct {
	mock.Mock
}

func (m *Repo) AuthUser(ctx context.Context, user models.AuthRequest) (int, error) {
	args := m.Called(ctx, user)
	return args.Int(0), args.Error(1)
}

func (m *Repo) BuyItem(ctx context.Context, userID int, itemName string) error {
	args := m.Called(ctx, userID, itemName)
	return args.Error(0)
}

func (m *Repo) GetUserInfo(ctx context.Context, userID int) (models.InfoResponse, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(models.InfoResponse), args.Error(1)
}

func (m *Repo) SendCoins(ctx context.Context, senderID int, req models.SendCoinRequest) error {
	args := m.Called(ctx, senderID, req)
	return args.Error(0)
}
