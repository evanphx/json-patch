package jsonpatch

import (
  "bytes"
  "encoding/json"
  "testing"
)

func reformatJSON(j string) string {
  buf := new(bytes.Buffer)

  json.Indent(buf, []byte(j), "", "  ")

  return buf.String()
}

func compareJSON(a, b string) bool {
  return Equal([]byte(a), []byte(b))
}

func applyPatch(doc, patch string) (string, error) {
  obj, err := DecodePatch([]byte(patch))

  if err != nil {
    panic(err)
  }

  out, err := obj.Apply([]byte(doc))

  if err != nil {
    return "", err
  }

  return string(out), nil
}

type Case struct {
  doc, patch, result string
}

var Cases = []Case {
  {
    `{ "foo": "bar"}`,
    `[
         { "op": "add", "path": "/baz", "value": "qux" }
     ]`,
     `{
       "baz": "qux",
       "foo": "bar"
     }`,
  },
  {
   `{ "foo": [ "bar", "baz" ] }`,
   `[
     { "op": "add", "path": "/foo/1", "value": "qux" }
    ]`,
    `{ "foo": [ "bar", "qux", "baz" ] }`,
  },
  {
    `{ "baz": "qux", "foo": "bar" }`,
    `[ { "op": "remove", "path": "/baz" } ]`,
    `{ "foo": "bar" }`,
  },
  {
    `{ "foo": [ "bar", "qux", "baz" ] }`,
    `[ { "op": "remove", "path": "/foo/1" } ]`,
    `{ "foo": [ "bar", "baz" ] }`,
  },
  {
    `{ "baz": "qux", "foo": "bar" }`,
    `[ { "op": "replace", "path": "/baz", "value": "boo" } ]`,
    `{ "baz": "boo", "foo": "bar" }`,
  },
  {
   `{
     "foo": {
       "bar": "baz",
       "waldo": "fred"
     },
     "qux": {
       "corge": "grault"
     }
   }`,
   `[ { "op": "move", "from": "/foo/waldo", "path": "/qux/thud" } ]`,
   `{
     "foo": {
       "bar": "baz"
     },
     "qux": {
       "corge": "grault",
       "thud": "fred"
     }
   }`,
  },
  {
    `{ "foo": [ "all", "grass", "cows", "eat" ] }`,
    `[ { "op": "move", "from": "/foo/1", "path": "/foo/3" } ]`,
    `{ "foo": [ "all", "cows", "eat", "grass" ] }`,
  },
  {
    `{ "foo": "bar" }`,
    `[ { "op": "add", "path": "/child", "value": { "grandchild": { } } } ]`,
    `{ "foo": "bar", "child": { "grandchild": { } } }`,
  },
  {
    `{ "foo": ["bar"] }`,
    `[ { "op": "add", "path": "/foo/-", "value": ["abc", "def"] } ]`,
    `{ "foo": ["bar", ["abc", "def"]] }`,
  },
}

type BadCase struct {
  doc, patch string
}

var BadCases = []BadCase {
  {
    `{ "foo": "bar" }`,
    `[ { "op": "add", "path": "/baz/bat", "value": "qux" } ]`,
  },
}

func TestAllCases(t *testing.T) {
  for _, c := range Cases {
   out, err := applyPatch(c.doc, c.patch)

   if err != nil {
     t.Errorf("Unable to apply patch: %s", err)
   }

   if !compareJSON(out, c.result) {
     t.Errorf("Patch did not apply. Expected:\n%s\n\nActual:\n%s",
      reformatJSON(c.result), reformatJSON(out))
   }
  }

  for _, c := range BadCases {
    _, err := applyPatch(c.doc, c.patch)

    if err == nil {
      t.Errorf("Patch should have failed to apply but it did not")
    }
  }
}

type TestCase struct {
  doc, patch string
  result bool
}

var TestCases = []TestCase {
  {
    `{
       "baz": "qux",
       "foo": [ "a", 2, "c" ]
     }`,
    `[
       { "op": "test", "path": "/baz", "value": "qux" },
       { "op": "test", "path": "/foo/1", "value": 2 }
     ]`,
     true,
  },
  {
    `{ "baz": "qux" }`,
    `[ { "op": "test", "path": "/baz", "value": "bar" } ]`,
    false,
  },
}


func TestAllTest(t *testing.T) {
  for _, c := range TestCases {
    _, err := applyPatch(c.doc, c.patch)

    if c.result && err != nil {
      t.Errorf("Testing failed when it should have passed: %s", err)
    } else if !c.result && err == nil {
      t.Errorf("Testing passed when it should have faild", err)
    }
  }
}
