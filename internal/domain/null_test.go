package domain

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"strings"
	"testing"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
	"github.com/google/go-cmp/cmp"
)

func TestNullIsZero(t *testing.T) {
	cases := []struct {
		name  string
		input Null[bool]
		want  bool
	}{
		{name: "zero", input: Null[bool]{V: false, Valid: false}, want: true},
		{name: "zero", input: Null[bool]{V: true, Valid: false}, want: true},
		{name: "nonzero", input: Null[bool]{V: true, Valid: true}, want: false},
		{name: "nonzero", input: Null[bool]{V: false, Valid: true}, want: false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.input.IsZero()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestNullMarshalerV2(t *testing.T) {
	v := new(Null[bool])
	var i interface{} = v
	_, ok := i.(json.MarshalerV2)
	if !ok {
		t.Fatal("expected json.MarshalerV2 interface to be satisfied")
	}
}

func TestNullMarshalJSONV2(t *testing.T) {
	cases := []struct {
		name  string
		input Null[bool]
		want  string
	}{
		{name: "valid", input: Null[bool]{V: true, Valid: true}, want: "true"},
		{name: "null", input: Null[bool]{V: false, Valid: false}, want: "null"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			buf := bytes.NewBuffer(nil)
			enc := jsontext.NewEncoder(buf)
			err := tc.input.MarshalJSONV2(enc, json.DefaultOptionsV2())
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			got := strings.TrimSuffix(buf.String(), "\n")
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestNullUnmarshalerV2(t *testing.T) {
	v := new(Null[bool])
	var i interface{} = v
	_, ok := i.(json.UnmarshalerV2)
	if !ok {
		t.Fatal("expected json.UnmarshalerV2 interface to be satisfied")
	}
}

func TestNullUnmarshalJSONV2(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  Null[bool]
	}{
		{name: "valid", input: "true", want: Null[bool]{V: true, Valid: true}},
		{name: "null", input: "null", want: Null[bool]{V: false, Valid: false}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var n Null[bool]
			buf := bytes.NewBuffer([]byte(tc.input))
			dec := jsontext.NewDecoder(buf)
			err := n.UnmarshalJSONV2(dec, json.DefaultOptionsV2())
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if diff := cmp.Diff(tc.want, n); diff != "" {
				t.Errorf("mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestNullScanner(t *testing.T) {
	v := new(Null[bool])
	var i interface{} = v
	_, ok := i.(sql.Scanner)
	if !ok {
		t.Fatal("expected sql.Scanner interface to be satisfied")
	}
}

func TestNullScan(t *testing.T) {
	cases := []struct {
		name  string
		input any
		want  Null[bool]
	}{
		{name: "valid", input: true, want: Null[bool]{V: true, Valid: true}},
		{name: "null", input: nil, want: Null[bool]{V: false, Valid: false}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var n Null[bool]
			err := n.Scan(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if diff := cmp.Diff(tc.want, n); diff != "" {
				t.Errorf("mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestNullValuer(t *testing.T) {
	v := new(Null[bool])
	var i interface{} = v
	_, ok := i.(driver.Valuer)
	if !ok {
		t.Fatal("expected driver.Valuer interface to be satisfied")
	}
}

func TestNullValue(t *testing.T) {
	cases := []struct {
		name  string
		input Null[bool]
		want  driver.Value
	}{
		{name: "valid", input: Null[bool]{V: true, Valid: true}, want: true},
		{name: "null", input: Null[bool]{V: false, Valid: false}, want: nil},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.input.Value()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestUndefinedIsZero(t *testing.T) {
	cases := []struct {
		name  string
		input Undefined[bool]
		want  bool
	}{
		{name: "zero", input: Undefined[bool]{V: false, Valid: false}, want: true},
		{name: "zero", input: Undefined[bool]{V: true, Valid: false}, want: true},
		{name: "nonzero", input: Undefined[bool]{V: true, Valid: true}, want: false},
		{name: "nonzero", input: Undefined[bool]{V: false, Valid: true}, want: false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.input.IsZero()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestUndefinedMarshalerV2(t *testing.T) {
	v := new(Undefined[bool])
	var i interface{} = v
	_, ok := i.(json.MarshalerV2)
	if !ok {
		t.Fatal("expected json.MarshalerV2 interface to be satisfied")
	}
}

func TestUndefinedMarshalJSONV2(t *testing.T) {
	cases := []struct {
		name  string
		input Undefined[bool]
		want  string
		err   bool
	}{
		{name: "valid", input: Undefined[bool]{V: true, Valid: true}, want: "true", err: false},
		{name: "undefined", input: Undefined[bool]{V: false, Valid: false}, want: "", err: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			buf := bytes.NewBuffer(nil)
			enc := jsontext.NewEncoder(buf)
			err := tc.input.MarshalJSONV2(enc, json.DefaultOptionsV2())
			if (err != nil && !tc.err) || (err == nil && tc.err) {
				t.Fatalf("unexpected error: %v", err)
			}

			got := strings.TrimSuffix(buf.String(), "\n")
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestUndefinedUnmarshalerV2(t *testing.T) {
	v := new(Undefined[bool])
	var i interface{} = v
	_, ok := i.(json.UnmarshalerV2)
	if !ok {
		t.Fatal("expected json.UnmarshalerV2 interface to be satisfied")
	}
}

func TestUndefinedUnmarshalJSONV2(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  Undefined[bool]
		err   bool
	}{
		{name: "valid", input: "true", want: Undefined[bool]{V: true, Valid: true}, err: false},
		{name: "undefined", input: "null", want: Undefined[bool]{V: false, Valid: false}, err: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var u Undefined[bool]
			buf := bytes.NewBuffer([]byte(tc.input))
			dec := jsontext.NewDecoder(buf)
			err := u.UnmarshalJSONV2(dec, json.DefaultOptionsV2())
			if (err != nil && !tc.err) || (err == nil && tc.err) {
				t.Fatalf("unexpected error: %v", err)
			}

			if diff := cmp.Diff(tc.want, u); diff != "" {
				t.Errorf("mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestUndefinedScanner(t *testing.T) {
	v := new(Undefined[bool])
	var i interface{} = v
	_, ok := i.(sql.Scanner)
	if !ok {
		t.Fatal("expected sql.Scanner interface to be satisfied")
	}
}

func TestUndefinedScan(t *testing.T) {
	cases := []struct {
		name  string
		input any
		want  Undefined[bool]
		err   bool
	}{
		{name: "valid", input: true, want: Undefined[bool]{V: true, Valid: true}, err: false},
		{name: "null", input: nil, want: Undefined[bool]{V: false, Valid: false}, err: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var u Undefined[bool]
			err := u.Scan(tc.input)
			if (err != nil && !tc.err) || (err == nil && tc.err) {
				t.Fatalf("unexpected error: %v", err)
			}

			if diff := cmp.Diff(tc.want, u); diff != "" {
				t.Errorf("mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestUndefinedValuer(t *testing.T) {
	v := new(Undefined[bool])
	var i interface{} = v
	_, ok := i.(driver.Valuer)
	if !ok {
		t.Fatal("expected driver.Valuer interface to be satisfied")
	}
}

func TestUndefinedValue(t *testing.T) {
	cases := []struct {
		name  string
		input Undefined[bool]
		want  driver.Value
		err   bool
	}{
		{name: "valid", input: Undefined[bool]{V: true, Valid: true}, want: true, err: false},
		{name: "null", input: Undefined[bool]{V: false, Valid: false}, want: nil, err: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.input.Value()
			if (err != nil && !tc.err) || (err == nil && tc.err) {
				t.Fatalf("unexpected error: %v", err)
			}

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestConvertAssign(t *testing.T) {
	// TODO
}

func TestStrconvErr(t *testing.T) {
	// TODO as part of TestConvertAssign
}

func TestAsString(t *testing.T) {
	// TODO as part of TestConvertAssign
}

func TestAsBytes(t *testing.T) {
	// TODO as part of TestConvertAssign
}
