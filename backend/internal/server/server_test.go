package server

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dbulyk/ai-for-developers-project-386/backend/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_RoutesAssembled(t *testing.T) {
	log := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelError}))
	cfg := config.Config{
		Port:               "8080",
		LogFormat:          "text",
		LogLevel:           "error",
		CORSAllowedOrigins: "http://localhost:4010",
		OwnerTimezone:      "Europe/Moscow",
	}

	srv := New(cfg, log)
	server := httptest.NewServer(srv.Handler)
	defer server.Close()

	t.Run("health", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/health")
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("public event types list", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/public/event-types")
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("admin event types create validates input", func(t *testing.T) {
		resp, err := http.Post(server.URL+"/admin/event-types", "application/json", strings.NewReader(`{}`))
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("cors preflight", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodOptions, server.URL+"/public/event-types", nil)
		require.NoError(t, err)
		req.Header.Set("Origin", "http://localhost:4010")

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
		assert.Equal(t, "http://localhost:4010", resp.Header.Get("Access-Control-Allow-Origin"))
	})

	t.Run("spa fallback serves index.html", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/admin")
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Contains(t, string(body), `<div id="root">`)
		assert.Contains(t, resp.Header.Get("Content-Type"), "text/html")
	})
}
