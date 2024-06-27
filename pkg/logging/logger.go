package logging

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"sync"
)

// contextKey is a private string type to prevent collisions in the context map.
type contextKey string

const (
	// loggerKey points to the value in the context where the logger is stored.
	loggerKey = contextKey("logger")
	// defaultLogLevel sets the default logger verbosity level.
	defaultLogLevel = slog.LevelInfo
	// defaultLogJSON sets the default logger output format.
	defaultLogJSON = true
)

var (
	// defaultLogger is the default logger. It is initialized once per package
	// include when calling [DefaultLogger].
	defaultLogger     *slog.Logger
	defaultLoggerOnce sync.Once
)

// NewLogger creates a [slog.Logger] with the given verbosity and output format.
func NewLogger(level slog.Level, json bool) *slog.Logger {
	var handler slog.Handler
	if json {
		handler = slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{Level: level},
		)
	} else {
		handler = slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{Level: level},
		)
	}
	return slog.New(handler)
}

// DefaultLogger returns the default logger for the package.
func DefaultLogger() *slog.Logger {
	defaultLoggerOnce.Do(func() {
		defaultLogger = NewLogger(defaultLogLevel, defaultLogJSON)
	})
	return defaultLogger
}

// WithLogger creates a new context with the provided logger attached.
func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// FromContext returns the logger stored in a context. If no such context
// exists, a default logger is returned.
func FromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
		return logger
	}
	return DefaultLogger()
}

// SlogLevel converts the given string to the appropriate log level. The
// supported input options are "info", "warn", "error", and "debug". All
// other inputs will result in an info level.
func SlogLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	case "debug":
		return slog.LevelDebug
	default:
		return slog.LevelInfo
	}
}
