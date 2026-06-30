package server

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/dbulyk/ai-for-developers-project-386/backend/internal/clock"
	"github.com/dbulyk/ai-for-developers-project-386/backend/internal/config"
	"github.com/dbulyk/ai-for-developers-project-386/backend/internal/handlers"
	"github.com/dbulyk/ai-for-developers-project-386/backend/internal/middleware"
	"github.com/dbulyk/ai-for-developers-project-386/backend/internal/store"
	"github.com/go-chi/chi/v5"
)

// New assembles and returns an http.Server with all application routes,
// middleware and dependencies wired together. The returned server is ready
// to be started with Serve or ListenAndServe.
func New(cfg config.Config, logger *slog.Logger) *http.Server {
	s := store.NewMemoryStore()
	c := clock.RealClock{}

	tz, err := time.LoadLocation(cfg.OwnerTimezone)
	if err != nil {
		logger.Warn("failed to load owner timezone, falling back to UTC",
			slog.String("timezone", cfg.OwnerTimezone),
			slog.Any("error", err))
		tz = time.UTC
	}

	r := chi.NewRouter()
	r.Use(middleware.CORS(cfg.CORSAllowedOrigins))
	r.Use(middleware.Logger(logger))
	r.Use(recoverer(logger))

	handlers.NewAdminEventTypesHandler(s).RegisterRoutes(r)
	handlers.NewPublicEventTypesHandler(s, c, tz).RegisterRoutes(r)
	handlers.NewPublicBookingsHandler(s, c, tz).RegisterRoutes(r)
	handlers.NewAdminBookingsHandler(s, c).RegisterRoutes(r)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("ok")); err != nil {
			logger.Error("failed to write health response", slog.Any("error", err))
		}
	})

	return &http.Server{
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
	}
}

// recoverer recovers from panics, logs them and returns 500 to the client.
func recoverer(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rcv := recover(); rcv != nil {
					logger.Error("panic recovered", slog.Any("panic", rcv))
					w.WriteHeader(http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
