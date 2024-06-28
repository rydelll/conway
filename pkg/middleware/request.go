package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

// reqIDKey points to the value in the context where the request ID is stored.
const reqIDKey = contextKey("requestID")

// RequestID populates the request context with a random UUID
// if one does not already exist.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if reqID := RequestIDFromContext(ctx); reqID == "" {
			reqID = uuid.NewString()
			ctx = context.WithValue(ctx, reqIDKey, reqID)
			r = r.WithContext(ctx)
		}
		next.ServeHTTP(w, r)
	})
}

// RequestIDFromContext pulls the request ID from the context. If one was not
// set, it returns the empty string.
func RequestIDFromContext(ctx context.Context) string {
	v := ctx.Value(reqIDKey)
	if reqID, ok := v.(string); ok {
		return reqID
	}
	return ""
}
