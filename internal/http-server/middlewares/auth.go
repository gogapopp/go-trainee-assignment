package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/gogapopp/go-trainee-assignment/internal/libs/jwt"
)

type CtxKeyUserID struct{}

func AuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" {
				http.Error(w, "authorization header is missing", http.StatusUnauthorized)
				return
			}

			headerParts := strings.Split(header, " ")
			if len(headerParts) != 2 || headerParts[0] != "Bearer" {
				http.Error(w, "invalid token format", http.StatusUnauthorized)
				return
			}

			token := headerParts[1]
			userID, err := jwt.ParseJWTToken(jwtSecret, token)
			if err != nil {
				http.Error(w, "invalid token or expired", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), CtxKeyUserID{}, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
