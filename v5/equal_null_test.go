package jsonpatch

import "testing"

func TestEqualArraysContainingNull(t *testing.T) {
	cases := []struct {
		name, left, right string
		want              bool
	}{
		{"null element", `[null]`, `[null]`, true},
		{"null versus number", `[null]`, `[1]`, false},
		{"number versus null", `[1]`, `[null]`, false},
		{"null after value", `[true,null]`, `[true,null]`, true},
		{"nested array", `[[null]]`, `[[null]]`, true},
		{"object key order", `{"name":"test","data":[true,null]}`, `{"data":[true,null],"name":"test"}`, true},
		{"different later value", `[null,1]`, `[null,2]`, false},
	}
	for i := range cases {
		tc := cases[i]
		t.Run(tc.name, func(t *testing.T) {
			if got := Equal([]byte(tc.left), []byte(tc.right)); got != tc.want {
				t.Fatalf("Equal(%s, %s) = %v, want %v", tc.left, tc.right, got, tc.want)
			}
		})
	}
}
