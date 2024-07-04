package middleware

import (
	"log/slog"
	"net/http"

	"github.com/rydelll/conway/pkg/logging"
)

// Logger populates the logger into the requests context and adds a request ID
// to the logger if one exists.
func Logger(originLogger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			logger := originLogger
			if reqID := RequestIDFromContext(ctx); reqID != "" {
				logger = logger.With(slog.String("requestID", reqID))
			}
			ctx = logging.WithLogger(ctx, logger)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

// LogRequest logs all incoming requests and when they are complete.
func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := logging.FromContext(ctx)
		logger.Info("incoming request", slog.Group("request",
			slog.String("proto", r.Proto),
			slog.String("method", r.Method),
			slog.String("URI", r.RequestURI),
			slog.String("remoteAddr", r.RemoteAddr),
		))
		next.ServeHTTP(w, r)
		logger.Info("request complete")
	})
}
