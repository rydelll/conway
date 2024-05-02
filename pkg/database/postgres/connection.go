package postgres

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	defaultMaxIdleConns     = 5
	defaultMaxOpenConns     = 25
	defaultConnsMaxIdleTime = time.Minute * 5
	defaultConnsMaxLifetime = time.Millisecond * 100
)

// DB is a database handle representing a pool of zero or more underlying connections.
// It is safe for concurrent use by multiple goroutines.
type DB struct {
	*sql.DB
}

// New creates a database handle that is safe for concurrent use.
func New(config Config) (*DB, error) {
	db, err := sql.Open("pgx", config.ConnectionURL())
	if err != nil {
		return nil, fmt.Errorf("parse database connection string: %v", err)
	}

	db.SetMaxIdleConns(defaultMaxIdleConns)
	db.SetMaxOpenConns(defaultMaxOpenConns)
	db.SetConnMaxIdleTime(defaultConnsMaxIdleTime)
	db.SetConnMaxLifetime(defaultConnsMaxLifetime)

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("database connection ping: %v", err)
	}

	return &DB{db}, nil
}

// NewFromEnv creates a database handle from environment variables.
// The handle is safe for concurrent use.
//
// Consumes DB_SCHEME, DB_HOST, DB_PORT, DB_NAME, DB_USER, DB_PASSWORD,
// DB_CONNECT_TIMEOUT, DB_SSLMODE, DB_SSLCERT, DB_SSLKEY, and DB_SSLROOTCERT
// for determining that database connection configuration.
func NewFromEnv() (*DB, error) {
	config := Config{
		Scheme:      os.Getenv("DB_SCHEME"),
		Host:        os.Getenv("DB_HOST"),
		Name:        os.Getenv("DB_NAME"),
		User:        os.Getenv("DB_USER"),
		Password:    os.Getenv("DB_PASSWORD"),
		SSLMode:     os.Getenv("DB_SSLMODE"),
		SSLCert:     os.Getenv("DB_SSLCERT"),
		SSLKey:      os.Getenv("DB_SSLKEY"),
		SSLRootCert: os.Getenv("DB_SSLROOTCERT"),
	}

	// default to zero by ignoring parsing errors
	config.Port, _ = strconv.Atoi(os.Getenv("DB_PORT"))
	config.ConnectTimeout, _ = strconv.Atoi(os.Getenv("DB_CONNECT_TIMEOUT"))

	return New(config)
}
