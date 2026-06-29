package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_Defaults(t *testing.T) {
	cleanup := clearEnv(t)
	defer cleanup()

	cfg, err := Load()
	require.NoError(t, err)

	assert.Equal(t, "8080", cfg.Port)
	assert.Equal(t, "text", cfg.LogFormat)
	assert.Equal(t, "info", cfg.LogLevel)
	assert.Equal(t, "http://localhost:4010", cfg.CORSAllowedOrigins)
	assert.Equal(t, "Europe/Moscow", cfg.OwnerTimezone)
}

func TestLoad_CustomValues(t *testing.T) {
	cleanup := clearEnv(t)
	defer cleanup()

	setEnv(t, "PORT", "3000")
	setEnv(t, "LOG_FORMAT", "json")
	setEnv(t, "LOG_LEVEL", "debug")
	setEnv(t, "CORS_ALLOWED_ORIGINS", "https://example.com,https://app.example.com")
	setEnv(t, "OWNER_TIMEZONE", "Europe/London")

	cfg, err := Load()
	require.NoError(t, err)

	assert.Equal(t, "3000", cfg.Port)
	assert.Equal(t, "json", cfg.LogFormat)
	assert.Equal(t, "debug", cfg.LogLevel)
	assert.Equal(t, "https://example.com,https://app.example.com", cfg.CORSAllowedOrigins)
	assert.Equal(t, "Europe/London", cfg.OwnerTimezone)
}

func clearEnv(t *testing.T) func() {
	t.Helper()

	keys := []string{"PORT", "LOG_FORMAT", "LOG_LEVEL", "CORS_ALLOWED_ORIGINS", "OWNER_TIMEZONE"}
	original := make(map[string]*string, len(keys))
	for _, k := range keys {
		v, ok := os.LookupEnv(k)
		if ok {
			original[k] = &v
		} else {
			original[k] = nil
		}
		require.NoError(t, os.Unsetenv(k))
	}

	return func() {
		for _, k := range keys {
			if original[k] == nil {
				_ = os.Unsetenv(k)
			} else {
				_ = os.Setenv(k, *original[k])
			}
		}
	}
}

func setEnv(t *testing.T, key, value string) {
	t.Helper()
	require.NoError(t, os.Setenv(key, value))
}
