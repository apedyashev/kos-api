package config

import (
	"errors"
	"os"
)

// Config contains runtime configuration for the API.
type Config struct {
	DatabaseURL string
}

// Load reads configuration from environment variables. DATABASE_URL should be
// the Supabase pooled PostgreSQL connection string.
func Load() (Config, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return Config{}, errors.New("DATABASE_URL is required")
	}

	return Config{DatabaseURL: databaseURL}, nil
}
