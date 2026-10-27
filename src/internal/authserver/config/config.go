package config

import (
	"log"
	"os"
)

type Config struct {
	DatabaseURL string
	Secret      string
}

func NewConfig() *Config {
	url := os.Getenv("AUTH_DATABASE_URL")
	if url == "" {
		log.Fatal("environment variable AUTH_DATABASE_URL is required but not set")
	}

	secret := os.Getenv("JWT_SECRET")

	if secret == "" {
		log.Fatal("environment variable JWT_SECRET is required but not set")
	}
	return &Config{DatabaseURL: url, Secret: secret}
}
