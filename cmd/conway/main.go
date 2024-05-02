package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/rydelll/conway/pkg/database/postgres"
	"github.com/rydelll/conway/pkg/logger"
	"github.com/rydelll/conway/pkg/server"
	"golang.org/x/sys/unix"
)

func main() {
	// Setup context to listen for signals
	ctx, cancel := signal.NotifyContext(context.Background(), unix.SIGINT, unix.SIGTERM, unix.SIGQUIT)
	defer cancel()

	// Get logger
	logger := logger.NewFromEnv()

	// Setup database connection
	db, err := postgres.NewFromEnv()
	if err != nil {
		logger.Error("establish database connection", "error", err)
		os.Exit(1)
	}

	// Domain object
	// objectRepo := object.NewRepository(dbConn)
	// objectService := object.NewService(objectRepo)
	// objectHandler := object.NewHandler(objectService)
	// ...

	// Setup HTTP router
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})

	// Start HTTP server
	server := server.New(logger, mux, 8080)
	if err := server.ListenAndServe(ctx); err != nil {
		logger.Error("server shutdown", "error", err)
		os.Exit(1)
	}

	logger.Info("application shutdown successful")
}
