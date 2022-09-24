package jsonpatch

import (
	"fmt"
	"testing"
)

func TestObjectJsonPatchExample(t *testing.T) {
	// Let's create a merge patch from these two documents...
	original := []byte(`{"name": "John", "age": 24, "height": 3.21}`)
	target := []byte(`{"name": "Jane", "age": 24, "class": 2217, "no": "33", "position": {"x":5,"y":2}}`)

	patch, err := CreateJsonPatch(original, target)
	if err != nil {
		t.Fatal(err)
	}

	// Now lets apply the patch against a different JSON document...

	alternative := []byte(`{"name": "Tina", "age": 28, "height": 3.75}`)
	modifiedAlternative, err := patch.Apply(alternative)
	if err != nil {
		panic(err)
	}

	fmt.Printf("patch document:   %s\n", patch)
	fmt.Printf("updated alternative doc: %s\n", modifiedAlternative)
}

func TestArrayJsonPatchExample(t *testing.T) {
	// Let's create a merge patch from these two documents...
	original := []byte(`["name",1,{"x":"x1"},[1]]`)
	target := []byte(`["name value",2,{"x":"x2","y":"y2"},[0,1]]`)

	patch, err := CreateJsonPatch(original, target)
	if err != nil {
		t.Fatal(err)
	}

	// Now lets apply the patch against a different JSON document...

	alternative := []byte(`["name",1,{"x":"x1","z":"z1"},[1,2,3,4,5]]`)
	modifiedAlternative, err := patch.Apply(alternative)
	if err != nil {
		panic(err)
	}

	fmt.Printf("patch document:   %s\n", patch)
	fmt.Printf("updated alternative doc: %s\n", modifiedAlternative)
}

func jsonPatch(doc, modified string) (string, error) {
	original := []byte(doc)
	target := []byte(modified)

	patch, err := CreateJsonPatch(original, target)
	if err != nil {
		return "", err
	}
	out, err := patch.ApplyWithOptions(original, &ApplyOptions{
		AllowMissingPathOnRemove: true,
	})
	if err != nil {
		return "", err
	}
	fmt.Printf("\noriginal:%s\nmodified:%s\nactual:%s\njson patch:%s\n", doc, modified, string(out), patch.String())
	return string(out), nil
}

func TestJsonPatchReplaceKey(t *testing.T) {
	doc := `{ "title": "hello" }`
	pat := `{ "title": "goodbye" }`

	res, err := jsonPatch(doc, pat)
	if err != nil {
		panic(err)
	}

	if !compareJSON(pat, res) {
		t.Fatalf("Key was not replaced")
	}
}

func TestJsonPatchIgnoresOtherValues(t *testing.T) {
	doc := `{ "title": "hello", "age": 18 }`
	pat := `{ "title": "goodbye", "class":2217}`

	res, err := jsonPatch(doc, pat)
	if err != nil {
		panic(err)
	}

	exp := `{ "title": "goodbye", "class": 2217 }`

	if !compareJSON(exp, res) {
		t.Fatalf("Key was not replaced")
	}
}

func TestJsonPatchNilDoc(t *testing.T) {
	doc := `{ "title": null }`
	pat := `{ "title": {"foo": "bar"} }`

	res, err := jsonPatch(doc, pat)
	if err != nil {
		panic(err)
	}

	exp := `{ "title": {"foo": "bar"} }`

	if !compareJSON(exp, res) {
		t.Fatalf("Key was not replaced")
	}
}

func TestJsonPatchNilArray(t *testing.T) {

	cases := []arrayCases{
		{`{"a": [ {"b":"c"} ] }`, `{"a": [1]}`, `{"a": [1]}`},
		{`{"a": [ {"b":"c"} ] }`, `{"a": [null, 1]}`, `{"a": [null, 1]}`},
		{`["a",null]`, `[null]`, `[null]`},
		{`["a"]`, `[null]`, `[]`},
		{`["a"]`, `["a",null]`, `["a",null]`},
		{`["a", "b"]`, `["a", null]`, `["a"]`},
		{`["bar","qux","baz"]`, `["bar","baz"]`, `["bar","baz"]`},
		{`{"a":["b"]}`, `{"a": ["b", null]}`, `{"a":["b", null]}`},
		{`{"a":[]}`, `{"a": ["b", null, null, "a"]}`, `{"a":["b", null, null, "a"]}`},
	}

	for _, c := range cases {
		act, err := jsonPatch(c.original, c.patch)
		if err != nil {
			panic(err)
		}

		if !compareJSON(c.res, act) {
			t.Errorf("null values not preserved in array")
		}
	}
}

func TestJsonPatchRecursesIntoObjects(t *testing.T) {
	doc := `{ "person": { "title": "hello", "age": 18 } }`
	pat := `{ "person": { "title": "goodbye", "class":2217} }`

	res, err := jsonPatch(doc, pat)
	if err != nil {
		panic(err)
	}

	exp := `{ "person": { "title": "goodbye", "class": 2217 } }`

	if !compareJSON(exp, res) {
		t.Fatalf("Key was not replaced: %s", res)
	}
}

func TestJsonPatchReplacesNonObjectsWholesale(t *testing.T) {
	a1 := `[1]`
	a2 := `[2]`
	o1 := `{ "a": 1 }`
	o2 := `{ "a": 2 }`
	o3 := `{ "a": 1, "b": 1 }`
	//o4 := `{ "a": 2, "b": 1 }`
	o4 := `{ "a": 2}` // update a's value to 2 and remove b

	cases := []nonObjectCases{
		{a1, a2, a2},
		{o1, a2, a2},
		{a1, o1, o1},
		{o3, o2, o4},
	}

	for _, c := range cases {
		act, err := jsonPatch(c.doc, c.pat)
		if err != nil {
			panic(err)
		}

		if !compareJSON(c.res, act) {
			t.Errorf("whole object replacement failed")
		}
	}
}

func TestJsonPatchReturnsErrorOnBadJSON(t *testing.T) {
	_, err := jsonPatch(`[[[[`, `1`)

	if err == nil {
		t.Errorf("Did not return an error for bad json: %s", err)
	}

	_, err = jsonPatch(`1`, `[[[[`)

	if err == nil {
		t.Errorf("Did not return an error for bad json: %s", err)
	}
}

func TestJsonPatchReturnsEmptyArrayOnEmptyArray(t *testing.T) {
	doc := `{ "array": ["one", "two"] }`
	pat := `{ "array": [] }`

	exp := `{ "array": [] }`

	res, err := MergePatch([]byte(doc), []byte(pat))

	if err != nil {
		t.Errorf("Unexpected error: %s, %s", err, string(res))
	}

	if !compareJSON(exp, string(res)) {
		t.Fatalf("Emtpy array did not return not return as empty array")
	}
}

var rfc6092Tests = []struct {
	target   string
	patch    string
	expected string
}{
	// test cases from https://tools.ietf.org/html/rfc7386#appendix-A
	{target: `{"a":"b"}`, patch: `{"a":"c"}`, expected: `{"a":"c"}`},
	{target: `{"a":"b"}`, patch: `{"b":"c"}`, expected: `{"a":"b","b":"c"}`},
	{target: `{"a":"b"}`, patch: `{"a":null}`, expected: `{}`},
	{target: `{"a":"b","b":"c"}`, patch: `{"a":null}`, expected: `{"b":"c"}`},
	{target: `{"a":["b"]}`, patch: `{"a":"c"}`, expected: `{"a":"c"}`},
	{target: `{"a":"c"}`, patch: `{"a":["b"]}`, expected: `{"a":["b"]}`},
	{target: `{"a":{"b": "c"}}`, patch: `{"a": {"b": "d","c": null}}`, expected: `{"a":{"b":"d"}}`},
	{target: `{"a":[{"b":"c"}]}`, patch: `{"a":[1]}`, expected: `{"a":[1]}`},
	{target: `["a","b"]`, patch: `["c","d"]`, expected: `["c","d"]`},
	{target: `{"a":"b"}`, patch: `["c"]`, expected: `["c"]`},
	{target: `{"a":"foo"}`, patch: `null`, expected: `null`},
	//{target: `{"a":"foo"}`, patch: `"bar"`, expected: `"bar"`},
	{target: `{"e":null}`, patch: `{"a":1}`, expected: `{"a":1,"e":null}`},
	{target: `[1,2]`, patch: `{"a":"b","c":null}`, expected: `{"a":"b"}`},
	{target: `{}`, patch: `{"a":{"bb":{"ccc":null}}}`, expected: `{"a":{"bb":{}}}`},
}

func TestJsonPatchRFC6902Cases(t *testing.T) {
	for i, c := range rfcTests {
		out, err := jsonPatch(c.target, c.expected)
		if err != nil {
			panic(err)
		}

		if !compareJSON(out, c.expected) {
			t.Errorf("case[%d], patch '%s' did not apply properly to '%s'. expected:\n'%s'\ngot:\n'%s'", i, c.patch, c.target, c.expected, out)
		}
	}
}
