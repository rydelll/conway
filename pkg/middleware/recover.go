package middleware

import (
	"log/slog"
	"net/http"

	"github.com/rydelll/conway/pkg/logging"
)

// Recover from panics in a [http.Handler] and return an internal server error.
func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := logging.FromContext(ctx)
		defer func() {
			if p := recover(); p != nil {
				logger.Error("http handler panic", slog.Any("panic", p))
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
