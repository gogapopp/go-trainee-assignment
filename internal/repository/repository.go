package repository

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrItemNotFound       = errors.New("item not found")
	ErrInsufficientFunds  = errors.New("insufficient funds")
	ErrUserNotFound       = errors.New("user not found")
	ErrSameUser           = errors.New("cant send to yourself")
	ErrNoInfo             = errors.New("user doesnt have info")
)
