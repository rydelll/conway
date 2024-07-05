package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRecover(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		handler http.Handler
		code    int
	}{
		{
			name: "default",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}),
			code: http.StatusOK,
		},
		{
			name: "panic",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic("oops")
			}),
			code: http.StatusInternalServerError,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			r := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()
			Recover(tc.handler).ServeHTTP(w, r)
			if diff := cmp.Diff(tc.code, w.Result().StatusCode); diff != "" {
				t.Errorf("mismatch (-want, +got):\n%s", diff)
			}
		})
	}

}
