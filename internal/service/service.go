package service

import (
	"context"

	"github.com/go-playground/validator"
	"github.com/gogapopp/go-trainee-assignment/internal/models"
)

type repo interface {
	AuthUser(ctx context.Context, user models.AuthRequest) (int, error)
	BuyItem(ctx context.Context, userID int, itemName string) error
	GetUserInfo(ctx context.Context, userID int) (models.InfoResponse, error)
	SendCoins(ctx context.Context, senderID int, req models.SendCoinRequest) error
}

type Service struct {
	repo      repo
	validator *validator.Validate
	jwtSecret string
}

func New(repo repo, jwtSecret string) *Service {
	return &Service{
		repo:      repo,
		validator: validator.New(),
		jwtSecret: jwtSecret,
	}
}
