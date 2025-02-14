package service

import (
	"context"
	"fmt"

	"github.com/gogapopp/go-trainee-assignment/internal/models"
)

func (s *Service) GetUserInfo(ctx context.Context, userID int) (models.InfoResponse, error) {
	const op = "internal.service.info.GetUserInfo"

	info, err := s.repo.GetUserInfo(ctx, userID)
	if err != nil {
		return models.InfoResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	return info, nil
}
