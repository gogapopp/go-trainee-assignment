package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gogapopp/go-trainee-assignment/internal/repository"
	"go.uber.org/zap"
)

type buyItemService interface {
	BuyItem(ctx context.Context, userID int, itemName string) error
}

func BuyItemHandler(logger *zap.SugaredLogger, buyItemService buyItemService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.http-server.handlers.buy.BuyItemHandler"
		ctx := r.Context()

		logger = logger.With("req_id", middleware.GetReqID(ctx))

		userID, err := getUserIDFromContext(ctx)
		if err != nil {
			logger.Errorf("%s: %w", op, err)
			errorJSONResponse(w, http.StatusUnauthorized, "invalid token")
			return
		}

		itemName := chi.URLParam(r, "item")
		if itemName == "" {
			logger.Errorf("%s: %w", op, errors.New("empty item name"))
			errorJSONResponse(w, http.StatusBadRequest, "empty item name")
			return
		}

		if err := buyItemService.BuyItem(ctx, userID, itemName); err != nil {
			logger.Errorf("%s: %w", op, err)

			switch {
			case errors.Is(err, repository.ErrItemNotFound):
				errorJSONResponse(w, http.StatusNotFound, "item not found")
			case errors.Is(err, repository.ErrInsufficientFunds):
				errorJSONResponse(w, http.StatusBadRequest, "insufficient funds")
			default:
				errorJSONResponse(w, http.StatusInternalServerError, "internal error")
			}
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
