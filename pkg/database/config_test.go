package database

import (
	"bytes"
	"log/slog"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/rydelll/conway/pkg/logging"
)

func TestConnectionURL(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		config PGConfig
		want   string
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
				SSLMode:            "require",
				SSLCert:            "db.crt",
				SSLKey:             "db.key",
				SSLRootCert:        "root.crt",
				PoolMinConnections: 2,
				PoolMaxConnections: 10,
				PoolMaxConnLife:    time.Minute * 5,
				PoolMaxConnIdle:    time.Minute * 2,
				PoolHealthcheck:    time.Minute,
			},
			want: "postgres://user:password@localhost:1234/database?connect_timeout=10&" +
				"pool_health_check_period=1m0s&pool_max_conn_idle_time=2m0s&" +
				"pool_max_conn_lifetime=5m0s&pool_max_conns=10&pool_min_conns=2&" +
				"sslcert=db.crt&sslkey=db.key&sslmode=require&sslrootcert=root.crt",
		},
		{
			name: "nopassword",
			config: PGConfig{
				Scheme: "postgres",
				Host:   "localhost",
				User:   "user",
			},
			want: "postgres://user@localhost",
		},
		{
			name: "nouser",
			config: PGConfig{
				Scheme:   "postgres",
				Host:     "localhost",
				Password: "password",
			},
			want: "postgres://localhost",
		},
		{
			name:   "empty",
			config: PGConfig{},
			want:   "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := tc.config.ConnectionURL()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestLogValuer(t *testing.T) {
	t.Parallel()

	v := new(PGConfig)
	var i interface{} = v
	_, ok := i.(slog.LogValuer)
	if !ok {
		t.Fatal("expected slog.LogValuer interface to be satisfied")
	}
}

func TestLogValue(t *testing.T) {
	t.Parallel()

	want := "{\"level\":\"INFO\",\"msg\":\"test\",\"config\":{\"Scheme\":\"postgres\"," +
		"\"Host\":\"localhost\",\"Port\":1234,\"Name\":\"database\",\"User\":\"user\"," +
		"\"Password\":\"[REDACTED]\",\"ConnectTimeout\":10,\"SSLMode\":\"require\"," +
		"\"SSLCert\":\"db.crt\",\"SSLKey\":\"db.key\",\"SSLRootCert\":\"root.crt\"," +
		"\"PoolMinConnections\":2,\"PoolMaxConnections\":10,\"PoolMaxConnLife\":300000000000," +
		"\"PoolMaxConnIdle\":120000000000,\"PoolHealthcheck\":60000000000}}\n"

	buf := bytes.NewBuffer(nil)
	logger := logging.NewLoggerTimeless(buf, slog.LevelInfo, true)
	config := PGConfig{
		Scheme:             "postgres",
		Host:               "localhost",
		Port:               1234,
		Name:               "database",
		User:               "user",
		Password:           "password",
		ConnectTimeout:     10,
		SSLMode:            "require",
		SSLCert:            "db.crt",
		SSLKey:             "db.key",
		SSLRootCert:        "root.crt",
		PoolMinConnections: 2,
		PoolMaxConnections: 10,
		PoolMaxConnLife:    time.Minute * 5,
		PoolMaxConnIdle:    time.Minute * 2,
		PoolHealthcheck:    time.Minute,
	}

	logger.Info("test", slog.Any("config", config))
	if diff := cmp.Diff(want, buf.String()); diff != "" {
		t.Errorf("mismatch (-want, +got):\n%s", diff)
	}
}
