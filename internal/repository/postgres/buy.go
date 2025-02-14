package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/gogapopp/go-trainee-assignment/internal/repository"
	"github.com/jackc/pgx/v5"
)

func (s *storage) BuyItem(ctx context.Context, userID int, itemName string) error {
	const op = "internal.repository.postgres.buy.BuyItem"

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback(ctx)

	// check item price
	var price int
	err = tx.QueryRow(ctx, "SELECT price FROM items WHERE name = $1", itemName).Scan(&price)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("%s: %w", op, repository.ErrItemNotFound)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	// check user balance
	var balance int
	err = tx.QueryRow(ctx, "SELECT coins FROM users WHERE id = $1 FOR UPDATE", userID).Scan(&balance)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if balance < price {
		return fmt.Errorf("%s: %w", op, repository.ErrInsufficientFunds)
	}

	// minus coins from user
	_, err = tx.Exec(ctx, "UPDATE users SET coins = coins - $1 WHERE id = $2", price, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// add item to user inventory
	_, err = tx.Exec(ctx,
		`INSERT INTO user_inventory (user_id, item_id, quantity) VALUES ($1, (SELECT id FROM items WHERE name = $2), 1)
        ON CONFLICT (user_id, item_id) DO UPDATE SET quantity = user_inventory.quantity + 1`,
		userID, itemName,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return tx.Commit(ctx)
}
