package httpserver

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gogapopp/go-trainee-assignment/internal/http-server/handlers"
	"github.com/gogapopp/go-trainee-assignment/internal/http-server/middlewares"
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
		middleware.Recoverer,
	)

	router.Post("/api/auth", handlers.AuthHandler(logger, service))
	router.Group(func(router chi.Router) {
		router.Use(middlewares.AuthMiddleware(cfg.JWTSecret))

		router.Get("/api/buy/{item}", handlers.BuyItemHandler(logger, service))
		router.Get("/api/info", handlers.InfoHandler(logger, service))
		router.Post("/api/sendCoin", handlers.SendCoinHandler(logger, service))
	})
	srv := &http.Server{
		Addr:              cfg.HTTPConifg.Addr,
		Handler:           router,
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
	}

	return srv
}
