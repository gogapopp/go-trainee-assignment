package postgres

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/gogapopp/go-trainee-assignment/internal/models"
	"github.com/gogapopp/go-trainee-assignment/internal/repository"
	"github.com/jackc/pgx/v5"
)

func (s *Storage) AuthUser(ctx context.Context, user models.AuthRequest) (int, error) {
	const op = "internal.repository.postgres.auth.AuthUser"

	hashedPassword := s.generatePasswordHash(user.Password)

	var userID int
	err := s.DB.QueryRow(ctx,
		"INSERT INTO users(username, password_hash) VALUES($1, $2) ON CONFLICT (username) DO NOTHING RETURNING id",
		user.Username, hashedPassword,
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
	err = s.DB.QueryRow(ctx,
		"SELECT id, password_hash FROM users WHERE username = $1", user.Username).
		Scan(&userID, &passwordHashFromDB)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	if passwordHashFromDB != hashedPassword {
		return 0, fmt.Errorf("%s: %w", op, repository.ErrInvalidCredentials)
	}

	return userID, nil
}

func (s *Storage) generatePasswordHash(password string) string {
	hash := sha256.New()
	_, _ = hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(s.passSecret)))
}
