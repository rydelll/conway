package server

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestWithReadTimeout(t *testing.T) {
	cases := []struct {
		time time.Duration
	}{
		{time: time.Millisecond},
		{time: time.Second},
		{time: time.Minute},
	}

	for _, tc := range cases {
		t.Run(tc.time.String(), func(t *testing.T) {
			srv := new(Server)
			srv.server = new(http.Server)
			WithReadTimeout(tc.time)(srv)
			if diff := cmp.Diff(tc.time, srv.server.ReadTimeout); diff != "" {
				t.Errorf("mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestWithWriteTimeout(t *testing.T) {
	cases := []struct {
		time time.Duration
	}{
		{time: time.Millisecond},
		{time: time.Second},
		{time: time.Minute},
	}

	for _, tc := range cases {
		t.Run(tc.time.String(), func(t *testing.T) {
			srv := new(Server)
			srv.server = new(http.Server)
			WithWriteTimeout(tc.time)(srv)
			if diff := cmp.Diff(tc.time, srv.server.WriteTimeout); diff != "" {
				t.Errorf("mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestWithIdleTimeout(t *testing.T) {
	cases := []struct {
		time time.Duration
	}{
		{time: time.Millisecond},
		{time: time.Second},
		{time: time.Minute},
	}

	for _, tc := range cases {
		t.Run(tc.time.String(), func(t *testing.T) {
			srv := new(Server)
			srv.server = new(http.Server)
			WithIdleTimeout(tc.time)(srv)
			if diff := cmp.Diff(tc.time, srv.server.IdleTimeout); diff != "" {
				t.Errorf("mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestWithShutdownTimeout(t *testing.T) {
	cases := []struct {
		time time.Duration
	}{
		{time: time.Millisecond},
		{time: time.Second},
		{time: time.Minute},
	}

	for _, tc := range cases {
		t.Run(tc.time.String(), func(t *testing.T) {
			srv := new(Server)
			WithShutdownTimeout(tc.time)(srv)
			if diff := cmp.Diff(tc.time, srv.shutdownTimeout); diff != "" {
				t.Errorf("mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}
