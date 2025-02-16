package config

import (
	"os"

	"github.com/joho/godotenv"
)

const (
	jwtSecret  = "JWT_SECRET"
	passSecret = "PASS_SECRET"
)

type Config struct {
	HTTPConifg *httpConfig
	PGConfig   *pgConfig
	JWTSecret  string
	PassSecret string
}

func New(path string) (*Config, error) {
	err := godotenv.Load(path)
	if err != nil {
		return nil, err
	}

	pgCfg := newPGConfig()

	httpCfg, err := newHTTPConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		HTTPConifg: httpCfg,
		PGConfig:   pgCfg,
		JWTSecret:  os.Getenv(jwtSecret),
		PassSecret: os.Getenv(passSecret),
	}, nil
}
