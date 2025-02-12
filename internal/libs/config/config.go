package config

import (
	"os"

	"github.com/joho/godotenv"
)

const jwtSecret = "JWT_SECRET"

type Config struct {
	HTTPConifg *httpConfig
	PGConfig   *pgConfig
	JWT_SECRET string
}

func New(path string) (*Config, error) {
	err := godotenv.Load(path)
	if err != nil {
		return nil, err
	}

	pgCfg, err := newPGConfig()
	if err != nil {
		return nil, err
	}

	httpCfg, err := newHTTPConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		HTTPConifg: httpCfg,
		PGConfig:   pgCfg,
		JWT_SECRET: os.Getenv(jwtSecret),
	}, nil
}
