package database

import (
	"fmt"
	"log/slog"
	"net/url"
	"strconv"
	"time"
)

// PGConfig holds the information required for a PostgreSQL database connection.
type PGConfig struct {
	Scheme             string
	Host               string
	Port               int
	Name               string
	User               string
	Password           string
	ConnectTimeout     int
	SSLMode            string
	SSLCert            string
	SSLKey             string
	SSLRootCert        string
	PoolMinConnections int
	PoolMaxConnections int
	PoolMaxConnLife    time.Duration
	PoolMaxConnIdle    time.Duration
	PoolHealthcheck    time.Duration
}

// ConnectionURL builds a connection string suitable for a database.
func (c PGConfig) ConnectionURL() string {
	host := c.Host
	if v := c.Port; v > 0 {
		host = fmt.Sprintf("%s:%d", host, v)
	}

	u := url.URL{
		Scheme: c.Scheme,
		Host:   host,
		Path:   c.Name,
	}

	if c.User != "" || c.Password != "" {
		u.User = url.UserPassword(c.User, c.Password)
	}

	q := u.Query()
	if v := c.ConnectTimeout; v > 0 {
		q.Add("connect_timeout", strconv.Itoa(v))
	}
	if v := c.SSLMode; v != "" {
		q.Add("sslmode", v)
	}
	if v := c.SSLCert; v != "" {
		q.Add("sslcert", v)
	}
	if v := c.SSLKey; v != "" {
		q.Add("sslkey", v)
	}
	if v := c.SSLRootCert; v != "" {
		q.Add("sslrootcert", v)
	}
	if v := c.PoolMinConnections; v > 0 {
		q.Add("pool_min_conns", strconv.Itoa(v))
	}
	if v := c.PoolMaxConnections; v > 0 {
		q.Add("pool_max_conns", strconv.Itoa(v))
	}
	if v := c.PoolMaxConnLife; v > 0 {
		q.Add("pool_max_conn_lifetime", v.String())
	}
	if v := c.PoolMaxConnIdle; v > 0 {
		q.Add("pool_max_conn_idle_time", v.String())
	}
	if v := c.PoolHealthcheck; v > 0 {
		q.Add("pool_health_check_period", v.String())
	}
	u.RawQuery = q.Encode()

	return u.String()
}

func (c PGConfig) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("Scheme", c.Scheme),
		slog.String("Host", c.Host),
		slog.Int("Port", c.Port),
		slog.String("Name", c.Name),
		slog.String("User", c.User),
		slog.Int("ConnectTimeout", c.ConnectTimeout),
		slog.String("SSLMode", c.SSLMode),
		slog.String("SSLCert", c.SSLCert),
		slog.String("SSLKey", c.SSLKey),
		slog.String("SSLRootCert", c.SSLRootCert),
		slog.Int("PoolMinConnections", c.PoolMinConnections),
		slog.Int("PoolMaxConnections", c.PoolMaxConnections),
		slog.Duration("PoolMaxConnLife", c.PoolMaxConnLife),
		slog.Duration("PoolMaxConnIdle", c.PoolMaxConnIdle),
		slog.Duration("PoolHealthcheck", c.PoolHealthcheck),
	)
}
