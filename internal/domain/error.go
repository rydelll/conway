package domain

import "errors"

var (
	// ErrInternal when an unknown internal service fails.
	ErrInternal = errors.New("internal server error")
	// ErrInvalidID when a resource ID is unable to be parsed.
	ErrInvalidID = errors.New("invalid resource id")
	// ErrNotFound when the requested resource is not found.
	ErrNotFound = errors.New("not found")
	// ErrConflict when data conflicts with the current state of the server.
	ErrConflict = errors.New("data conflict")
	// ErrNoUpdate when no data is provider for an update.
	ErrNoUpdate = errors.New("no update data")
	// ErrNull when an option is null.
	ErrNull = errors.New("option is null")
	// ErrUndefined when an option is undefined.
	ErrUndefined = errors.New("option is undefined")
)
