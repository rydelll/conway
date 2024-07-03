package domain

import (
	"database/sql"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

// Null represents a value that may be null.
type Null[T any] sql.Null[T]

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
