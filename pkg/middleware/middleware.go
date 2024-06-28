package middleware

import "net/http"

// contextKey is a private string type to prevent collisions in the context map.
type contextKey string

// MiddlewareFunc recieves and returns a [http.Handler]. Typically, the returned
// handler does something with the [http.ResponseWriter] and [http.Request]
// passed to it, then calls the handler passed in to the MiddlewareFunc.
type MiddelwareFunc func(next http.Handler) http.Handler

// Use wraps the given handler with the middleware in the order passed in.
func Use(h http.Handler, mwf ...MiddelwareFunc) http.Handler {
	for _, fn := range mwf {
		h = fn(h)
	}
	return h
}
