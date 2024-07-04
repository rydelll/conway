package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestUse(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("X"))
	})
	mw1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Y"))
			next.ServeHTTP(w, r)
			w.Write([]byte("Y"))
		})
	}
	mw2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Z"))
			next.ServeHTTP(w, r)
			w.Write([]byte("Z"))
		})
	}

	cases := []struct {
		name string
		use  http.Handler
		want string
	}{
		{name: "none", use: Use(handler), want: "X"},
		{name: "single", use: Use(handler, mw1), want: "YXY"},
		{name: "multiple", use: Use(handler, mw1, mw2), want: "ZYXYZ"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			r := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()

			tc.use.ServeHTTP(w, r)
			if diff := cmp.Diff(tc.want, w.Body.String()); diff != "" {
				t.Errorf("mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}
