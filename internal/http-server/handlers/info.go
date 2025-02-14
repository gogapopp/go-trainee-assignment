package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/gogapopp/go-trainee-assignment/internal/models"
	"github.com/gogapopp/go-trainee-assignment/internal/repository"
	"go.uber.org/zap"
)

type infoService interface {
	GetUserInfo(ctx context.Context, userID int) (models.InfoResponse, error)
}

func InfoHandler(logger *zap.SugaredLogger, infoService infoService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.http-server.handlers.info.InfoHandler"
		ctx := r.Context()

		logger = logger.With("req_id", middleware.GetReqID(ctx))

		userID, err := getUserIDFromContext(ctx)
		if err != nil {
			logger.Errorf("%s: %w", op, err)
			errorJSONResponse(w, http.StatusUnauthorized, "invalid token")
			return
		}

		info, err := infoService.GetUserInfo(ctx, userID)
		if err != nil {
			logger.Errorf("%s: %w", op, err)

			switch {
			case errors.Is(err, repository.ErrUserNotFound):
				errorJSONResponse(w, http.StatusNotFound, "user not found")
			default:
				errorJSONResponse(w, http.StatusInternalServerError, "internal error")
			}
			return
		}

		jsonResponse(w, http.StatusOK, info)
	}
}
