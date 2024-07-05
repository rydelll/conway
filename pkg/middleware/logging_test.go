package middleware

import (
	"bytes"
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/rydelll/conway/pkg/logging"
)

func TestLogger(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		ctx  context.Context
		want string
	}{
		{
			name: "string",
			ctx:  context.WithValue(context.Background(), reqIDKey, "123"),
			want: "{\"level\":\"INFO\",\"msg\":\"test\",\"requestID\":\"123\"}\n",
		},
		{
			name: "int",
			ctx:  context.WithValue(context.Background(), reqIDKey, 123),
			want: "{\"level\":\"INFO\",\"msg\":\"test\"}\n",
		},
		{
			name: "empty",
			ctx:  context.Background(),
			want: "{\"level\":\"INFO\",\"msg\":\"test\"}\n",
		},
	}

	buf := bytes.NewBuffer(nil)
	logger := logging.NewLoggerTimeless(buf, slog.LevelInfo, true)
	handler := Logger(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := logging.FromContext(ctx)
		logger.Info("test")
	}))

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			buf.Reset()
			r := httptest.NewRequest(http.MethodGet, "/", nil).WithContext(tc.ctx)
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, r)
			if diff := cmp.Diff(tc.want, buf.String()); diff != "" {
				t.Errorf("mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestLogRequest(t *testing.T) {
	t.Parallel()

	want := "{\"level\":\"INFO\",\"msg\":\"incoming request\"," +
		"\"request\":{\"proto\":\"HTTP/1.1\",\"method\":\"GET\"," +
		"\"URI\":\"/\",\"remoteAddr\":\"192.0.2.1:1234\"}}\n" +
		"{\"level\":\"INFO\",\"msg\":\"request complete\"}\n"

	buf := bytes.NewBuffer(nil)
	logger := logging.NewLoggerTimeless(buf, slog.LevelInfo, true)
	ctx := logging.WithLogger(context.Background(), logger)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	r := httptest.NewRequest(http.MethodGet, "/", nil).WithContext(ctx)
	w := httptest.NewRecorder()
	LogRequest(handler).ServeHTTP(w, r)
	if diff := cmp.Diff(want, buf.String()); diff != "" {
		t.Errorf("mismatch (-want, +got):\n%s", diff)
	}
}
