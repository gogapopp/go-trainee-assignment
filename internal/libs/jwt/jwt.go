package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	errInvalidToken = errors.New("invalid token")
	errTokenExpired = errors.New("token expired")
)

const tokenExp = time.Minute * 100000000

type tokenClaims struct {
	UserID int
	jwt.RegisteredClaims
}

func GenerateJWTToken(jwtSecret string, userID int) (string, error) {
	claims := tokenClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func ParseJWTToken(jwtSecret, tokenString string) (int, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&tokenClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		},
	)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return 0, errTokenExpired
		}
		return 0, errInvalidToken
	}

	if claims, ok := token.Claims.(*tokenClaims); ok && token.Valid {
		return claims.UserID, nil
	}

	return 0, errInvalidToken
}
