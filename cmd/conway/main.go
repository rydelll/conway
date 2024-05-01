package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"github.com/rydelll/conway/internal/logger"
	"github.com/rydelll/conway/internal/server"
	"golang.org/x/sys/unix"
)

func main() {
	// Setup context to listen for signals
	ctx, cancel := signal.NotifyContext(context.Background(), unix.SIGINT, unix.SIGTERM, unix.SIGQUIT)
	defer cancel()

	// Get logger
	logger := logger.NewFromEnv()

	// Start the application
	if err := run(ctx, logger); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	logger.Info("application shutdown successful")
}

func run(ctx context.Context, logger *slog.Logger) error {
	// Database connection
	// ...

	// Database migration
	// ...

	// Domain object
	// objectRepo := object.NewRepository(dbConn)
	// objectService := object.NewService(objectRepo)
	// objectHandler := object.NewHandler(objectService)
	// ...

	// Router
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})

	// Server
	server := server.New(logger, mux, 8080)
	return server.ServeHTTP(ctx)
}
