package domain

import (
	"database/sql"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

type Null[T any] sql.Null[T]

func (n Null[T]) MarshalJSONV2(enc *jsontext.Encoder, opts json.Options) error {
	if n.Valid {
		return json.MarshalEncode(enc, n.V, opts)
	}
	return enc.WriteToken(jsontext.Null)
}

func (n *Null[T]) UnmarshalJSONV2(dec *jsontext.Decoder, opts json.Options) error {
	if dec.PeekKind() == jsontext.Null.Kind() {
		n.Valid = false
		return dec.SkipValue()
	}
	n.Valid = true
	return json.UnmarshalDecode(dec, &n.V, opts)
}
