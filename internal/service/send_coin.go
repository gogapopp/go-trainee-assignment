package service

import (
	"context"
	"fmt"

	"github.com/gogapopp/go-trainee-assignment/internal/models"
)

func (s *Service) SendCoins(ctx context.Context, senderID int, req models.SendCoinRequest) error {
	const op = "internal.service.send_coin.SendCoins"

	if err := s.validator.Struct(req); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := s.repo.SendCoins(ctx, senderID, req); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
