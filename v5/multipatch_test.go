package jsonpatch

import (
	"encoding/json"
	"testing"
)

func TestNewMultiPatch(t *testing.T) {
	type args struct {
		doc []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "empty doc",
			args:    args{doc: nil},
			wantErr: true,
		},
		{
			name:    "broken doc",
			args:    args{doc: []byte(`"`)},
			wantErr: true,
		},
		{
			name: "happy path",
			args: args{doc: []byte(`{"a": 1}`)},
			want: []byte(`{"a":1}`),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMultiPatch(tt.args.doc)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMultiPatch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return
			}

			if got == nil {
				t.Errorf("NewMultiPatch() got is nil while should be defined")
				return
			}

			gotDoc, err := got.Document()
			if err != nil {
				t.Errorf("MultiPatch.Document() error: %s", err)
				return
			}

			if string(gotDoc) != string(tt.want) {
				t.Errorf("MultiPatch.Document() got = %v, want %v", string(gotDoc), string(tt.want))
			}
		})
	}
}

func TestNewMultiPatchIndent(t *testing.T) {
	type args struct {
		doc    []byte
		indent string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "empty doc",
			args:    args{doc: nil, indent: "<indent>"},
			wantErr: true,
		},
		{
			name:    "broken doc",
			args:    args{doc: []byte(`"`), indent: "<indent>"},
			wantErr: true,
		},
		{
			name: "happy path",
			args: args{doc: []byte(`{"a": 1}`), indent: "<indent>"},
			want: []byte(`{
<indent>"a": 1
}`),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMultiPatchIndent(tt.args.doc, tt.args.indent)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMultiPatch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return
			}

			if got == nil {
				t.Errorf("NewMultiPatch() got is nil while should be defined")
				return
			}

			gotDoc, err := got.Document()
			if err != nil {
				t.Errorf("MultiPatch.Document() error: %s", err)
				return
			}

			if string(gotDoc) != string(tt.want) {
				t.Errorf("MultiPatch.Document() got = %v, want %v", string(gotDoc), string(tt.want))
			}
		})
	}
}

func Test_multiPatch_Apply(t *testing.T) {
	toJRM := func(msg string) *json.RawMessage {
		m := json.RawMessage(msg)
		return &m
	}

	type args struct {
		doc []byte
		p   []Patch
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "no patches",
			args: args{
				doc: []byte(`{"a": 1}`),
				p:   nil,
			},
			want: []byte(`{"a":1}`),
		},
		{
			name: "happy path",
			args: args{
				doc: []byte(`{"a": 1}`),
				p: []Patch{
					{
						Operation{
							"op":    toJRM(`"add"`),
							"path":  toJRM(`"/b"`),
							"value": toJRM(`2`),
						},
						Operation{
							"op":    toJRM(`"add"`),
							"path":  toJRM(`"/c"`),
							"value": toJRM(`3`),
						},
					},
					{
						Operation{
							"op":    toJRM(`"replace"`),
							"path":  toJRM(`"/b"`),
							"value": toJRM(`20`),
						},
						Operation{
							"op":    toJRM(`"add"`),
							"path":  toJRM(`"/z"`),
							"value": toJRM(`99`),
						},
					},
				},
			},
			want: []byte(`{"a":1,"b":20,"c":3,"z":99}`),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			m, err := NewMultiPatch(tt.args.doc)
			if err != nil {
				t.Errorf("NewMultiPatch() error = %s", err)
				return
			}

			for _, p := range tt.args.p {
				err := m.Apply(p)
				if err != nil {
					t.Errorf("Apply(%v) error = %s", p, err)
					return
				}
			}

			gotDoc, err := m.Document()
			if err != nil {
				t.Errorf("Document() error = %s", err)
				return
			}

			if string(tt.want) != string(gotDoc) {
				t.Errorf("Expected %q, actual %q", string(tt.want), string(gotDoc))
			}
		})
	}
}

func Test_multiPatch_ApplyWithOptions(t *testing.T) {
	toJRM := func(msg string) *json.RawMessage {
		m := json.RawMessage(msg)
		return &m
	}

	type patchWithOptions struct {
		p Patch
		o ApplyOptions
	}
	type args struct {
		doc []byte
		po  []patchWithOptions
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "no patches",
			args: args{
				doc: []byte(`{"a": 1}`),
				po:  nil,
			},
			want: []byte(`{"a":1}`),
		},
		{
			name: "happy path",
			args: args{
				doc: []byte(`{"a": 1}`),
				po: []patchWithOptions{
					{
						p: Patch{
							Operation{
								"op":    toJRM(`"add"`),
								"path":  toJRM(`"/b"`),
								"value": toJRM(`2`),
							},
						},
						o: ApplyOptions{},
					},
					{
						p: Patch{
							Operation{
								"op":    toJRM(`"add"`),
								"path":  toJRM(`"/c/d"`),
								"value": toJRM(`3`),
							},
						},
						o: ApplyOptions{EnsurePathExistsOnAdd: true},
					},
					{
						p: Patch{
							Operation{
								"op":    toJRM(`"replace"`),
								"path":  toJRM(`"/b"`),
								"value": toJRM(`20`),
							},
							Operation{
								"op":   toJRM(`"remove"`),
								"path": toJRM(`"/not_found"`),
							},
							Operation{
								"op":    toJRM(`"add"`),
								"path":  toJRM(`"/z"`),
								"value": toJRM(`99`),
							},
						},
						o: ApplyOptions{AllowMissingPathOnRemove: true},
					},
				},
			},
			want: []byte(`{"a":1,"b":20,"c":{"d":3},"z":99}`),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			m, err := NewMultiPatch(tt.args.doc)
			if err != nil {
				t.Errorf("NewMultiPatch() error = %s", err)
				return
			}

			for _, po := range tt.args.po {
				err := m.ApplyWithOptions(po.p, &po.o)
				if err != nil {
					t.Errorf("Apply(%v) error = %s", po, err)
					return
				}
			}

			gotDoc, err := m.Document()
			if err != nil {
				t.Errorf("Document() error = %s", err)
				return
			}

			if string(tt.want) != string(gotDoc) {
				t.Errorf("Expected %q, actual %q", string(tt.want), string(gotDoc))
			}
		})
	}
}

func Test_multiPatch_Get(t *testing.T) {
	type args struct {
		doc  string
		path string
	}
	tests := []struct {
		name         string
		args         args
		wantExists   bool
		wantIsArray  bool
		wantIsMap    bool
		wantIsString bool
		wantString   string
		wantIsBool   bool
		wantBool     bool
		wantIsNumber bool
		wantNumber   float64
	}{
		{
			name: "empty doc",
			args: args{
				doc:  `{}`,
				path: `"/a"`,
			},
		},
		{
			name: "small doc, node exists",
			args: args{
				doc:  `{"a": 1}`,
				path: `/a`,
			},
			wantExists:   true,
			wantIsNumber: true,
			wantNumber:   1,
		},
		{
			name: "array",
			args: args{
				doc:  `{"a": [1, 2]}`,
				path: `/a`,
			},
			wantExists:  true,
			wantIsArray: true,
		},
		{
			name: "map",
			args: args{
				doc:  `{"a": {"b": [1, 2]}}`,
				path: `/a`,
			},
			wantExists: true,
			wantIsMap:  true,
		},
		{
			name: "element inside of array",
			args: args{
				doc:  `{"a": [1, 2]}`,
				path: `/a/1`,
			},
			wantExists:   true,
			wantIsNumber: true,
			wantNumber:   2,
		},
		{
			name: "element inside of map",
			args: args{
				doc:  `{"a": {"b": [1, 2]}}`,
				path: `/a/b`,
			},
			wantExists:  true,
			wantIsArray: true,
		},
		{
			name: "string value",
			args: args{
				doc:  `{"a": "b"}`,
				path: `/a`,
			},
			wantExists:   true,
			wantIsString: true,
			wantString:   "b",
		},
		{
			name: "empty string value",
			args: args{
				doc:  `{"a": ""}`,
				path: `/a`,
			},
			wantExists:   true,
			wantIsString: true,
			wantString:   "",
		},
		{
			name: "bool value",
			args: args{
				doc:  `{"a": true}`,
				path: `/a`,
			},
			wantExists: true,
			wantIsBool: true,
			wantBool:   true,
		},
		{
			name: "bool false value",
			args: args{
				doc:  `{"a": false}`,
				path: `/a`,
			},
			wantExists: true,
			wantIsBool: true,
			wantBool:   false,
		},
		{
			name: "number value",
			args: args{
				doc:  `{"a": 99.9}`,
				path: `/a`,
			},
			wantExists:   true,
			wantIsNumber: true,
			wantNumber:   99.9,
		},
		{
			name: "0 number value",
			args: args{
				doc:  `{"a": 0}`,
				path: `/a`,
			},
			wantExists:   true,
			wantIsNumber: true,
			wantNumber:   0,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			m, err := NewMultiPatch([]byte(tt.args.doc))
			if err != nil {
				t.Errorf("NewMultiPatch() error = %s", err)
				return
			}

			got := m.Get(tt.args.path)
			if tt.wantExists != got.Exists() {
				t.Errorf("Exists() expected = %t, actual = %t", tt.wantExists, got.Exists())
			}
			if !tt.wantExists {
				return
			}

			if tt.wantIsArray != got.IsArray() {
				t.Errorf("IsArray() expected = %t, actual = %t", tt.wantIsArray, got.IsArray())
			}

			if tt.wantIsMap != got.IsMap() {
				t.Errorf("IsMap() expected = %t, actual = %t", tt.wantIsMap, got.IsMap())
			}

			if tt.wantIsString != got.IsString() {
				t.Errorf("IsString() expected = %t, actual = %t", tt.wantIsString, got.IsString())
			}

			if tt.wantString != got.String() {
				t.Errorf("String() expected = %s, actual = %s", tt.wantString, got.String())
			}

			if tt.wantIsNumber != got.IsNumber() {
				t.Errorf("IsNumber() expected = %t, actual = %t", tt.wantIsNumber, got.IsNumber())
			}

			if tt.wantNumber != got.Number() {
				t.Errorf("Number() expected = %f, actual = %f", tt.wantNumber, got.Number())
			}

			if tt.wantIsBool != got.IsBool() {
				t.Errorf("IsBool() expected = %t, actual = %t", tt.wantIsBool, got.IsBool())
			}

			if tt.wantBool != got.Bool() {
				t.Errorf("Bool() expected = %t, actual = %t", tt.wantBool, got.Bool())
			}
		})
	}
}
