package logging

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewLogger(t *testing.T) {
	t.Parallel()

	cases := []struct {
		level slog.Level
		json  bool
	}{
		{level: slog.LevelDebug, json: true},
		{level: slog.LevelDebug, json: false},
		{level: slog.LevelError, json: true},
		{level: slog.LevelError, json: false},
		{level: slog.LevelWarn, json: true},
		{level: slog.LevelWarn, json: false},
		{level: slog.LevelInfo, json: true},
		{level: slog.LevelInfo, json: false},
	}

	for _, tc := range cases {
		t.Run(tc.level.String(), func(t *testing.T) {
			t.Parallel()

			if NewLogger(nil, tc.level, tc.json) == nil {
				t.Errorf("expected logger to never be nil")
			}
		})
	}
}

func TestNewLoggerTimeless(t *testing.T) {
	t.Parallel()

	cases := []struct {
		level slog.Level
		json  bool
	}{
		{level: slog.LevelDebug, json: true},
		{level: slog.LevelDebug, json: false},
		{level: slog.LevelError, json: true},
		{level: slog.LevelError, json: false},
		{level: slog.LevelWarn, json: true},
		{level: slog.LevelWarn, json: false},
		{level: slog.LevelInfo, json: true},
		{level: slog.LevelInfo, json: false},
	}

	for _, tc := range cases {
		t.Run(tc.level.String(), func(t *testing.T) {
			t.Parallel()

			if NewLoggerTimeless(nil, tc.level, tc.json) == nil {
				t.Errorf("expected logger to never be nil")
			}
		})
	}

	want := "{\"level\":\"INFO\",\"msg\":\"test\"}\n"
	buf := bytes.NewBuffer(nil)
	logger := NewLoggerTimeless(buf, slog.LevelInfo, true)
	logger.Info("test")
	if diff := cmp.Diff(want, buf.String()); diff != "" {
		t.Errorf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestDefaultLogger(t *testing.T) {
	t.Parallel()

	logger1 := DefaultLogger()
	if logger1 == nil {
		t.Fatal("expected logger to never be nil")
	}
	logger2 := DefaultLogger()
	if logger2 == nil {
		t.Fatal("expected logger to never be nil")
	}
	if logger1 != logger2 {
		t.Errorf("expected %#v to be %#v", logger1, logger2)
	}
}

func TestContext(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	logger1 := DefaultLogger()
	logger2 := FromContext(ctx)
	if logger1 != logger2 {
		t.Errorf("expected %#v to be %#v", logger1, logger2)
	}

	logger1 = NewLogger(nil, slog.LevelDebug, false)
	ctx = WithLogger(ctx, logger1)
	logger2 = FromContext(ctx)
	if logger1 != logger2 {
		t.Errorf("expected %#v to be %#v", logger1, logger2)
	}
}

func TestSlogToLevel(t *testing.T) {
	t.Parallel()

	cases := []struct {
		input string
		want  slog.Level
	}{
		{input: "debug", want: slog.LevelDebug},
		{input: "DEBUG", want: slog.LevelDebug},
		{input: "error", want: slog.LevelError},
		{input: "ERROR", want: slog.LevelError},
		{input: "warn", want: slog.LevelWarn},
		{input: "WARN", want: slog.LevelWarn},
		{input: "info", want: slog.LevelInfo},
		{input: "INFO", want: slog.LevelInfo},
		{input: "other", want: slog.LevelInfo},
		{input: "OTHER", want: slog.LevelInfo},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			t.Parallel()

			got := SlogLevel(tc.input)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}
