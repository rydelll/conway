package domain

import "errors"

var (
	// ErrInternal when an unknown internal service fails.
	ErrInternal = errors.New("internal server error")
	// ErrConflict when an a unique constrains on the data is violated.
	ErrConflict = errors.New("unique data conflict")
	// ErrNotFound when the requested resource is not found.
	ErrNotFound = errors.New("not found")
	// ErrNoUpdate when no data is provider for an update.
	ErrNoUpdate = errors.New("no update data")
	// ErrNull when an option is null.
	ErrNull = errors.New("option is null")
	// ErrUndefined when an option is undefined.
	ErrUndefined = errors.New("option is undefined")
)
