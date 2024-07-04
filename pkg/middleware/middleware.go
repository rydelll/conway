package middleware

import "net/http"

// contextKey is a private string type to prevent collisions in the context map.
type contextKey string

// Use wraps the given handler with the middleware in the order passed in.
//
// For example Use(handler, A, B) first wraps handler in A, then in B. This
// results in B executing, then A executing, then handler executing. When the
// handler is done, execution returns to A, then finally returns to B.
func Use(h http.Handler, mwf ...func(http.Handler) http.Handler) http.Handler {
	for _, fn := range mwf {
		h = fn(h)
	}
	return h
}
