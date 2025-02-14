package service

import (
	"context"
	"fmt"
)

func (s *Service) BuyItem(ctx context.Context, userID int, itemName string) error {
	const op = "internal.service.buy.BuyItem"

	if err := s.repo.BuyItem(ctx, userID, itemName); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
