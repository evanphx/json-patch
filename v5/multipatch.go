package jsonpatch

import (
	"encoding/json"
	"github.com/pkg/errors"
)

type MultiPatch interface {
	Apply(p Patch) error
	ApplyWithOptions(p Patch, options *ApplyOptions) error
	Document() ([]byte, error)
	Get(path string) Result
}

type Result interface {
	IsArray() bool
	Array() []Result
	IsMap() bool
	Map() map[string]Result
	String() string
	Number() float64
	Bool() bool
	Exists() bool
}

type multiPatch struct {
	pd                  container
	indent              string
	defaultApplyOptions *ApplyOptions
}

func NewMultiPatch(doc []byte) (MultiPatch, error) {
	return NewMultiPatchIndent(doc, "")
}

func NewMultiPatchIndent(doc []byte, indent string) (MultiPatch, error) {
	if len(doc) == 0 {
		return nil, errors.New("empty doc")
	}

	pd, err := createContainer(doc)
	if err != nil {
		return nil, err
	}

	return &multiPatch{
		pd:                  pd,
		indent:              indent,
		defaultApplyOptions: NewApplyOptions(),
	}, nil
}

// Apply mutates a JSON document according to the patch.
// It returns an error if the patch failed.
func (m *multiPatch) Apply(p Patch) error {
	return m.ApplyWithOptions(p, m.defaultApplyOptions)
}

// ApplyWithOptions mutates a JSON document according to the patch and the passed in ApplyOptions.
// It returns an error if the patch failed.
func (m *multiPatch) ApplyWithOptions(p Patch, options *ApplyOptions) error {
	return p.applyIndentWithOptions(m.pd, options)
}

// Document converts the internal representation of the JSON into a serialized document
func (m *multiPatch) Document() ([]byte, error) {
	if m.indent != "" {
		return json.MarshalIndent(m.pd, "", m.indent)
	}

	return json.Marshal(m.pd)
}

func (m *multiPatch) Get(path string) Result {
	con, key := findObject(&m.pd, path, m.defaultApplyOptions)
	node, err := con.get(key, m.defaultApplyOptions)

	if err != nil {
		return &result{
			node:   nil,
			exists: false,
		}
	}

	return &result{
		node:   node,
		exists: true,
	}
}

type result struct {
	node   *lazyNode
	exists bool
}

func (r result) IsArray() bool {
	return r.exists && r.node.tryAry()
}

func (r result) Array() []Result {
	if !r.IsArray() {
		return nil
	}
	partialArray, err := r.node.intoAry()
	if err != nil || partialArray == nil {
		return nil
	}

	results := make([]Result, 0, len(*partialArray))
	for _, node := range *partialArray {
		results = append(results, result{node: node, exists: true})
	}

	return results
}

func (r result) IsMap() bool {
	return r.exists && r.node.tryDoc()
}

func (r result) Map() map[string]Result {
	resultMap := make(map[string]Result)
	if !r.IsMap() {
		return resultMap
	}

	partialDoc, err := r.node.intoDoc()
	if err != nil {
		return resultMap
	}

	for _, key := range partialDoc.keys {
		if node, ok := partialDoc.obj[key]; ok {
			resultMap[key] = &result{node: node, exists: true}
		}
	}

	return resultMap
}

func (r result) String() string {
	if !r.Exists() || r.IsMap() || r.IsArray() || r.node.raw == nil {
		return ""
	}

	var value interface{}
	err := json.Unmarshal(*r.node.raw, &value)
	if err != nil {
		return ""
	}

	strValue, ok := value.(string)
	if !ok {
		return ""
	}

	return strValue
}

func (r result) Number() float64 {
	if !r.Exists() || r.IsMap() || r.IsArray() || r.node.raw == nil {
		return 0.0
	}

	var value interface{}
	err := json.Unmarshal(*r.node.raw, &value)
	if err != nil {
		return 0.0
	}

	numValue, ok := value.(float64)
	if !ok {
		return 0.0
	}

	return numValue
}

func (r result) Bool() bool {
	if !r.Exists() || r.IsMap() || r.IsArray() || r.node.raw == nil {
		return false
	}

	var value interface{}
	err := json.Unmarshal(*r.node.raw, &value)
	if err != nil {
		return false
	}

	boolValue, ok := value.(bool)
	if !ok {
		return false
	}

	return boolValue
}

func (r result) Exists() bool {
	return r.exists
}
