package bif

import (
	"testing"

	"github.com/britojr/lkbn/vars"
)

func TestBifStruct(t *testing.T) {
	vs := []*vars.Var{
		vars.New(0, 3, "Akt", false),
		vars.New(1, 3, "Erk", false),
		vars.New(2, 3, "Jnk", false),
		vars.New(3, 3, "Mek", false),
		vars.New(4, 3, "P38", false),
		vars.New(5, 3, "PIP2", false),
		vars.New(6, 3, "PIP3", false),
		vars.New(7, 3, "PKA", false),
		vars.New(8, 3, "PKC", false),
		vars.New(9, 3, "Plcg", false),
		vars.New(10, 3, "Raf", false),
	}
	cases := []struct {
		fname                   string
		roots, leafs, internals vars.VarList
	}{{
		"sachs.bif",
		vars.VarList{vs[8], vs[9]},
		vars.VarList{vs[0], vs[2], vs[4], vs[5]},
		vars.VarList{vs[1], vs[3], vs[6], vs[7], vs[10]},
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
