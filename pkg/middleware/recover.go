package middleware

import "net/http"

// Recover from panics in a [http.Handler] and return an internal server error.
func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				js := []byte("{\"error\":\"internal server error\"}")
				w.Header().Set("content-type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(js)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
