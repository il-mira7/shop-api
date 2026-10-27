package config

import (
	"log"
	"os"
)

type Config struct {
	DatabaseURL string
}

func NewConfig() *Config {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		log.Fatal("environment variable DATABASE_URL is required but not set")
	}
	return &Config{DatabaseURL: url}
}
