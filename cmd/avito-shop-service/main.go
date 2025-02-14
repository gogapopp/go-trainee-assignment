package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpserver "github.com/gogapopp/go-trainee-assignment/internal/http-server"
	"github.com/gogapopp/go-trainee-assignment/internal/libs/config"
	"github.com/gogapopp/go-trainee-assignment/internal/libs/logger"
	"github.com/gogapopp/go-trainee-assignment/internal/repository/postgres"
	"github.com/gogapopp/go-trainee-assignment/internal/service"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	envPath        = ".env"
	migrationsPath = "./migrations"
)

func main() {
	ctx := context.Background()

	var (
		logger = must(logger.New())
		config = must(config.New(envPath))

		repository = must(postgres.New(config.PGConfig.DSN))
		migrations = must(migrate.New(
			fmt.Sprintf("file://%s", migrationsPath),
			config.PGConfig.DSN,
		))
	)
	defer logger.Sync()
	defer repository.Close(ctx)

	if err := migrations.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Fatal(err)
	}

	service := service.New(repository, config.JWTSecret)

	srv := httpserver.New(config, logger, service)

	go func() {
		logger.Info("Server started")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("listen: %s\n", err)
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutdown Server...")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("Server Shutdown: %s", err)
	}

	select {
	case <-ctx.Done():
		logger.Info("timeout of 5 seconds.")
	}
	logger.Info("Server exiting")
}

func must[T any](v T, err error) T {
	if err != nil {
		log.Fatal(err)
	}
	return v
}
