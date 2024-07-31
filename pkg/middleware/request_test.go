package middleware

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRequestID(t *testing.T) {
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
		io.WriteString(w, reqID)
	})

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()
			RequestID(handler).ServeHTTP(w, r)
			reqID, _ := io.ReadAll(w.Result().Body)
			if len(reqID) == 0 {
				t.Error("expected request ID to be non empty string")
			}
		})
	}
}

func TestRequestIDFromContext(t *testing.T) {
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
			ctx := context.WithValue(context.Background(), reqIDKey, tc.reqID)
			got := RequestIDFromContext(ctx)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}
