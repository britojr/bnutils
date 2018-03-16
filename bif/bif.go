package bif

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"

	"github.com/britojr/lkbn/vars"
	"github.com/britojr/utl/conv"
	"github.com/britojr/utl/ioutl"
)

// Struct defines a bif structure
type Struct struct {
	vs                vars.VarList
	parents, children map[string]vars.VarList
}

// NewStruct creates a new empty bif structure
func NewStruct() *Struct {
	b := new(Struct)
	b.vs = vars.VarList{}
	b.parents = make(map[string]vars.VarList)
	b.children = make(map[string]vars.VarList)
	return b
}

func (b *Struct) String() string {
	s := ""
	s += fmt.Sprintf("Vars: (%v)\n", b.vs)
	for _, v := range b.vs {
		s += fmt.Sprintf("PA[%v]:%v\n", v.Name(), b.parents[v.Name()])
	}
	for _, v := range b.vs {
		s += fmt.Sprintf("CH[%v]:%v\n", v.Name(), b.children[v.Name()])
	}
	return s
}

// Variables returns list of network variables
func (b *Struct) Variables() vars.VarList {
	return b.vs
}

// Leafs returns a list of variables that have no children
func (b *Struct) Leafs() (ls vars.VarList) {
	for _, v := range b.vs {
		if len(b.children[v.Name()]) == 0 {
			ls.Add(v)
		}
	}
	return
}

// Roots returns a list of variables that have no parents
func (b *Struct) Roots() (rs vars.VarList) {
	for _, v := range b.vs {
		if len(b.parents[v.Name()]) == 0 {
			rs.Add(v)
		}
	}
	return
}

// Internals returns network internal nodes
func (b *Struct) Internals() (is vars.VarList) {
	for _, v := range b.vs {
		if len(b.parents[v.Name()]) != 0 && len(b.children[v.Name()]) != 0 {
			is.Add(v)
		}
	}
	return
}

// ParseStruct creates a bif struct from a file
func ParseStruct(fname string) *Struct {
	f := ioutl.OpenFile(fname)
	defer f.Close()
	b := NewStruct()
	reVarName := regexp.MustCompile(`variable \s*(\w*)\s*`)
	reCard := regexp.MustCompile(`discrete \[\s*(\d*)\s*\]`)
	reFamily := regexp.MustCompile(`probability\s*\(\s*(\w*)\s*[|]*(.*)\)`)

	scanner := bufio.NewScanner(f)
	var name string
	i := 0
	for scanner.Scan() {
		vName := reVarName.FindStringSubmatch(scanner.Text())
		if len(vName) > 1 {
			name = vName[1]
		}
		vNState := reCard.FindStringSubmatch(scanner.Text())
		if len(vNState) > 1 {
			b.vs.Add(vars.New(i, conv.Atoi(vNState[1]), name, false))
			i++
		}

		vFamily := reFamily.FindStringSubmatch(scanner.Text())
		if len(vFamily) > 2 {
			vFamily[1] = strings.TrimSpace(vFamily[1])
			vFamily[2] = strings.TrimSpace(strings.Replace(vFamily[2], ",", " ", -1))
			vx := b.vs.FindByName(vFamily[1])
			pavx := vars.VarList{}
			for _, vname := range strings.Fields(vFamily[2]) {
				pavx.Add(b.vs.FindByName(vname))
			}
			b.parents[vx.Name()] = pavx
			for _, v := range pavx {
				ch := b.children[v.Name()]
				b.children[v.Name()] = ch.Add(vx)
			}
		}
	}
	return b
}