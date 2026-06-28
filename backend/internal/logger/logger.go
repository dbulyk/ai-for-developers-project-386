package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
)

// Setup creates a slog.Logger with the requested format and level.
// If w is nil, logs are written to os.Stdout.
func Setup(format, level string, w io.Writer) (*slog.Logger, error) {
	if w == nil {
		w = os.Stdout
	}

	lvl, err := parseLevel(level)
	if err != nil {
		return nil, err
	}

	opts := &slog.HandlerOptions{Level: lvl}

	var handler slog.Handler
	switch strings.ToLower(format) {
	case "text":
		handler = slog.NewTextHandler(w, opts)
	case "json":
		handler = slog.NewJSONHandler(w, opts)
	default:
		return nil, fmt.Errorf("unsupported log format %q", format)
	}

	return slog.New(handler), nil
}

func parseLevel(level string) (slog.Level, error) {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug, nil
	case "info":
		return slog.LevelInfo, nil
	case "warn":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return slog.Level(0), fmt.Errorf("unsupported log level %q", level)
	}
}
