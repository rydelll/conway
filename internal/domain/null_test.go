package domain

import (
	"bytes"
	"testing"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
	"github.com/google/go-cmp/cmp"
)

func TestNullMarshal(t *testing.T) {
	cases := []struct {
		name  string
		input Null[int]
		want  string
	}{
		{name: "valid", input: Null[int]{V: 69, Valid: true}, want: "69\n"},
		{name: "null", input: Null[int]{V: 69, Valid: false}, want: "null\n"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			buf := bytes.NewBuffer(nil)
			enc := jsontext.NewEncoder(buf)
			err := tc.input.MarshalJSONV2(enc, json.DefaultOptionsV2())
			if err != nil {
				t.Fatalf("unexpected marshal error: %v", err)
			}

			got := buf.String()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("mismatch (-want, +got):\n%s", diff)
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
				t.Fatalf("unexpected unmarshal error: %v", err)
			}

			if diff := cmp.Diff(tc.value, n.V); diff != "" {
				t.Errorf("mismatch (-want, +got):\n%s", diff)
			}
			if diff := cmp.Diff(tc.valid, n.Valid); diff != "" {
				t.Errorf("mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}
