package service

import (
	"context"

	"github.com/go-playground/validator"
	"github.com/gogapopp/go-trainee-assignment/internal/models"
	"go.uber.org/zap"
)

type repo interface {
	AuthUser(ctx context.Context, user models.AuthRequest) (int, error)
}

type Service struct {
	logger    *zap.SugaredLogger
	repo      repo
	validator *validator.Validate
	jwtSecret string
}

func New(logger *zap.SugaredLogger, repo repo, jwtSecret string) *Service {
	return &Service{
		logger:    logger,
		repo:      repo,
		validator: validator.New(),
		jwtSecret: jwtSecret,
	}
}
