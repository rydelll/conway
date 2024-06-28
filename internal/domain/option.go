package domain

import (
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

// state represents if an [Option] is valid, null, or undefined.
type state byte

const (
	valid state = iota
	null
	undefined
)

// Option is a generic type, which implements a value that can be in one of
// three states: T, null, or undefined.
//
// If the field is expected to be optional, add the JSON tags `omitzero`
// or `omitempty` as required. Do NOT use *Option[T].
type Option[T any] struct {
	v     T
	state state
}

// NewOption creates an [Option] with a valid value.
func NewOption[T any](v T) Option[T] {
	return Option[T]{
		v:     v,
		state: valid,
	}
}

// NewOptionNull creates an [Option] of T that is initialized to null.
func NewOptionNull[T any]() Option[T] {
	return Option[T]{
		v:     *new(T),
		state: null,
	}
}

// NewOptionUndefined creates an [Option] of T that is initialized to undefined.
func NewOptionUndefined[T any]() Option[T] {
	return Option[T]{
		v:     *new(T),
		state: undefined,
	}
}

// Get the value, or error if the value is null or undefined.
func (o Option[T]) Get() (T, error) {
	switch o.state {
	case valid:
		return o.v, nil
	case null:
		return o.v, ErrNull
	default:
		return o.v, ErrUndefined
	}
}

// MustGet the value, or panic if the value is null or undefined.
func (o Option[T]) MustGet() T {
	switch o.state {
	case valid:
		return o.v
	case null:
		panic(ErrNull)
	default:
		panic(ErrUndefined)
	}
}

// Set a valid value.
func (o *Option[T]) Set(v T) {
	o.v = v
	o.state = valid
}

// SetNull sets the value to null.
func (o *Option[T]) SetNull() {
	o.v = *new(T)
	o.state = null
}

// SetUndefined sets the value to undefined.
func (o *Option[T]) SetUndefined() {
	o.v = *new(T)
	o.state = undefined
}

// IsNull indicates if the value is null.
func (o Option[T]) IsNull() bool {
	return o.state == null
}

// IsUndefined indicates if the value is undefined.
func (o Option[T]) IsUndefined() bool {
	return o.state == undefined
}

// IsZero implements the [json.isZeroer] interface.
func (o Option[T]) IsZero() bool {
	return o.IsUndefined()
}

// MarshalJSONV2 implements the [json.MarshalerV2] interface.
func (o Option[T]) MarshalJSONV2(enc *jsontext.Encoder, opts json.Options) error {
	switch o.state {
	case valid:
		return json.MarshalEncode(enc, o.v, opts)
	default:
		return enc.WriteToken(jsontext.Null)
	}
}

// UnmarshalJSONV2 implements the [json.UnmarshalerV2] interface.
func (o *Option[T]) UnmarshalJSONV2(dec *jsontext.Decoder, opts json.Options) error {
	if dec.PeekKind() == jsontext.Null.Kind() {
		o.SetNull()
		return dec.SkipValue()
	}
	// TODO: do I need to handle undefined case?
	o.state = valid
	return json.UnmarshalDecode(dec, &o.v, opts)
}

// // Value implements the [sql.Valuer] interface.
// func (o Option[T]) Value() (driver.Value, error) {

// }

// // Scan implements the [sql.Scanner] interface.
// func (o *Option[T]) Scan(v any) error {

// }
