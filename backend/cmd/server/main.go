package main

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dbulyk/ai-for-developers-project-386/backend/internal/config"
	"github.com/dbulyk/ai-for-developers-project-386/backend/internal/logger"
	"github.com/dbulyk/ai-for-developers-project-386/backend/internal/server"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		slog.Default().Error("failed to load config", slog.Any("error", err))
		os.Exit(1)
	}

	log, err := logger.Setup(cfg.LogFormat, cfg.LogLevel, nil)
	if err != nil {
		slog.Default().Error("failed to setup logger", slog.Any("error", err))
		os.Exit(1)
	}

	listener, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		log.Error("failed to listen", slog.Any("error", err))
		os.Exit(1)
	}

	if err := run(ctx, listener, cfg, log); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error("server error", slog.Any("error", err))
		os.Exit(1)
	}

	log.Info("server stopped gracefully")
}

func run(ctx context.Context, listener net.Listener, cfg config.Config, log *slog.Logger) error {
	srv := server.New(cfg, log)

	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.Serve(listener)
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			return err
		}
		return nil
	}
}
