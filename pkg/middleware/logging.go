package middleware

import (
	"log/slog"
	"net/http"

	"github.com/rydelll/conway/pkg/logging"
)

// Logger populates the logger into the requests context and adds a request ID
// to the logger if one exists.
func Logger(logger *slog.Logger) MiddelwareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			logger := logger
			if reqID := RequestIDFromContext(ctx); reqID != "" {
				logger = logger.With(slog.String("requestID", reqID))
			}
			ctx = logging.WithLogger(ctx, logger)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// LogRequest logs all incoming requests.
func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := logging.FromContext(ctx)
		logger = logger.With(slog.Group("request",
			slog.String("proto", r.Proto),
			slog.String("method", r.Method),
			slog.String("remoteAddr", r.RemoteAddr),
			slog.String("URI", r.RequestURI),
		))
		logger.Info("incoming")
		next.ServeHTTP(w, r)
	})
}
