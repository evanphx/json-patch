package json

import (
	"bytes"
	stdlib "encoding/json"
)

var useApi API

type API interface {
	Marshal(v interface{}) ([]byte, error)
	MarshalIndent(v interface{}, prefix, indent string) ([]byte, error)
	Unmarshal(data []byte, v interface{}) error
}

// The standard JSON API does not return a struct,
// so we provide a pass-through implementation
// that implements the standard API's methods.
type stdApi struct {
}

func (api stdApi) Marshal(v interface{}) ([]byte, error) {
	return stdlib.Marshal(v)
}
func (api stdApi) MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return stdlib.MarshalIndent(v, prefix, indent)
}
func (api stdApi) Unmarshal(data []byte, v interface{}) error {
	return stdlib.Unmarshal(data, v)
}

// Methods which we will defer to standard library for
func Compact(dst *bytes.Buffer, src []byte) error {
	return stdlib.Compact(dst, src)
}
func Indent(dst *bytes.Buffer, src []byte, prefix, indent string) error {
	return stdlib.Indent(dst, src, prefix, indent)
}

// Methods which custom interface MUST implement
func Marshal(v interface{}) ([]byte, error) {
	return GetAPI().Marshal(v)
}
func MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return GetAPI().MarshalIndent(v, prefix, indent)
}
func Unmarshal(data []byte, v interface{}) error {
	return GetAPI().Unmarshal(data, v)
}

// GetAPI returns the json API instance being used
func GetAPI() API {
	if useApi == nil {
		useApi = stdApi{}
	}
	return useApi
}

// SetAPI changes the json API instance being used
func SetAPI(newApi API) {
	useApi = newApi
}

// ResetAPI sets back to using standard json API
func ResetAPI() {
	useApi = stdApi{}
}
