package service

import (
	"context"
	"fmt"

	"github.com/gogapopp/go-trainee-assignment/internal/libs/jwt"
	"github.com/gogapopp/go-trainee-assignment/internal/models"
)

func (s *Service) AuthUser(ctx context.Context, req models.AuthRequest) (string, error) {
	const op = "internal.service.auth.AuthUser"

	if err := s.validator.Struct(req); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	userID, err := s.repo.AuthUser(ctx, req)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	token, err := jwt.GenerateJWTToken(s.jwtSecret, userID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}
