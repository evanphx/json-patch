package jsonpatch_test

import (
	"bytes"
	"errors"
	"testing"

	jsonpatch "github.com/evanphx/json-patch/v5"
)

func TestPatchTestNumericBoundaries(t *testing.T) {
	cases := []struct {
		name, document, value string
		wantSuccess           bool
	}{
		{"negative_zero", `{"value":-0}`, `0`, true},
		{"trailing_zeros", `{"value":1.2300}`, `1.23`, true},
		{"positive_exponent", `{"value":-1.23e+2}`, `-123`, true},
		{"negative_exponent", `{"value":123e-2}`, `1.23`, true},
		{"enormous_scale", `{"value":1e100000000000000000000}`, `10e99999999999999999999`, true},
		{"unequal_scale", `{"value":1e1}`, `1`, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			patch, err := jsonpatch.DecodePatch([]byte(`[{"op":"test","path":"/value","value":` + tc.value + `}]`))
			if err != nil {
				t.Fatalf("DecodePatch failed: %v", err)
			}

			result, err := patch.Apply([]byte(tc.document))
			if tc.wantSuccess {
				if err != nil {
					t.Fatalf("equal numbers should pass: %v", err)
				}
				if !bytes.Equal(result, []byte(tc.document)) {
					t.Fatalf("test operation changed the document: %s", result)
				}
				return
			}

			if !errors.Is(err, jsonpatch.ErrTestFailed) {
				t.Fatalf("unequal numbers should return ErrTestFailed, got: %v", err)
			}
			if result != nil {
				t.Fatalf("failed test operation returned a document: %s", result)
			}
		})
	}
}

func TestEqualNumericValues(t *testing.T) {
	cases := []struct {
		name, left, right string
		wantEqual         bool
	}{
		{"decimal", `1`, `1.0`, true},
		{"recursive", `{"n":[1]}`, `{"n":[1e0]}`, true},
		{"distinct_large_integers", `9007199254740992`, `9007199254740993`, false},
		{"number_vs_string", `"1"`, `1.0`, false},
		{"invalid_leading_zero", `01`, `1`, false},
		{"invalid_exponent", `1e`, `1`, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := jsonpatch.Equal([]byte(tc.left), []byte(tc.right)); got != tc.wantEqual {
				t.Fatalf("Equal(%s, %s) = %v, want %v", tc.left, tc.right, got, tc.wantEqual)
			}
		})
	}
}
