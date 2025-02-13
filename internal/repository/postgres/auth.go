package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/gogapopp/go-trainee-assignment/internal/models"
	"github.com/gogapopp/go-trainee-assignment/internal/repository"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 12

func (s *storage) AuthUser(ctx context.Context, user models.AuthRequest) (int, error) {
	const op = "internal.repository.postgres.auth.AuthUser"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcryptCost)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	var userID int
	err = s.db.QueryRow(ctx,
		"INSERT INTO users(username, password_hash) VALUES($1, $2) ON CONFLICT (username) DO NOTHING RETURNING id",
		user.Username, string(hashedPassword),
	).Scan(&userID)
	// user is created successfully
	if err == nil {
		return userID, nil
	}
	// other internal error
	if !errors.Is(err, pgx.ErrNoRows) {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	// if user exists
	var passwordHashFromDB string
	err = s.db.QueryRow(ctx,
		"SELECT id, password_hash FROM users WHERE username = $1", user.Username).
		Scan(&userID, &passwordHashFromDB)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHashFromDB), []byte(user.Password)); err != nil {
		return 0, fmt.Errorf("%s: %w", op, repository.ErrInvalidCredentials)
	}

	return userID, nil
}
