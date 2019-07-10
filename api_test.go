package jsonpatch

import (
	stdlib "encoding/json"
	"strings"
	"testing"
)

type customJsonApi struct {
}

func (api customJsonApi) Marshal(v interface{}) ([]byte, error) {
	output, err := stdlib.Marshal(v)
	if err != nil {
		return nil, err
	}
	stringOutput := string(output)
	if stringOutput != "" && strings.Contains(stringOutput, "bar") {
		newString := strings.Replace(stringOutput, "bar", "BAR", -1)
		return []byte(newString), nil
	}
	return output, nil
}

func (api customJsonApi) Unmarshal(data []byte, v interface{}) error {
	return stdlib.Unmarshal(data, v)
}

func (api customJsonApi) MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return stdlib.MarshalIndent(v, prefix, indent)
}

func TestMergePatchWithCustomMarshal(t *testing.T) {

	SetAPI(customJsonApi{})

	doc := `{ "hello": "goodbye", "foo": "bar" }`
	patch := `{ "bed": "time", "good": "night" }`
	expectedOutput := `{"bed":"time","foo":"BAR","good":"night","hello":"goodbye"}`

	out, err := MergePatch([]byte(doc), []byte(patch))

	if err != nil {
		t.Fatalf("Error performing merge patch")
	} else {
		result := string(out)
		if result != expectedOutput {
			t.Errorf(
				`expected "%s" but received "%s"`,
				expectedOutput,
				result,
			)
		}
	}

	// Verify that we can reset to using the standard API
	ResetAPI()

	expectedOutput = `{"bed":"time","foo":"bar","good":"night","hello":"goodbye"}`

	out, err = MergePatch([]byte(doc), []byte(patch))

	if err != nil {
		t.Fatalf("Error performing merge patch")
	} else {
		result := string(out)
		if result != expectedOutput {
			t.Errorf(
				`expected "%s" but received "%s"`,
				expectedOutput,
				result,
			)
		}
	}

}
