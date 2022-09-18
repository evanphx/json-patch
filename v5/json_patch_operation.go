package jsonpatch

import (
	"encoding/json"
)

func raw(value interface{}) (*json.RawMessage, error) {
	raw := new(json.RawMessage)
	bs, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	if err := raw.UnmarshalJSON(bs); err != nil {
		return nil, err
	}
	return raw, nil
}

// { "op": "test", "path": "/a/b/c", "value": "foo" }
func NewTest(path string, value interface{}) (Operation, error) {
	o := make(map[string]*json.RawMessage, 3)
	var err error
	if o["op"], err = raw("test"); err != nil {
		return nil, err
	}
	if o["path"], err = raw(path); err != nil {
		return nil, err
	}
	if o["value"], err = raw(value); err != nil {
		return nil, err
	}
	return o, nil
}

// { "op": "remove", "path": "/a/b/c" }
func NewRemove(path string) (Operation, error) {
	o := make(map[string]*json.RawMessage, 2)
	var err error
	if o["op"], err = raw("remove"); err != nil {
		return nil, err
	}
	if o["path"], err = raw(path); err != nil {
		return nil, err
	}
	return o, nil
}

// { "op": "add", "path": "/a/b/c", "value": [ "foo", "bar" ] }
func NewAdd(path string, value interface{}) (Operation, error) {
	o := make(map[string]*json.RawMessage, 3)
	var err error
	if o["op"], err = raw("add"); err != nil {
		return nil, err
	}
	if o["path"], err = raw(path); err != nil {
		return nil, err
	}
	if o["value"], err = raw(value); err != nil {
		return nil, err
	}
	return o, nil
}

// { "op": "replace", "path": "/a/b/c", "value": 42 }
func NewReplace(path string, value interface{}) (Operation, error) {
	o := make(map[string]*json.RawMessage, 3)
	var err error
	if o["op"], err = raw("replace"); err != nil {
		return nil, err
	}
	if o["path"], err = raw(path); err != nil {
		return nil, err
	}
	if o["value"], err = raw(value); err != nil {
		return nil, err
	}
	return o, nil
}

// { "op": "move", "from": "/a/b/c", "path": "/a/b/d" }
func NewMove(from, path string) (Operation, error) {
	o := make(map[string]*json.RawMessage, 3)
	var err error
	if o["op"], err = raw("move"); err != nil {
		return nil, err
	}
	if o["from"], err = raw(from); err != nil {
		return nil, err
	}
	if o["path"], err = raw(path); err != nil {
		return nil, err
	}
	return o, nil
}

// { "op": "copy", "from": "/a/b/d", "path": "/a/b/e" }
func NewCopy(from, path string) (Operation, error) {
	o := make(map[string]*json.RawMessage, 3)
	var err error
	if o["op"], err = raw("copy"); err != nil {
		return nil, err
	}
	if o["from"], err = raw(from); err != nil {
		return nil, err
	}
	if o["path"], err = raw(path); err != nil {
		return nil, err
	}
	return o, nil
}
