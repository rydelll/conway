package database

import (
	"context"
	"testing"
	"time"
)

func TestNewPostgres(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		config PGConfig
	}{
		{
			name: "all",
			config: PGConfig{
				Scheme:             "postgres",
				Host:               "localhost",
				Port:               1234,
				Name:               "database",
				User:               "user",
				Password:           "password",
				ConnectTimeout:     10,
				SSLMode:            "disable",
				SSLCert:            "db.crt",
				SSLKey:             "db.key",
				SSLRootCert:        "root.crt",
				PoolMinConnections: 2,
				PoolMaxConnections: 10,
				PoolMaxConnLife:    time.Minute * 5,
				PoolMaxConnIdle:    time.Minute * 2,
				PoolHealthcheck:    time.Minute,
			},
		},
		{
			name:   "empty",
			config: PGConfig{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			pool, err := NewPostgres(context.Background(), tc.config)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if pool == nil {
				t.Fatal("expected pool to never be nil")
			}
		})
	}
}

func TestNewPostgresErr(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		config PGConfig
	}{
		{
			name: "missing",
			config: PGConfig{
				Host: "localhost",
				Name: "database",
			},
		},
		{
			name: "sslrequire",
			config: PGConfig{
				Scheme:      "postgres",
				Host:        "localhost",
				Name:        "database",
				SSLMode:     "require",
				SSLCert:     "db.crt",
				SSLKey:      "db.key",
				SSLRootCert: "root.crt",
			},
		},
		{
			name: "sslcert",
			config: PGConfig{
				Scheme:  "postgres",
				Host:    "localhost",
				Name:    "database",
				SSLCert: "db.crt",
			},
		},
		{
			name: "poolsize",
			config: PGConfig{
				Scheme:             "postgres",
				Host:               "localhost",
				Name:               "database",
				PoolMaxConnections: 999999999999999,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewPostgres(context.Background(), tc.config)
			if err == nil {
				t.Fatal("expected an error")
			}
		})
	}
}
