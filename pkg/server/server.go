package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

const (
	defaultReadTimeout     = time.Second * 1
	defaultWriteTimeout    = time.Second * 5
	defaultIdleTimeout     = time.Second * 60
	defaultShutdownTimeout = time.Second * 10
)

// Server wraps and extends the [http.Server] functionality.
type Server struct {
	server          *http.Server
	logger          *slog.Logger
	shutdownTimeout time.Duration
}

// New creates a server with optional configuration.
func New(logger *slog.Logger, handler http.Handler, port int, opts ...Option) *Server {
	s := &Server{
		logger:          logger,
		shutdownTimeout: defaultShutdownTimeout,
	}

	s.server = &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      handler,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
		IdleTimeout:  defaultIdleTimeout,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelWarn),
	}

	// apply optional configuration
	for _, opt := range opts {
		opt(s)
	}

	return s
}

// ListenAndServe starts a server and blocks until the context is cancelled.
// When the context is cancelled, the server is gracefully stopped with the
// configured timeout.
//
// Once it has been stopped it is NOT safe for reuse.
func (s *Server) ListenAndServe(ctx context.Context) error {
	shutdownErrorChan := make(chan error, 1)
	logger := s.logger.WithGroup("server").With(slog.String("addr", s.server.Addr))

	go func() {
		<-ctx.Done()

		logger.Debug("shutdown signal recieved")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
		defer cancel()

		logger.Debug("shutting down", slog.String("timeout", s.shutdownTimeout.String()))
		shutdownErrorChan <- s.server.Shutdown(shutdownCtx)
	}()

	logger.Info("listening and serving HTTP")
	defer logger.Info("stopped listening and serving HTTP")

	err := s.server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	if err := <-shutdownErrorChan; err != nil {
		return err
	}

	return nil
}
