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

type sendCoinService interface {
	SendCoins(ctx context.Context, senderID int, req models.SendCoinRequest) error
}

func SendCoinHandler(logger *zap.SugaredLogger, sendCoinService sendCoinService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.http-server.handlers.send_coin.SendCoinHandler"
		ctx := r.Context()

		logger = logger.With("req_id", middleware.GetReqID(ctx))

		userID, err := getUserIDFromContext(ctx)
		if err != nil {
			logger.Errorf("%s: %w", op, err)
			errorJSONResponse(w, http.StatusUnauthorized, "invalid token")
			return
		}

		var req models.SendCoinRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.Errorf("%s: %w", op, err)
			errorJSONResponse(w, http.StatusBadRequest, "invalid request")
			return
		}

		if err := sendCoinService.SendCoins(ctx, userID, req); err != nil {
			var valErr validator.ValidationErrors
			logger.Errorf("%s: %w", op, err)

			switch {
			case errors.Is(err, repository.ErrInsufficientFunds):
				errorJSONResponse(w, http.StatusBadRequest, "insufficient funds")
			case errors.Is(err, repository.ErrUserNotFound):
				errorJSONResponse(w, http.StatusNotFound, "user not found")
			case errors.Is(err, repository.ErrSameUser):
				errorJSONResponse(w, http.StatusBadRequest, "cant send to yourself")
			case errors.As(err, &valErr):
				errorJSONResponse(w, http.StatusBadRequest, "validation error")
			default:
				errorJSONResponse(w, http.StatusInternalServerError, "internal error")
			}
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
