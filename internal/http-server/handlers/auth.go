package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator"
	"github.com/gogapopp/go-trainee-assignment/internal/models"
	"github.com/gogapopp/go-trainee-assignment/internal/repository"
	"go.uber.org/zap"
)

type authService interface {
	AuthUser(ctx context.Context, req models.AuthRequest) (string, error)
}

func AuthHandler(logger *zap.SugaredLogger, authService authService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.http-server.handler.auth.authHandler"
		ctx := r.Context()

		logger = logger.With("req_id", middleware.GetReqID(ctx))

		var req models.AuthRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.Errorf("%s: %w", op, err)
			errorJSONResponse(w, http.StatusBadRequest, "invalid request")
			return
		}

		token, err := authService.AuthUser(ctx, req)
		if err != nil {
			var valErr validator.ValidationErrors
			logger.Errorf("%s: %w", op, err)

			switch {
			case errors.Is(err, repository.ErrInvalidCredentials):
				errorJSONResponse(w, http.StatusUnauthorized, "invalid credentials")
			case errors.As(err, &valErr):
				errorJSONResponse(w, http.StatusBadRequest, "validation error")
			default:
				errorJSONResponse(w, http.StatusInternalServerError, "internal server error")
			}
			return
		}

		jsonResponse(w, http.StatusOK, models.AuthResponse{Token: token})
	}
}
