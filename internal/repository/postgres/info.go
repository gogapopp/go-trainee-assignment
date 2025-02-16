package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/gogapopp/go-trainee-assignment/internal/models"
	"github.com/gogapopp/go-trainee-assignment/internal/repository"
	"github.com/jackc/pgx/v5"
)

func (s *Storage) GetUserInfo(ctx context.Context, userID int) (models.InfoResponse, error) {
	const op = "internal.repository.postgres.info.GetUserInfo"

	var info models.InfoResponse

	err := s.DB.QueryRow(ctx, "SELECT coins FROM users WHERE id = $1", userID).Scan(&info.Coins)
	if err != nil {
		return models.InfoResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	// user inventory
	rows, err := s.DB.Query(ctx,
		`SELECT i.name, ui.quantity FROM user_inventory ui
        JOIN items i ON ui.item_id = i.id
        WHERE ui.user_id = $1`,
		userID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.InfoResponse{}, fmt.Errorf("%s: %w", op, repository.ErrNoInfo)
		}
		return models.InfoResponse{}, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var item models.InventoryItem
		if err := rows.Scan(&item.Type, &item.Quantity); err != nil {
			return models.InfoResponse{}, fmt.Errorf("%s: %w", op, err)
		}
		info.Inventory = append(info.Inventory, item)
	}

	// transactions when user is receiver
	receivedRows, err := s.DB.Query(ctx,
		`SELECT u.username, uch.amount FROM user_coin_history uch
	    JOIN users u ON uch.from_user_id = u.id
	    WHERE uch.to_user_id = $1`,
		userID,
	)
	if err != nil {
		return models.InfoResponse{}, fmt.Errorf("%s: %w", op, err)
	}
	defer receivedRows.Close()

	for receivedRows.Next() {
		var t models.ReceivedTransaction
		if err := receivedRows.Scan(&t.FromUser, &t.Amount); err != nil {
			return models.InfoResponse{}, fmt.Errorf("%s: %w", op, err)
		}
		info.CoinHistory.Received = append(info.CoinHistory.Received, t)
	}

	// transactions when user is sender
	sentRows, err := s.DB.Query(ctx,
		`SELECT u.username, uch.amount FROM user_coin_history uch
	    JOIN users u ON uch.to_user_id = u.id
		WHERE uch.from_user_id = $1`,
		userID,
	)
	if err != nil {
		return models.InfoResponse{}, fmt.Errorf("%s: %w", op, err)
	}
	defer sentRows.Close()

	for sentRows.Next() {
		var t models.SentTransaction
		if err := sentRows.Scan(&t.ToUser, &t.Amount); err != nil {
			return models.InfoResponse{}, fmt.Errorf("%s: %w", op, err)
		}
		info.CoinHistory.Sent = append(info.CoinHistory.Sent, t)
	}

	return info, nil
}
