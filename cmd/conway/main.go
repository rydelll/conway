package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/rydelll/conway/internal/rest/middleware"
	"github.com/rydelll/conway/pkg/database"
	"github.com/rydelll/conway/pkg/logging"
	"github.com/rydelll/conway/pkg/server"
	"golang.org/x/sys/unix"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), unix.SIGINT, unix.SIGTERM, unix.SIGQUIT)
	defer cancel()

	if err := run(ctx, os.Args, os.Getenv, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

// run parses arguments and environment variables, initializes dependencies,
// and starts the application.
func run(ctx context.Context, args []string, getenv func(string) string, stderr io.Writer) error {
	// Arguments
	fs := flag.NewFlagSet("", flag.ExitOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {
		fmt.Fprintf(stderr, "Conway's Game of Life\n\n")
		fmt.Fprintf(stderr, "Usage:\n\n")
		fmt.Fprintf(stderr, "\t%s [options]\n\n", args[0])
		fmt.Fprintf(stderr, "Options:\n")
		fs.PrintDefaults()
		fmt.Fprintln(stderr)
	}

	var port int
	fs.IntVar(&port, "port", 8080, "port for the server to listen on")
	fs.Parse(args[1:])

	// Environment variables
	logLevel := logging.SlogLevel(getenv("LOG_LEVEL"))
	logJSON := strings.ToLower(getenv("LOG_MODE")) == "json"

	pgConfig := database.PGConfig{
		Scheme:      getenv("DB_SCHEME"),
		Host:        getenv("DB_HOST"),
		Name:        getenv("DB_NAME"),
		User:        getenv("DB_USER"),
		Password:    getenv("DB_PASSWORD"),
		SSLMode:     getenv("DB_SSLMODE"),
		SSLCert:     getenv("DB_SSLCERT"),
		SSLKey:      getenv("DB_SSLKEY"),
		SSLRootCert: getenv("DB_SSLROOTCERT"),
	}
	pgConfig.Port, _ = strconv.Atoi(getenv("DB_PORT"))
	pgConfig.ConnectTimeout, _ = strconv.Atoi(getenv("DB_CONNECT_TIMEOUT"))
	pgConfig.PoolMinConnections, _ = strconv.Atoi(getenv("DB_POOL_MIN_CONNS"))
	pgConfig.PoolMaxConnections, _ = strconv.Atoi(getenv("DB_POOL_MAX_CONNS"))
	pgConfig.PoolMaxConnLife, _ = time.ParseDuration(getenv("DB_POOL_MAX_CONN_LIFE"))
	pgConfig.PoolMaxConnIdle, _ = time.ParseDuration(getenv("DB_POOL_MAX_CONN_IDLE"))
	pgConfig.PoolHealthcheck, _ = time.ParseDuration(getenv("DB_POOL_HEALTHCHECK"))

	// Logging
	logger := logging.NewLogger(logLevel, logJSON)

	// Database
	db, err := database.NewPostgres(ctx, pgConfig)
	if err != nil {
		return err
	}
	defer db.Close()

	// Router and middleware
	rootMux := http.NewServeMux()
	subMux := http.NewServeMux()
	rootMux.Handle("/api/", http.StripPrefix("/api", middleware.Recover(subMux)))

	// Hello
	subMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})

	// Server
	server := server.New(logger, rootMux, port)
	if err := server.ListenAndServe(ctx); err != nil {
		return err
	}

	return nil
}
