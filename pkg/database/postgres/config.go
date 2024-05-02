package postgres

import (
	"fmt"
	"net/url"
	"strconv"
)

// Config settings for a database connection.
type Config struct {
	Scheme         string
	Host           string
	Port           int
	Name           string
	User           string
	Password       string
	ConnectTimeout int
	SSLMode        string
	SSLCert        string
	SSLKey         string
	SSLRootCert    string
}

// ConnectionURL builds a connection string suitable for a database.
func (c Config) ConnectionURL() string {
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
	u.RawQuery = q.Encode()

	return u.String()
}
