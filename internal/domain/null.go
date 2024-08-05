package domain

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

// errNilPtr when the destination pointer is nil.
var errNilPtr = errors.New("destination pointer is nil")

// Null represents a value that may be null. This allows for better handling of
// nullable values while ensuring proper serialization.
//
// If the value is expected to be omitted, add the `omitzero` or `omitempty`
// JSON tags as required. Do NOT use `*Null[T]`.
type Null[T any] struct {
	V     T
	Valid bool
}

// IsZero represents whether a value is null.
func (n Null[T]) IsZero() bool {
	return !n.Valid
}

// MarshalJSONV2 implements the [json.MarshalerV2] interface.
func (n Null[T]) MarshalJSONV2(enc *jsontext.Encoder, opts json.Options) error {
	if n.Valid {
		return json.MarshalEncode(enc, n.V, opts)
	}
	return enc.WriteToken(jsontext.Null)
}

// UnmarshalJSONV2 implements the [json.UnmarshalerV2] interface.
func (n *Null[T]) UnmarshalJSONV2(dec *jsontext.Decoder, opts json.Options) error {
	if dec.PeekKind() == jsontext.Null.Kind() {
		n.V, n.Valid = *new(T), false
		return dec.SkipValue()
	}
	n.Valid = true
	return json.UnmarshalDecode(dec, &n.V, opts)
}

// Scan implements the [sql.Scanner] interface.
func (n *Null[T]) Scan(value any) error {
	if value == nil {
		n.V, n.Valid = *new(T), false
		return nil
	}
	n.Valid = true
	return convertAssign(&n.V, value)
}

// Value implements the [driver.Valuer] interface.
func (n Null[T]) Value() (driver.Value, error) {
	if n.Valid {
		return n.V, nil
	}
	return nil, nil
}

// Undefined represents a value that may be undefined. This allows for better
// handling of undefined values while ensuring proper serialization.
//
// If the value is expected to be omitted, add the `omitzero` or `omitempty`
// JSON tags as required. Do NOT use `*Undefined[T]`.
type Undefined[T any] struct {
	V     any
	Valid bool
}

// IsZero represents whether a value is undefined.
func (u Undefined[T]) IsZero() bool {
	return !u.Valid
}

// MarshalJSONV2 implements the [json.MarshalerV2] interface.
func (u Undefined[T]) MarshalJSONV2(enc *jsontext.Encoder, opts json.Options) error {
	if u.Valid {
		return json.MarshalEncode(enc, u.V, opts)
	}
	return &json.SemanticError{JSONKind: jsontext.Null.Kind(), GoType: reflect.TypeOf(u)}
}

// UnmarshalJSONV2 implements the [json.UnmarshalerV2] interface.
func (u *Undefined[T]) UnmarshalJSONV2(dec *jsontext.Decoder, opts json.Options) error {
	if dec.PeekKind() == jsontext.Null.Kind() {
		u.V, u.Valid = *new(T), false
		return &json.SemanticError{JSONKind: jsontext.Null.Kind(), GoType: reflect.TypeOf(u)}
	}
	u.Valid = true
	return json.UnmarshalDecode(dec, &u.V, opts)
}

// Scan implements the [sql.Scanner] interface.
func (u *Undefined[T]) Scan(value any) error {
	if value == nil {
		u.V, u.Valid = *new(T), false
		return errors.New("converting NULL to undefined is unsupported")
	}
	u.Valid = true
	return convertAssign(&u.V, value)
}

// Value implements the [driver.Valuer] interface.
func (u *Undefined[T]) Value() (driver.Value, error) {
	if u.Valid {
		return u.V, nil
	}
	return nil, errors.New("converting undefined to driver.Value is unsupported")
}

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
	o.v, o.state = v, valid
}

// SetNull sets the value to null.
func (o *Option[T]) SetNull() {
	o.v, o.state = *new(T), null
}

// SetUndefined sets the value to undefined.
func (o *Option[T]) SetUndefined() {
	o.v, o.state = *new(T), undefined
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
	return nil
}

// UnmarshalJSONV2 implements the [json.UnmarshalerV2] interface.
func (o *Option[T]) UnmarshalJSONV2(dec *jsontext.Decoder, opts json.Options) error {
	return nil
}

// Scan implements the [sql.Scanner] interface.
func (o *Option[T]) Scan(v any) error {
	return nil
}

// Value implements the [driver.Valuer] interface.
func (o Option[T]) Value() (driver.Value, error) {
	return nil, nil
}

// convertAssign copies to dest the value in src, converting it if possible.
// An error is returned if the copy would result in loss of information.
// dest should be a pointer type.
//
// This is a slighly modified copy of convertAssign from database/sql.
func convertAssign(dest, src any) error {
	switch s := src.(type) {
	case string:
		switch d := dest.(type) {
		case *string:
			if d == nil {
				return errNilPtr
			}
			*d = s
			return nil
		case *[]byte:
			if d == nil {
				return errNilPtr
			}
			*d = []byte(s)
			return nil
		case *sql.RawBytes:
			if d == nil {
				return errNilPtr
			}
			*d = append((*d)[:0], s...)
			return nil
		}
	case []byte:
		switch d := dest.(type) {
		case *string:
			if d == nil {
				return errNilPtr
			}
			*d = string(s)
			return nil
		case *any:
			if d == nil {
				return errNilPtr
			}
			*d = bytes.Clone(s)
			return nil
		case *[]byte:
			if d == nil {
				return errNilPtr
			}
			*d = bytes.Clone(s)
			return nil
		case *sql.RawBytes:
			if d == nil {
				return errNilPtr
			}
			*d = s
			return nil
		}
	case time.Time:
		switch d := dest.(type) {
		case *time.Time:
			*d = s
			return nil
		case *string:
			*d = s.Format(time.RFC3339Nano)
			return nil
		case *[]byte:
			if d == nil {
				return errNilPtr
			}
			*d = []byte(s.Format(time.RFC3339Nano))
			return nil
		case *sql.RawBytes:
			if d == nil {
				return errNilPtr
			}
			*d = s.AppendFormat((*d)[:0], time.RFC3339Nano)
			return nil
		}
	case nil:
		switch d := dest.(type) {
		case *any:
			if d == nil {
				return errNilPtr
			}
			*d = nil
			return nil
		case *[]byte:
			if d == nil {
				return errNilPtr
			}
			*d = nil
			return nil
		case *sql.RawBytes:
			if d == nil {
				return errNilPtr
			}
			*d = nil
			return nil
		}
	}

	var sv reflect.Value
	switch d := dest.(type) {
	case *string:
		sv = reflect.ValueOf(src)
		switch sv.Kind() {
		case reflect.Bool,
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64:
			*d = asString(src)
			return nil
		}
	case *[]byte:
		sv = reflect.ValueOf(src)
		if b, ok := asBytes(nil, sv); ok {
			*d = b
			return nil
		}
	case *sql.RawBytes:
		sv = reflect.ValueOf(src)
		if b, ok := asBytes([]byte(*d)[:0], sv); ok {
			*d = sql.RawBytes(b)
			return nil
		}
	case *bool:
		bv, err := driver.Bool.ConvertValue(src)
		if err == nil {
			*d = bv.(bool)
		}
		return err
	case *any:
		*d = src
		return nil
	}

	if scanner, ok := dest.(sql.Scanner); ok {
		return scanner.Scan(src)
	}

	dpv := reflect.ValueOf(dest)
	if dpv.Kind() != reflect.Pointer {
		return errors.New("destination not a pointer")
	}
	if dpv.IsNil() {
		return errNilPtr
	}

	if !sv.IsValid() {
		sv = reflect.ValueOf(src)
	}

	dv := reflect.Indirect(dpv)
	if sv.IsValid() && sv.Type().AssignableTo(dv.Type()) {
		switch b := src.(type) {
		case []byte:
			dv.Set(reflect.ValueOf(bytes.Clone(b)))
		default:
			dv.Set(sv)
		}
		return nil
	}

	if dv.Kind() == sv.Kind() && sv.Type().ConvertibleTo(dv.Type()) {
		dv.Set(sv.Convert(dv.Type()))
		return nil
	}

	switch dv.Kind() {
	case reflect.Pointer:
		if src == nil {
			dv.SetZero()
			return nil
		}
		dv.Set(reflect.New(dv.Type().Elem()))
		return convertAssign(dv.Interface(), src)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if src == nil {
			return fmt.Errorf("converting NULL to %s is unsupported", dv.Kind())
		}
		s := asString(src)
		i64, err := strconv.ParseInt(s, 10, dv.Type().Bits())
		if err != nil {
			err = strconvErr(err)
			return fmt.Errorf("converting driver.Value type %T (%q) to a %s: %v", src, s, dv.Kind(), err)
		}
		dv.SetInt(i64)
		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if src == nil {
			return fmt.Errorf("converting NULL to %s is unsupported", dv.Kind())
		}
		s := asString(src)
		u64, err := strconv.ParseUint(s, 10, dv.Type().Bits())
		if err != nil {
			err = strconvErr(err)
			return fmt.Errorf("converting driver.Value type %T (%q) to a %s: %v", src, s, dv.Kind(), err)
		}
		dv.SetUint(u64)
		return nil
	case reflect.Float32, reflect.Float64:
		if src == nil {
			return fmt.Errorf("converting NULL to %s is unsupported", dv.Kind())
		}
		s := asString(src)
		f64, err := strconv.ParseFloat(s, dv.Type().Bits())
		if err != nil {
			err = strconvErr(err)
			return fmt.Errorf("converting driver.Value type %T (%q) to %s: %v", src, s, dv.Kind(), err)
		}
		dv.SetFloat(f64)
		return nil
	case reflect.String:
		if src == nil {
			return fmt.Errorf("converting NULL to %s is unsupported", dv.Kind())
		}
		switch v := src.(type) {
		case string:
			dv.SetString(v)
			return nil
		case []byte:
			dv.SetString(string(v))
			return nil
		}
	}

	return fmt.Errorf("converting type %T to type %T is unsupported", src, dest)
}

// strconvErr returns the reason for why a strconv failed.
func strconvErr(err error) error {
	if ne, ok := err.(*strconv.NumError); ok {
		return ne.Err
	}
	return err
}

// asString converts any type to its string representation.
func asString(src any) string {
	switch v := src.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	}
	rv := reflect.ValueOf(src)
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(rv.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(rv.Uint(), 10)
	case reflect.Float64:
		return strconv.FormatFloat(rv.Float(), 'g', -1, 64)
	case reflect.Float32:
		return strconv.FormatFloat(rv.Float(), 'g', -1, 32)
	case reflect.Bool:
		return strconv.FormatBool(rv.Bool())
	}
	return fmt.Sprintf("%v", src)
}

// asBytes converts a numeric, string, or bool to a byte slice.
func asBytes(buf []byte, rv reflect.Value) (b []byte, ok bool) {
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.AppendInt(buf, rv.Int(), 10), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.AppendUint(buf, rv.Uint(), 10), true
	case reflect.Float32:
		return strconv.AppendFloat(buf, rv.Float(), 'g', -1, 32), true
	case reflect.Float64:
		return strconv.AppendFloat(buf, rv.Float(), 'g', -1, 64), true
	case reflect.Bool:
		return strconv.AppendBool(buf, rv.Bool()), true
	case reflect.String:
		s := rv.String()
		return append(buf, s...), true
	}
	return
}
