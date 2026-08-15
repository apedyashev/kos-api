package config

import (
	"errors"
	"os"
	"strconv"
)

// Config contains runtime configuration for the API.
type Config struct {
	DatabaseURL string
	DebugErrors bool
}

// Load reads configuration from environment variables. DATABASE_URL should be
// the Supabase pooled PostgreSQL connection string.
func Load() (Config, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return Config{}, errors.New("DATABASE_URL is required")
	}

	debugErrors, err := strconv.ParseBool(os.Getenv("DEBUG_ERRORS"))
	if err != nil {
		debugErrors = false
	}

	return Config{DatabaseURL: databaseURL, DebugErrors: debugErrors}, nil
}
