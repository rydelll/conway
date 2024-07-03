package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRequestID(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		ctx  context.Context
	}{
		{name: "string", ctx: context.WithValue(context.Background(), reqIDKey, "123")},
		{name: "empty", ctx: context.Background()},
		{name: "int", ctx: context.WithValue(context.Background(), reqIDKey, 123)},
	}

	// use handler to expose request ID outside request call stack
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		reqID := RequestIDFromContext(ctx)
		w.Write([]byte(reqID))
	})

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			r := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()

			RequestID(handler).ServeHTTP(w, r)
			if w.Body.String() == "" {
				t.Error("expected request ID to be non empty string")
			}
		})
	}
}

func TestRequestIDFromContext(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		reqID any
		want  string
	}{
		{name: "string", reqID: "123", want: "123"},
		{name: "empty", reqID: nil, want: ""},
		{name: "int", reqID: 123, want: ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.WithValue(context.Background(), reqIDKey, tc.reqID)

			got := RequestIDFromContext(ctx)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}
