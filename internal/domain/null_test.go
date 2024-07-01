package domain

import (
	"bytes"
	"testing"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

func TestNullMarshal(t *testing.T) {
	cases := []struct {
		name  string
		input Null[int]
		want  []byte
	}{
		{name: "valid", input: Null[int]{V: 69, Valid: true}, want: []byte("69\n")},
		{name: "null", input: Null[int]{V: 69, Valid: false}, want: []byte("null\n")},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			buf := bytes.NewBuffer(nil)
			enc := jsontext.NewEncoder(buf)
			err := tc.input.MarshalJSONV2(enc, json.DefaultOptionsV2())
			if err != nil {
				t.Errorf("unexpected marshal error: %v", err)
			}

			out := buf.Bytes()
			if !bytes.Equal(out, tc.want) {
				t.Fatalf("mismatch (want, got):\n%s, %s", tc.want, out)
			}
		})
	}
}

func TestNullUnmarshal(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		input []byte
		value int
		valid bool
	}{
		{name: "valid", input: []byte("69"), value: 69, valid: true},
		{name: "null", input: []byte("null"), value: 0, valid: false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var n Null[int]
			buf := bytes.NewBuffer(tc.input)
			dec := jsontext.NewDecoder(buf)
			err := n.UnmarshalJSONV2(dec, json.DefaultOptionsV2())
			if err != nil {
				t.Errorf("unexpected unmarshal error: %v", err)
			}

			if n.V != tc.value {
				t.Errorf("mismatch (want, got):\n%d, %d", tc.value, n.V)
			}
			if n.Valid != tc.valid {
				t.Errorf("mismatch (want, got):\n%t, %t", tc.valid, n.Valid)
			}
		})
	}
}
