package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

// Config holds all application configuration loaded from environment variables.
type Config struct {
	Port               string `env:"PORT" envDefault:"8080"`
	LogFormat          string `env:"LOG_FORMAT" envDefault:"text"`
	LogLevel           string `env:"LOG_LEVEL" envDefault:"info"`
	CORSAllowedOrigins string `env:"CORS_ALLOWED_ORIGINS" envDefault:"http://localhost:4010"`
	OwnerTimezone      string `env:"OWNER_TIMEZONE" envDefault:"Europe/Moscow"`
}

// Load parses environment variables into Config.
func Load() (Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return Config{}, fmt.Errorf("failed to load config: %w", err)
	}
	return cfg, nil
}
