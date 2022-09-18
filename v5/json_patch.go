package jsonpatch

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
)

// CreateJsonPatch will return a JSON patch document capable of converting
// the original document(s) to the modified document(s).
// The parameters can be bytes of either two JSON Documents, or two arrays of
// JSON documents.
// The JSON patch returned follows the specification defined at https://www.rfc-editor.org/rfc/rfc6902
func CreateJsonPatch(originalJSON, modifiedJSON []byte) (*Patch, error) {
	originalType := testJsonType(originalJSON)
	modifiedType := testJsonType(modifiedJSON)
	patch := new(Patch)

	// Do both byte-slices seem like JSON arrays?
	if originalType == JSONArray && modifiedType == JSONArray {
		if err := createArrayJsonPatch(originalJSON, modifiedJSON, patch, ""); err != nil {
			return nil, err
		}
		return patch, nil
	}

	// Are both byte-slices are not arrays? Then they are likely JSON objects...
	if originalType == JSONObject && modifiedType == JSONObject {
		if err := createObjectJsonPatch(originalJSON, modifiedJSON, patch, ""); err != nil {
			return nil, err
		}
		return patch, nil
	}

	// replace root
	var root interface{}
	if err := json.Unmarshal(modifiedJSON, &root); err != nil {
		return nil, err
	}
	if o, err := NewReplace("", root); err != nil {
		return nil, err
	} else {
		patch.AddOperation(o)
	}

	return patch, nil
}

// createArrayJsonPatch will return an array of json-patch documents capable
// of converting the original document to the modified document for each
// pair of JSON documents provided in the arrays.
// Arrays of mismatched sizes will result in an error.
func createArrayJsonPatch(originalJSON, modifiedJSON []byte, patch *Patch, path string) error {
	originalDocs := []json.RawMessage{}
	modifiedDocs := []json.RawMessage{}

	err := json.Unmarshal(originalJSON, &originalDocs)
	if err != nil {
		return errBadJSONDoc
	}

	err = json.Unmarshal(modifiedJSON, &modifiedDocs)
	if err != nil {
		return errBadJSONDoc
	}

	minLen := int(math.Min(float64(len(originalDocs)), float64(len(modifiedDocs))))
	for i := 0; i < minLen; i++ {
		original := originalDocs[i]
		modified := modifiedDocs[i]

		if err := createValueJsonPatch(original, modified, patch, fmt.Sprintf("%s/%d", path, i)); err != nil {
			return err
		}
	}

	if minLen < len(modifiedDocs) {
		for i := minLen; i < len(modifiedDocs); i++ {
			if o, err := NewAdd(fmt.Sprintf("%s/%d", path, i), modifiedDocs[i]); err != nil {
				return err
			} else {
				patch.AddOperation(o)
			}
		}
	}

	if minLen < len(originalDocs) {
		for i := minLen; i < len(originalDocs); i++ {
			if o, err := NewRemove(fmt.Sprintf("%s/%d", path, i)); err != nil {
				return err
			} else {
				patch.AddOperation(o)
			}
		}
	}

	return nil
}

// createObjectJsonPatch will return a json-patch document capable of
// converting the original document to the modified document.
func createObjectJsonPatch(originalJSON, modifiedJSON []byte, patch *Patch, path string) error {
	originalDoc := map[string]interface{}{}
	modifiedDoc := map[string]interface{}{}

	err := json.Unmarshal(originalJSON, &originalDoc)
	if err != nil {
		return errBadJSONDoc
	}

	err = json.Unmarshal(modifiedJSON, &modifiedDoc)
	if err != nil {
		return errBadJSONDoc
	}

	return buildJsonDiff(originalDoc, modifiedDoc, patch, path)
}

// createValueJsonPatch will return a json-patch document capable of
// converting the original document to the modified document.
func createValueJsonPatch(originJSON, modifiedJSON json.RawMessage, patch *Patch, path string) error {
	var originalDoc interface{}
	var modifiedDoc interface{}

	err := json.Unmarshal(originJSON, &originalDoc)
	if err != nil {
		fmt.Println(err)
		return errBadJSONDoc
	}

	err = json.Unmarshal(modifiedJSON, &modifiedDoc)
	if err != nil {
		return errBadJSONDoc
	}

	if originalDoc != nil && modifiedDoc == nil {
		if o, err := NewRemove(path); err != nil {
			return err
		} else {
			patch.AddOperation(o)
		}
		return nil
	}
	if originalDoc == nil && modifiedDoc != nil {
		if o, err := NewAdd(path, modifiedDoc); err != nil {
			return err
		} else {
			patch.AddOperation(o)
		}
		return nil
	}
	at := reflect.TypeOf(originalDoc).Kind()
	if at != reflect.TypeOf(modifiedDoc).Kind() {
		if o, err := NewReplace(path, modifiedDoc); err != nil {
			return err
		} else {
			patch.AddOperation(o)
		}
		return nil
	}
	switch at {
	case reflect.Array, reflect.Slice:
		return createArrayJsonPatch(originJSON, modifiedJSON, patch, path)
	case reflect.Map:
		return buildJsonDiff(originalDoc.(map[string]interface{}), modifiedDoc.(map[string]interface{}), patch, path)
	case reflect.String, reflect.Float64, reflect.Bool:
		if !matchesValue(originalDoc, modifiedDoc) {
			if o, err := NewReplace(path, modifiedDoc); err != nil {
				return err
			} else {
				patch.AddOperation(o)
			}
		}
	default:
		return fmt.Errorf("Unknown type:%T in key %s", originalDoc, path)
	}
	return nil
}

// buildJsonDiff returns the (recursive) difference between a and b as a patch rule in Patch.
func buildJsonDiff(a, b map[string]interface{}, patch *Patch, path string) error {
	for key, bv := range b {
		av, ok := a[key]
		// value was added
		if !ok {
			if o, err := NewAdd(path+"/"+key, bv); err != nil {
				return err
			} else {
				patch.AddOperation(o)
			}
			continue
		}
		// If types have changed, replace completely
		if reflect.TypeOf(av) != reflect.TypeOf(bv) {
			if o, err := NewReplace(path+"/"+key, bv); err != nil {
				return err
			} else {
				patch.AddOperation(o)
			}
			continue
		}
		// Types are the same, compare values
		switch at := av.(type) {
		case map[string]interface{}:
			if err := buildJsonDiff(av.(map[string]interface{}), bv.(map[string]interface{}), patch, path+"/"+key); err != nil {
				return err
			}
		case string, float64, bool:
			if !matchesValue(av, bv) {
				if o, err := NewReplace(path+"/"+key, bv); err != nil {
					return err
				} else {
					patch.AddOperation(o)
				}
			}
		case []interface{}:
			bt := bv.([]interface{})
			if !matchesArray(at, bt) {
				if o, err := NewReplace(path+"/"+key, bv); err != nil {
					return err
				} else {
					patch.AddOperation(o)
				}
			}
		case nil:
			switch bv.(type) {
			case nil:
				// Both nil, fine.
			default:
				if o, err := NewReplace(path+"/"+key, bv); err != nil {
					return err
				} else {
					patch.AddOperation(o)
				}
			}
		default:
			return fmt.Errorf("Unknown type:%T in key %s", av, path+"/"+key)
		}
	}
	// Now add all deleted values as nil
	for key := range a {
		_, found := b[key]
		if !found {
			if o, err := NewRemove(path + "/" + key); err != nil {
				return err
			} else {
				patch.AddOperation(o)
			}
		}
	}
	return nil
}
