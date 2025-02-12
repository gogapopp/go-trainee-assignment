package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gogapopp/go-trainee-assignment/internal/libs/logger"
)

func main() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	var (
		logger = must(logger.New())
	)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutdown Server...")
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
