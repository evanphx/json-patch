package jsonpatch

import (
	"encoding/json"
	"strings"
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
		p       Patch
		o       ApplyOptions
		wantErr string
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
						o:       ApplyOptions{},
						wantErr: "",
					},
					{
						p: Patch{
							Operation{
								"op":    toJRM(`"add"`),
								"path":  toJRM(`"/c/d"`),
								"value": toJRM(`3`),
							},
						},
						o:       ApplyOptions{EnsurePathExistsOnAdd: true},
						wantErr: "unexpected kind: unknown",
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
					if po.wantErr == "" {
						t.Errorf("Apply(%v) error = %s", po, err)
						return
					} else if !strings.Contains(err.Error(), po.wantErr) {
						t.Errorf("Apply(%v) expected error = %s, actual error = %s", po, po.wantErr, err)
						return
					}
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
