package logger

import (
	"log/slog"
	"os"
	"strings"
)

const modeJSON = "json"

const (
	levelDebug   = "debug"
	levelInfo    = "info"
	levelWarning = "warning"
	levelError   = "error"
)

// New creates a logger with the given verbosity and output format.
func New(level slog.Level, json bool) *slog.Logger {
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

// NewFromEnv creates a logger from environment variables.
//
// Consumes LOG_LEVEL for determining the verbosity and consumes LOG_MODE
// for determining the log output format.
func NewFromEnv() *slog.Logger {
	level := slogLevel(strings.ToLower(os.Getenv("LOG_LEVEL")))
	json := strings.ToLower(os.Getenv("LOG_MODE")) == modeJSON

	return New(level, json)
}

// slogLevel converts the input to the appropriate slog level.
func slogLevel(level string) slog.Level {
	switch level {
	case levelDebug:
		return slog.LevelDebug
	case levelInfo:
		return slog.LevelInfo
	case levelWarning:
		return slog.LevelWarn
	case levelError:
		return slog.LevelError
	}

	return slog.LevelInfo
}
