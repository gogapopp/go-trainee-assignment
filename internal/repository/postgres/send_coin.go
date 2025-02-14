package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/gogapopp/go-trainee-assignment/internal/models"
	"github.com/gogapopp/go-trainee-assignment/internal/repository"
	"github.com/jackc/pgx/v5"
)

func (s *storage) SendCoins(ctx context.Context, senderID int, req models.SendCoinRequest) error {
	const op = "internal.repository.postgres.send_coin.SendCoins"

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback(ctx)

	// check if recipient exists
	var recipientID int
	err = tx.QueryRow(ctx, "SELECT id FROM users WHERE username = $1", req.ToUser).Scan(&recipientID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("%s: %w", op, repository.ErrUserNotFound)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	if recipientID == senderID {
		return fmt.Errorf("%s: %w", op, repository.ErrSameUser)
	}

	// check user balance
	var balance int
	err = tx.QueryRow(ctx, "SELECT coins FROM users WHERE id = $1 FOR UPDATE", senderID).Scan(&balance)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if balance < req.Amount {
		return fmt.Errorf("%s: %w", op, repository.ErrInsufficientFunds)
	}

	// minus coins from sender
	_, err = tx.Exec(ctx, "UPDATE users SET coins = coins - $1 WHERE id = $2", req.Amount, senderID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// add coins to recipient
	_, err = tx.Exec(ctx, "UPDATE users SET coins = coins + $1 WHERE id = $2", req.Amount, recipientID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = tx.Exec(ctx, `INSERT INTO user_coin_history (from_user_id, to_user_id, amount) VALUES ($1, $2, $3)`,
		senderID, recipientID, req.Amount)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return tx.Commit(ctx)
}
