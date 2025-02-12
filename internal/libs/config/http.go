package config

import (
	"errors"
	"fmt"
	"os"
)

const serverPort = "SERVER_PORT"

type httpConfig struct {
	Addr string
}

func newHTTPConfig() (*httpConfig, error) {
	port := os.Getenv(serverPort)

	if len(port) == 0 {
		return nil, errors.New("empty SERVER_PORT env")
	}

	return &httpConfig{
		Addr: fmt.Sprintf(":%s", port),
	}, nil
}
