package httpserver

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gogapopp/go-trainee-assignment/internal/http-server/handlers"
	"github.com/gogapopp/go-trainee-assignment/internal/libs/config"
	"github.com/gogapopp/go-trainee-assignment/internal/service"
	"go.uber.org/zap"
)

const (
	readHeaderTimeout = 5 * time.Second
	readTimeout       = 10 * time.Second
	writeTimeout      = 10 * time.Second
)

func New(cfg *config.Config, logger *zap.SugaredLogger, service *service.Service) *http.Server {
	router := chi.NewRouter()
	router.Use(
		middleware.Logger,
		middleware.RequestID,
	)

	router.Post("/api/auth", handlers.AuthHandler(logger, service))

	srv := &http.Server{
		Addr:              cfg.HTTPConifg.Addr,
		Handler:           router,
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
	}

	return srv
}
