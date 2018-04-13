package bif

import (
	"bufio"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/britojr/lkbn/factor"
	"github.com/britojr/lkbn/vars"
	"github.com/britojr/utl/conv"
	"github.com/britojr/utl/floats"
	"github.com/britojr/utl/ioutl"
)

// Struct defines a bif structure
type Struct struct {
	vs                vars.VarList
	parents, children map[string]vars.VarList
	factors           map[string]*factor.Factor
}

// NewStruct creates a new empty bif structure
func NewStruct() *Struct {
	b := new(Struct)
	b.vs = vars.VarList{}
	b.parents = make(map[string]vars.VarList)
	b.children = make(map[string]vars.VarList)
	b.factors = make(map[string]*factor.Factor)
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

// Factor returns the corresponding factor for a given variable
func (b *Struct) Factor(vname string) *factor.Factor {
	return b.factors[vname]
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
			stateline := scanner.Text()
			j := strings.Index(stateline, "{")
			if j >= 0 {
				stateline = stateline[j+1:]
				j = strings.Index(stateline, "}")
				for j < 0 {
					scanner.Scan()
					stateline += scanner.Text()
					j = strings.Index(stateline, "}")
				}
				stateline = stateline[:j]
				states := strings.Split(stateline, ",")
				for i, s := range states {
					states[i] = strings.TrimSpace(s)
				}
				b.vs[i].SetStates(states)
			}
			i++
		}

		vFamily := reFamily.FindStringSubmatch(scanner.Text())
		if len(vFamily) > 2 {
			vFamily[1] = strings.TrimSpace(vFamily[1])
			vFamily[2] = strings.TrimSpace(strings.Replace(vFamily[2], ",", " ", -1))
			vx := b.vs.FindByName(vFamily[1])
			pavx := vars.VarList{}
			varOrd := vars.VarList{vx}
			for _, vname := range strings.Fields(vFamily[2]) {
				p := b.vs.FindByName(vname)
				pavx.Add(p)
				varOrd = append(varOrd, p)
			}
			b.parents[vx.Name()] = pavx
			for _, v := range pavx {
				ch := b.children[v.Name()]
				b.children[v.Name()] = ch.Add(vx)
			}
			family := pavx.Union(vars.VarList{vx})
			arranged := make([]float64, family.NStates())

			for scanner.Scan() {
				if strings.TrimSpace(scanner.Text()) == "}" {
					break
				}
				line := strings.Trim(scanner.Text(), ";")
				line = strings.Replace(line, ",", " ", -1)
				if i := strings.Index(line, "table"); i >= 0 {
					values := make([]float64, family.NStates())
					line = line[i+len("table"):]
					for i, v := range strings.Fields(line) {
						values[i] = conv.Atof(strings.TrimSpace(v))
					}
					ixf := vars.NewOrderedIndex(varOrd, family)
					for _, v := range values {
						arranged[ixf.I()] = v
						ixf.Next()
					}
				}
				if i := strings.Index(line, ")"); i >= 0 {
					ixf := vars.NewIndexFor(family, family)
					attrMap := make(map[int]int)
					atts := strings.TrimSpace(line[:i])
					atts = strings.Trim(atts, "()")
					for i, v := range strings.Fields(atts) {
						stName := strings.TrimSpace(v)
						pa := varOrd[i+1]
						attrMap[pa.ID()] = pa.StateID(stName)
					}
					line = line[i+1:]
					acc := 0.0
					for i, v := range strings.Fields(line) {
						attrMap[vx.ID()] = i
						fv := conv.Atof(strings.TrimSpace(v))
						arranged[ixf.AttrbIndex(attrMap)] = fv
						acc += fv
					}
					if !floats.AlmostEqual(acc, 1.0, 1e-6) {
						log.Printf("warnig: unnormalized distribution (%v)\n", acc)
					}
				}
			}
			b.factors[vx.Name()] = factor.New(family...).SetValues(arranged)
		}
	}
	return b
}
