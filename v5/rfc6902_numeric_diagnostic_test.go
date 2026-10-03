package jsonpatch_test

import (
	"bytes"
	"errors"
	"testing"

	jsonpatch "github.com/evanphx/json-patch/v5"
)

func TestRFC6902NumericDiagnostic(t *testing.T) {
	cases := []struct {
		name, document, value string
		wantSuccess           bool
	}{
		{"integer_decimal", `{"value":1}`, `1.0`, true},
		{"integer_exponent", `{"value":1}`, `1e0`, true},
		{"nested_object", `{"value":{"n":1}}`, `{"n":1.0}`, true},
		{"nested_array", `{"value":[1,{"n":1}]}`, `[1.0,{"n":1e0}]`, true},
		{"identical_number", `{"value":1}`, `1`, true},
		{"unequal_number", `{"value":1}`, `2`, false},
		{"number_vs_string", `{"value":1}`, `"1"`, false},
		{"distinct_above_float64_precision", `{"value":9007199254740992}`, `9007199254740993`, false},
		{"distinct_int64_values", `{"value":9223372036854775807}`, `9223372036854775806`, false},
		{"identical_large_integer", `{"value":9223372036854775807}`, `9223372036854775807`, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			patchJSON := `[{"op":"test","path":"/value","value":` + tc.value + `}]`
			patch, err := jsonpatch.DecodePatch([]byte(patchJSON))
			if err != nil {
				t.Fatalf("DecodePatch setup failed: %v", err)
			}

			result, err := patch.Apply([]byte(tc.document))
			if tc.wantSuccess {
				if err != nil {
					t.Fatalf("numeric equality should pass through Patch.Apply: %v", err)
				}
				if !bytes.Equal(result, []byte(tc.document)) {
					t.Fatalf("test operation changed the document: %s", result)
				}
				return
			}

			if !errors.Is(err, jsonpatch.ErrTestFailed) {
				t.Fatalf("unequal values should return ErrTestFailed, got: %v", err)
			}
			if result != nil {
				t.Fatalf("failed test operation returned a document: %s", result)
			}
		})
	}
}
