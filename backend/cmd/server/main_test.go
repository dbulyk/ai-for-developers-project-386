package main

import (
	"context"
	"io"
	"log/slog"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/dbulyk/ai-for-developers-project-386/backend/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRun_GracefulShutdown(t *testing.T) {
	log := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelError}))
	cfg := config.Config{
		Port:               "8080",
		LogFormat:          "text",
		LogLevel:           "error",
		CORSAllowedOrigins: "http://localhost:4010",
		OwnerTimezone:      "Europe/Moscow",
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer func() { _ = listener.Close() }()

	ctx, cancel := context.WithCancel(context.Background())
	errCh := make(chan error, 1)
	go func() {
		errCh <- run(ctx, listener, cfg, log)
	}()

	baseURL := "http://" + listener.Addr().String()
	var resp *http.Response
	for i := 0; i < 50; i++ {
		resp, err = http.Get(baseURL + "/health")
		if err == nil {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NoError(t, resp.Body.Close())

	cancel()
	select {
	case err := <-errCh:
		assert.NoError(t, err)
	case <-time.After(2 * time.Second):
		t.Fatal("server did not shut down in time")
	}
}
