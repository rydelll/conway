package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Database represents a PostgreSQL database connection pool contract.
type Database interface {
	// Begin acquires a connection from the Pool and starts a transaction.
	// Unlike database/sql, the context only affects the begin command. i.e.
	// there is no auto-rollback on context cancellation. Begin initiates a
	// transaction block without explicitly setting a transaction mode for the
	// block (see BeginTx with TxOptions if transaction mode is required).
	// d*pgxpool.Tx is returned, which implements the pgx.Tx interface. Commit
	// or Rollback must be called on the returned transaction to finalize the
	// transaction block.
	Begin(ctx context.Context) (pgx.Tx, error)
	// Exec acquires a connection from the Pool and executes the given SQL.
	// SQL can be either a prepared statement name or an SQL string. Arguments
	// should be positional from the SQL string as $1, $2, etc. The acquired
	// connection is returned to the pool when the Exec function returns.
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	// Query acquires a connection and executes a query that returns pgx.Rows.
	// Arguments should be positional from the SQL string as $1, $2, etc. See
	// pgx.Rows documentation to close the returned Rows and return the acquired
	// connection to the Pool.
	//
	// If there is an error, the returned pgx.Rows will be returned in an error
	// state. If preferred, ignore the error returned from Query and handle
	// errors using the returned pgx.Rows.
	//
	// For extra control over how the query is executed, the types
	// QuerySimpleProtocol, QueryResultFormats, and QueryResultFormatsByOID may
	// be used as the first args to control exactly how the query is executed.
	// This is rarely needed. See the documentation for those types for details.
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	// QueryRow acquires a connection and executes a query that is expected to
	// return at most one row (pgx.Row). Errors are deferred until pgx.Row's
	// Scan method is called. If the query selects no rows, pgx.Row's Scan will
	// return ErrNoRows. Otherwise, pgx.Row's Scan scans the first selected row
	// and discards the rest. The acquired connection is returned to the Pool
	// when pgx.Row's Scan method is called.
	//
	// Arguments should be positional from the SQL string as $1, $2, etc.
	//
	// For extra control over how the query is executed, the types
	// QuerySimpleProtocol, QueryResultFormats, and QueryResultFormatsByOID may
	// be used as the first args to control exactly how the query is executed.
	// This is rarely needed. See the documentation for those types for details.
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}
