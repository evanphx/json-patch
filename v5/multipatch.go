package jsonpatch

import (
	"encoding/json"
	"github.com/pkg/errors"
)

type MultiPatch interface {
	Apply(p Patch) error
	ApplyWithOptions(p Patch, options *ApplyOptions) error
	Document() ([]byte, error)
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
		pd:     pd,
		indent: indent,
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
