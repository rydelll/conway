package logger

import (
	"log/slog"
	"os"
)

// New creates a [slog.Logger] with the given verbosity and output format.
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
