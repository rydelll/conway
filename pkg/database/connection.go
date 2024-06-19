package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPostgres creates a PostgreSQL database handle that uses a pool of
// connections to allow for safe concurrent use and verifies a connection
// can be established.
func NewPostgres(ctx context.Context, config PGConfig) (*pgxpool.Pool, error) {
	pgxConfig, err := pgxpool.ParseConfig(config.ConnectionURL())
	if err != nil {
		return nil, fmt.Errorf("parse connection string: %v", err)
	}

	pgxConfig.BeforeAcquire = func(ctx context.Context, conn *pgx.Conn) bool {
		return conn.Ping(ctx) == nil
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		return nil, fmt.Errorf("establish connection: %v", err)
	}

	return pool, nil
}
