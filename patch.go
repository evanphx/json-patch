package jsonpatch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

const (
	eRaw = iota
	eDoc = iota
	eAry = iota
)

type lazyNode struct {
	raw   *json.RawMessage
	doc   partialDoc
	ary   partialArray
	which int
}

func newLazyNode(raw *json.RawMessage) *lazyNode {
	return &lazyNode{raw: raw, doc: nil, ary: nil, which: eRaw}
}

func (n *lazyNode) MarshalJSON() ([]byte, error) {
	switch n.which {
	case eRaw:
		return *n.raw, nil
	case eDoc:
		return json.Marshal(n.doc)
	case eAry:
		return json.Marshal(n.ary)
	default:
		return nil, fmt.Errorf("Unknown type")
	}
}

func (n *lazyNode) UnmarshalJSON(data []byte) error {
	dest := make(json.RawMessage, len(data))
	copy(dest, data)
	n.raw = &dest
	n.which = eRaw
	return nil
}

func (n *lazyNode) IntoDoc() (*partialDoc, error) {
	if n.which == eDoc {
		return &n.doc, nil
	}

	err := json.Unmarshal(*n.raw, &n.doc)

	if err != nil {
		return nil, err
	}

	n.which = eDoc
	return &n.doc, nil
}

func (n *lazyNode) IntoAry() (*partialArray, error) {
	if n.which == eAry {
		return &n.ary, nil
	}

	err := json.Unmarshal(*n.raw, &n.ary)

	if err != nil {
		return nil, err
	}

	n.which = eAry
	return &n.ary, nil
}

func (n *lazyNode) compact() []byte {
	buf := new(bytes.Buffer)

	err := json.Compact(buf, *n.raw)

	if err != nil {
		return *n.raw
	}

	return buf.Bytes()
}

func (n *lazyNode) tryDoc() bool {
	err := json.Unmarshal(*n.raw, &n.doc)

	if err != nil {
		return false
	}

	n.which = eDoc
	return true
}

func (n *lazyNode) tryAry() bool {
	err := json.Unmarshal(*n.raw, &n.ary)

	if err != nil {
		return false
	}

	n.which = eAry
	return true
}

func (n *lazyNode) equal(o *lazyNode) bool {
	if n.which == eRaw {
		if !n.tryDoc() && !n.tryAry() {
			if o.which != eRaw {
				return false
			}

			return bytes.Equal(n.compact(), o.compact())
		}
	}

	if n.which == eDoc {
		if o.which == eRaw {
			if !o.tryDoc() {
				return false
			}
		}

		if o.which != eDoc {
			return false
		}

		for k, v := range n.doc {
			ov, ok := o.doc[k]

			if !ok {
				return false
			}

			if !v.equal(ov) {
				return false
			}
		}

		return true
	}

	if o.which != eAry && !o.tryAry() {
		return false
	}

	if len(n.ary) != len(o.ary) {
		return false
	}

	for idx, val := range n.ary {
		if !val.equal(o.ary[idx]) {
			return false
		}
	}

	return true
}

// Indicate if 2 JSON documents have the same structural equality
func Equal(a, b []byte) bool {
	ra := make(json.RawMessage, len(a))
	copy(ra, a)
	la := newLazyNode(&ra)

	rb := make(json.RawMessage, len(b))
	copy(rb, b)
	lb := newLazyNode(&rb)

	return la.equal(lb)
}

type Operation map[string]*json.RawMessage
type Patch []Operation

type partialDoc map[string]*lazyNode
type partialArray []*lazyNode

type Container interface {
	Get(key string) (*lazyNode, error)
	Set(key string, val *lazyNode) error
	Remove(key string) error
}

func DecodePatch(buf []byte) (Patch, error) {
	var p Patch

	err := json.Unmarshal(buf, &p)

	if err != nil {
		return nil, err
	}

	return p, nil
}

func (p Patch) Operation(i int) Operation {
	if i < 0 || i >= len(p) {
		return nil
	}

	return p[i]
}

func (o Operation) Kind() string {
	if obj, ok := o["op"]; ok {
		var op string

		err := json.Unmarshal(*obj, &op)

		if err != nil {
			return "unknown"
		}

		return op
	}

	return "unknown"
}

func (o Operation) Path() string {
	if obj, ok := o["path"]; ok {
		var op string

		err := json.Unmarshal(*obj, &op)

		if err != nil {
			return "unknown"
		}

		return op
	}

	return "unknown"
}

func (o Operation) From() string {
	if obj, ok := o["from"]; ok {
		var op string

		err := json.Unmarshal(*obj, &op)

		if err != nil {
			return "unknown"
		}

		return op
	}

	return "unknown"
}

func (o Operation) Value() *lazyNode {
	if obj, ok := o["value"]; ok {
		return newLazyNode(obj)
	}

	return nil
}

func isArray(buf []byte) bool {
	for _, c := range buf {
		switch c {
		case ' ':
		case '\n':
		case '\t':
			continue
		case '[':
			return true
		default:
			break
		}
	}

	return false
}

func findObject(doc *partialDoc, path string) (Container, string) {
	split := strings.Split(path, "/")

	parts := split[1 : len(split)-1]

	key := split[len(split)-1]

	var err error

	for idx, part := range parts {
		next, ok := (*doc)[part]
		if !ok {
			return nil, ""
		}

		if isArray(*next.raw) {
			if idx == len(parts)-1 {
				ary, err := next.IntoAry()

				if err != nil {
					return nil, ""
				}

				return ary, key
			} else {
				return nil, ""
			}
		} else {
			doc, err = next.IntoDoc()

			if err != nil {
				return nil, ""
			}
		}
	}

	return doc, key
}

func (d *partialDoc) Set(key string, val *lazyNode) error {
	(*d)[key] = val
	return nil
}

func (d *partialDoc) Get(key string) (*lazyNode, error) {
	return (*d)[key], nil
}

func (d *partialDoc) Remove(key string) error {
	delete(*d, key)
	return nil
}

func (d *partialArray) Set(key string, val *lazyNode) error {
	if key == "-" {
		*d = append(*d, val)
		return nil
	}

	idx, err := strconv.Atoi(key)

	if err != nil {
		return err
	}

	ary := make([]*lazyNode, len(*d)+1)

	cur := *d

	copy(ary[0:idx], cur[0:idx])
	ary[idx] = val
	copy(ary[idx+1:], cur[idx:])

	*d = ary
	return nil
}

func (d *partialArray) Get(key string) (*lazyNode, error) {
	idx, err := strconv.Atoi(key)

	if err != nil {
		return nil, err
	}

	return (*d)[idx], nil
}

func (d *partialArray) Remove(key string) error {
	idx, err := strconv.Atoi(key)

	if err != nil {
		return err
	}

	cur := *d

	ary := make([]*lazyNode, len(cur)-1)

	copy(ary[0:idx], cur[0:idx])
	copy(ary[idx:], cur[idx+1:])

	*d = ary
	return nil

}

func (p Patch) add(doc *partialDoc, op Operation) error {
	path := op.Path()

	con, key := findObject(doc, path)

	if con == nil {
		return fmt.Errorf("Missing container: %s", path)
	}

	con.Set(key, op.Value())

	return nil
}

func (p Patch) remove(doc *partialDoc, op Operation) error {
	path := op.Path()

	con, key := findObject(doc, path)

	return con.Remove(key)
}

func (p Patch) replace(doc *partialDoc, op Operation) error {
	path := op.Path()

	con, key := findObject(doc, path)

	con.Set(key, op.Value())

	return nil
}

func (p Patch) move(doc *partialDoc, op Operation) error {
	from := op.From()

	con, key := findObject(doc, from)

	val, err := con.Get(key)

	if err != nil {
		return err
	}

	con.Remove(key)

	path := op.Path()

	con, key = findObject(doc, path)

	con.Set(key, val)

	return nil
}

var eTestFailed = fmt.Errorf("Testing value failed")

func (p Patch) test(doc *partialDoc, op Operation) error {
	path := op.Path()

	con, key := findObject(doc, path)

	val, err := con.Get(key)

	if err != nil {
		return err
	}

	if val.equal(op.Value()) {
		return nil
	}

	return eTestFailed
}

func (p Patch) Apply(doc []byte) ([]byte, error) {
	pd := new(partialDoc)

	err := json.Unmarshal(doc, pd)

	if err != nil {
		return nil, err
	}

	err = nil

	for _, op := range p {
		switch op.Kind() {
		case "add":
			err = p.add(pd, op)
		case "remove":
			err = p.remove(pd, op)
		case "replace":
			err = p.replace(pd, op)
		case "move":
			err = p.move(pd, op)
		case "test":
			err = p.test(pd, op)
		default:
			err = fmt.Errorf("Unexpected kind: %s", op.Kind())
		}

		if err != nil {
			return nil, err
		}
	}

	return json.Marshal(pd)
}
