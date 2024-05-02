package server

import (
	"time"
)

// Option configures server by overriding a default setting.
type Option func(*Server)

// WithReadTimeout sets the maximum duration for reading HTTP requests.
// A zero or negative value means there is no timeout.
func WithReadTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.server.ReadTimeout = timeout
	}
}

// WithWriteTimeout sets the maximum duration for writing HTTP responses.
// A zero or negative value means there is no timeout.
func WithWriteTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.server.WriteTimeout = timeout
	}
}

// WithIdleTimeout sets the maximum duration to wait for the next request
// when keep-alive is enabled. A zero value means the value of ReadTimeout is
// used. If both are negative there is no timeout.
func WithIdleTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.server.IdleTimeout = timeout
	}
}

// WithShutdownTimeout sets the maximum duration to wait for a graceful
// shutdown before an immediate shutdown is forced.
func WithShutdownTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.shutdownTimeout = timeout
	}
}
