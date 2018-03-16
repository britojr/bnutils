package bif

import (
	"testing"

	"github.com/britojr/lkbn/vars"
)

func TestBifStruct(t *testing.T) {
	cases := []struct {
		fname                   string
		roots, leafs, internals vars.VarList
	}{{
		"sachs.bif",
		vars.NewList([]int{8, 9}, []int{3}),
		vars.NewList([]int{0, 2, 4, 5}, []int{3}),
		vars.NewList([]int{1, 3, 6, 7, 10}, []int{3}),
	}}
	for _, tt := range cases {
		got := ParseStruct(tt.fname)
		if got == nil {
			t.Errorf("got nil structure for file %v", tt.fname)
		}
		if !tt.leafs.Equal(got.Leafs()) {
			t.Errorf("wrong leafs, want %v, got %v", tt.leafs, got.Leafs())
		}
		if !tt.roots.Equal(got.Roots()) {
			t.Errorf("wrong roots, want %v, got %v", tt.roots, got.Roots())
		}
		if !tt.internals.Equal(got.Internals()) {
			t.Errorf("wrong internals, want %v, got %v", tt.internals, got.Internals())
		}
	}
}
