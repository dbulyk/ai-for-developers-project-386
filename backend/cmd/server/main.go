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
	"github.com/dbulyk/ai-for-developers-project-386/backend/internal/middleware"
	"github.com/go-chi/chi/v5"
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

	if err := run(ctx, listener, log); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error("server error", slog.Any("error", err))
		os.Exit(1)
	}

	log.Info("server stopped gracefully")
}

func run(ctx context.Context, listener net.Listener, log *slog.Logger) error {
	r := chi.NewRouter()
	r.Use(middleware.CORS(os.Getenv("CORS_ALLOWED_ORIGINS")))
	r.Use(middleware.Logger(log))
	r.Use(recoverer(log))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("ok")); err != nil {
			log.Error("failed to write health response", slog.Any("error", err))
		}
	})

	srv := &http.Server{
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

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

func recoverer(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rcv := recover(); rcv != nil {
					log.Error("panic recovered", slog.Any("panic", rcv))
					w.WriteHeader(http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
