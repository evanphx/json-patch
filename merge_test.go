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

		if err != eBadJSONPatch {
			t.Errorf("error not returned properly: %s, %s", err, string(out))
		}
	}

}
