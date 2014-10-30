package jsonpatch

import (
	"strings"
	"testing"
)

func mergePatch(doc, patch string) string {
	out, err := MergePatch([]byte(doc), []byte(patch))

	if err != nil {
		panic(err)
	}

	return string(out)
}

func TestMergePatchReplaceKey(t *testing.T) {
	doc := `{ "title": "hello" }`
	pat := `{ "title": "goodbye" }`

	res := mergePatch(doc, pat)

	if !compareJSON(pat, res) {
		t.Fatalf("Key was not replaced")
	}
}

func TestMergePatchIgnoresOtherValues(t *testing.T) {
	doc := `{ "title": "hello", "age": 18 }`
	pat := `{ "title": "goodbye" }`

	res := mergePatch(doc, pat)

	exp := `{ "title": "goodbye", "age": 18 }`

	if !compareJSON(exp, res) {
		t.Fatalf("Key was not replaced")
	}
}

func TestMergePatchRecursesIntoObjects(t *testing.T) {
	doc := `{ "person": { "title": "hello", "age": 18 } }`
	pat := `{ "person": { "title": "goodbye" } }`

	res := mergePatch(doc, pat)

	exp := `{ "person": { "title": "goodbye", "age": 18 } }`

	if !compareJSON(exp, res) {
		t.Fatalf("Key was not replaced")
	}
}

type nonObjectCases struct {
	doc, pat, res string
}

func TestMergePatchReplacesNonObjectsWholesale(t *testing.T) {
	a1 := `[1]`
	a2 := `[2]`
	o1 := `{ "a": 1 }`
	o2 := `{ "a": 2 }`
	o3 := `{ "a": 1, "b": 1 }`
	o4 := `{ "a": 2, "b": 1 }`

	cases := []nonObjectCases{
		{a1, a2, a2},
		{o1, a2, a2},
		{a1, o1, o1},
		{o3, o2, o4},
	}

	for _, c := range cases {
		act := mergePatch(c.doc, c.pat)

		if !compareJSON(c.res, act) {
			t.Errorf("whole object replacement failed")
		}
	}
}

func TestMergePatchReturnsErrorOnBadJSON(t *testing.T) {
	_, err := MergePatch([]byte(`[[[[`), []byte(`1`))

	if err == nil {
		t.Errorf("Did not return an error for bad json: %s", err)
	}

	_, err = MergePatch([]byte(`1`), []byte(`[[[[`))

	if err == nil {
		t.Errorf("Did not return an error for bad json: %s", err)
	}
}

var rfcTests = `
     {"a":"b"}   |   {"a":"c"}    |  {"a":"c"}
     {"a":"b"}   |   {"b":"c"}    |  {"a":"b", "b":"c"}
     {"a":"b"}   |   {"a":null}   |  {}
     {"a":["b"]} |   {"a":"c"}    |  {"a":"c"}
     {"a":"c"}   |   {"a":["b"]}  |  {"a":["b"]}
     ["a","b"]   |   ["c","d"]    |  ["c","d"]
     {"a":"b"}   |   ["c"]        |  ["c"]
     {"e":null}  |   {"a":1}      |  {"e":null, "a":1}

     {"a":"b", "b":"c"}  |   {"a":null}   |  {"b":"c"}
     {"a": [{"b":"c"}]}  |   {"a": [1]}   |  {"a": [1]}

     [1,2]       |   {"a":"b","c":null}   |  {"a":"b"}

     {"a": { "b": "c" } } | { "a": { "b": "d", "c": null } } | { "a": { "b": "d" } }

     {}          | {"a": { "bb": { "ccc": null }}} | {"a": { "bb": {}}}

     {"a":"foo"} | {"b": [3, null, {"x": null}]} | {"a":"foo", "b": [3, {}]}


     [1,2]       | [1,null,3]     | [1,3]

     [1,2]       | [1,null,2]     | [1,2]

     {"a":"b"}   | {"a": [ {"z":1, "b":null}]} | {"a": [ {"z":1}]}
`

func TestMergePatchRFCCases(t *testing.T) {
	tests := strings.Split(rfcTests, "\n")

	for _, c := range tests {
		if strings.TrimSpace(c) == "" {
			continue
		}

		parts := strings.SplitN(c, "|", 3)

		doc := strings.TrimSpace(parts[0])
		pat := strings.TrimSpace(parts[1])
		res := strings.TrimSpace(parts[2])

		out := mergePatch(doc, pat)

		if !compareJSON(out, res) {
			t.Errorf("patch '%s' did not apply properly to '%s': '%s'", pat, doc, out)
		}
	}
}

var rfcFailTests = `
     {"a":"foo"}  |   null
     {"a":"foo"}  |   "bar"
`

func TestMergePatchFailRFCCases(t *testing.T) {
	tests := strings.Split(rfcFailTests, "\n")

	for _, c := range tests {
		if strings.TrimSpace(c) == "" {
			continue
		}

		parts := strings.SplitN(c, "|", 2)

		doc := strings.TrimSpace(parts[0])
		pat := strings.TrimSpace(parts[1])

		out, err := MergePatch([]byte(doc), []byte(pat))

		if err != errBadJSONPatch {
			t.Errorf("error not returned properly: %s, %s", err, string(out))
		}
	}

}

func TestMergeReplaceKey(t *testing.T) {
	doc := `{ "title": "hello", "nested": {"one": 1, "two": 2} }`
	pat := `{ "title": "goodbye", "nested": {"one": 2, "two": 2}  }`

	exp := `{ "title": "goodbye", "nested": {"one": 2}  }`

	res, err := CreateMergePatch([]byte(doc), []byte(pat))

	if err != nil {
		t.Errorf("Unexpected error: %s, %s", err, string(res))
	}

	if !compareJSON(exp, string(res)) {
		t.Fatalf("Key was not replaced")
	}
}

func TestMergeGetArray(t *testing.T) {
	doc := `{ "title": "hello", "array": ["one", "two"], "notmatch": [1, 2, 3] }`
	pat := `{ "title": "hello", "array": ["one", "two", "three"], "notmatch": [1, 2, 3]  }`

	exp := `{ "array": ["one", "two", "three"] }`

	res, err := CreateMergePatch([]byte(doc), []byte(pat))

	if err != nil {
		t.Errorf("Unexpected error: %s, %s", err, string(res))
	}

	if !compareJSON(exp, string(res)) {
		t.Fatalf("Array was not added")
	}
}

func TestMergeGetObjArray(t *testing.T) {
	doc := `{ "title": "hello", "array": [{"banana": true}, {"evil": false}], "notmatch": [{"one":1}, {"two":2}, {"three":3}] }`
	pat := `{ "title": "hello", "array": [{"banana": false}, {"evil": true}], "notmatch": [{"one":1}, {"two":2}, {"three":3}] }`

	exp := `{  "array": [{"banana": false}, {"evil": true}] }`

	res, err := CreateMergePatch([]byte(doc), []byte(pat))

	if err != nil {
		t.Errorf("Unexpected error: %s, %s", err, string(res))
	}

	if !compareJSON(exp, string(res)) {
		t.Fatalf("Object array was not added")
	}
}

func TestMergeDeleteKey(t *testing.T) {
	doc := `{ "title": "hello", "nested": {"one": 1, "two": 2} }`
	pat := `{ "title": "hello", "nested": {"one": 1}  }`

	exp := `{"nested":{"two":null}}`

	res, err := CreateMergePatch([]byte(doc), []byte(pat))

	if err != nil {
		t.Errorf("Unexpected error: %s, %s", err, string(res))
	}

	// We cannot use "compareJSON", since Equals does not report a difference if the value is null
	if exp != string(res) {
		t.Fatalf("Key was not removed")
	}
}

func TestMergeEmptyArray(t *testing.T) {
	doc := `{ "array": null }`
	pat := `{ "array": [] }`

	exp := `{"array":[]}`

	res, err := CreateMergePatch([]byte(doc), []byte(pat))

	if err != nil {
		t.Errorf("Unexpected error: %s, %s", err, string(res))
	}

	// We cannot use "compareJSON", since Equals does not report a difference if the value is null
	if exp != string(res) {
		t.Fatalf("Key was not removed")
	}
}

func TestMergeObjArray(t *testing.T) {
	doc := `{ "array": [ {"a": {"b": 2}}, {"a": {"b": 3}} ]}`
	exp := `{}`

	res, err := CreateMergePatch([]byte(doc), []byte(doc))

	if err != nil {
		t.Errorf("Unexpected error: %s, %s", err, string(res))
	}

	// We cannot use "compareJSON", since Equals does not report a difference if the value is null
	if exp != string(res) {
		t.Fatalf("Array was not empty, was " + string(res))
	}
}

func TestMergeComplexMatch(t *testing.T) {
	doc := `{"hello": "world","t": true ,"f": false, "n": null,"i": 123,"pi": 3.1416,"a": [1, 2, 3, 4], "nested": {"hello": "world","t": true ,"f": false, "n": null,"i": 123,"pi": 3.1416,"a": [1, 2, 3, 4]} }`
	empty := `{}`
	res, err := CreateMergePatch([]byte(doc), []byte(doc))

	if err != nil {
		t.Errorf("Unexpected error: %s, %s", err, string(res))
	}

	// We cannot use "compareJSON", since Equals does not report a difference if the value is null
	if empty != string(res) {
		t.Fatalf("Did not get empty result, was:%s", string(res))
	}
}

func TestMergeComplexAddAll(t *testing.T) {
	doc := `{"hello": "world","t": true ,"f": false, "n": null,"i": 123,"pi": 3.1416,"a": [1, 2, 3, 4], "nested": {"hello": "world","t": true ,"f": false, "n": null,"i": 123,"pi": 3.1416,"a": [1, 2, 3, 4]} }`
	empty := `{}`
	res, err := CreateMergePatch([]byte(empty), []byte(doc))

	if err != nil {
		t.Errorf("Unexpected error: %s, %s", err, string(res))
	}

	if !compareJSON(doc, string(res)) {
		t.Fatalf("Did not get everything as, it was:\n%s", string(res))
	}
}

func TestMergeComplexRemoveAll(t *testing.T) {
	doc := `{"hello": "world","t": true ,"f": false, "n": null,"i": 123,"pi": 3.1416,"a": [1, 2, 3, 4], "nested": {"hello": "world","t": true ,"f": false, "n": null,"i": 123,"pi": 3.1416,"a": [1, 2, 3, 4]} }`
	exp := `{"a":null,"f":null,"hello":null,"i":null,"n":null,"nested":null,"pi":null,"t":null}`
	empty := `{}`
	res, err := CreateMergePatch([]byte(doc), []byte(empty))

	if err != nil {
		t.Errorf("Unexpected error: %s, %s", err, string(res))
	}

	if exp != string(res) {
		t.Fatalf("Did not get result, was:%s", string(res))
	}

	// FIXME: Crashes if using compareJSON like this:
	/*
		if !compareJSON(doc, string(res)) {
			t.Fatalf("Did not get everything as, it was:\n%s", string(res))
		}
	*/
}
