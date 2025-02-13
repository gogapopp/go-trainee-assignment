package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var errUnknownClaimsType = errors.New("unknown claims type")

const tokenExp = time.Minute * 10

type tokenClaims struct {
	userID int
	jwt.RegisteredClaims
}

func GenerateJWTToken(jwtSecret string, userID int) (string, error) {
	claims := tokenClaims{
		userID,
		jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return ss, nil
}

func ParseJWTToken(jwtSecret, userJWTToken string) (int, error) {
	token, err := jwt.ParseWithClaims(userJWTToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return 0, err
	} else if claims, ok := token.Claims.(*tokenClaims); ok {
		return claims.userID, nil
	}
	return 0, errUnknownClaimsType
}
