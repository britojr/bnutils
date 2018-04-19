package bif

import (
	"reflect"
	"testing"

	"github.com/britojr/lkbn/factor"
	"github.com/britojr/lkbn/vars"
)

var sachVars = vars.VarList{
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

func init() {
	for _, v := range sachVars {
		v.SetStates([]string{"LOW", "AVG", "HIGH"})
	}
}

func TestBifVariables(t *testing.T) {
	cases := []struct {
		fname string
		vs    vars.VarList
	}{
		{"sachs.bif", sachVars},
	}
	for _, tt := range cases {
		got, _ := ParseStruct(tt.fname)
		if got == nil {
			t.Errorf("got nil structure for file %v\n", tt.fname)
		}
		if !tt.vs.Equal(got.Variables()) {
			t.Errorf("got different vars\n%v\n!=\n%v\n", tt.vs, got.Variables())
		}
		for i, v := range tt.vs {
			if !reflect.DeepEqual(got.Variables()[i].States(), v.States()) {
				t.Errorf("got different states\n%v\n!=\n%v\n", v.States(), got.Variables()[i].States())
			}
		}
	}
}

func TestBifStruct(t *testing.T) {
	cases := []struct {
		fname                   string
		roots, leafs, internals vars.VarList
	}{{
		"sachs.bif",
		vars.VarList{sachVars[8], sachVars[9]},
		vars.VarList{sachVars[0], sachVars[2], sachVars[4], sachVars[5]},
		vars.VarList{sachVars[1], sachVars[3], sachVars[6], sachVars[7], sachVars[10]},
	}}
	for _, tt := range cases {
		got, _ := ParseStruct(tt.fname)
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

func TestBifFactors(t *testing.T) {
	cases := []struct {
		fname string
		names []string
		fs    []*factor.Factor
	}{{
		"sachs.bif",
		[]string{"PKC", "PKA", "Raf"},
		[]*factor.Factor{
			factor.New(sachVars.FindByName("PKC")).SetValues([]float64{0.42313152, 0.48163920, 0.09522928}),
			factor.New(sachVars.FindByName("PKA"), sachVars.FindByName("PKC")).SetValues(
				[]float64{0.3864255, 0.3794243, 0.2341501,
					0.06039638, 0.92264651, 0.01695712,
					0.01577014, 0.95873839, 0.02549147,
				},
			),
			factor.New(
				sachVars.FindByName("PKA"), sachVars.FindByName("PKC"), sachVars.FindByName("Raf")).SetValues(
				[]float64{
					0.06232176, 0.4475056, 0.84288483, 0.3694012, 0.55082326, 0.74895046, 0.86757991, 8.842572e-01, 0.841807910,
					0.14724878, 0.3125747, 0.12714563, 0.3312117, 0.39291391, 0.15952981, 0.12785388, 1.156677e-01, 0.155367232,
					0.79042946, 0.2399197, 0.02996955, 0.2993871, 0.05626283, 0.09151973, 0.00456621, 7.510891e-05, 0.002824859,
				},
			),
		},
	}}

	for _, tt := range cases {
		got, _ := ParseStruct(tt.fname)
		if got == nil {
			t.Errorf("got nil structure for file %v", tt.fname)
		}
		for i, f := range tt.fs {
			g := got.Factor(tt.names[i])
			if !f.Equal(g) {
				t.Errorf("wrong factor of var %v:\n%v\n%v\n%v\n%v",
					tt.names[i], f.Variables(), g.Variables(), f.Values(), g.Values(),
				)
			}
		}
	}
}
