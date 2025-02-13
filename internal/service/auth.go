package service

import (
	"context"
	"fmt"

	"github.com/gogapopp/go-trainee-assignment/internal/libs/jwt"
	"github.com/gogapopp/go-trainee-assignment/internal/models"
)

func (a *Service) AuthUser(ctx context.Context, req models.AuthRequest) (string, error) {
	const op = "internal.service.auth.AuthUser"

	if err := a.validator.Struct(req); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	userID, err := a.repo.AuthUser(ctx, req)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	token, err := jwt.GenerateJWTToken(a.jwtSecret, userID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}
