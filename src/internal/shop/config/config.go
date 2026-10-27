package config

import (
	"log"
	"os"
)

type Config struct {
	DatabaseURL string
	AuthAddr    string
}

func NewConfig() *Config {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		log.Fatal("environment variable DATABASE_URL is required but not set")
	}

	authAddr := os.Getenv("AUTH_GRPC_ADDR")
	if authAddr == "" {
		authAddr = "localhost:50051"
	}

	return &Config{DatabaseURL: url, AuthAddr: authAddr}
}
