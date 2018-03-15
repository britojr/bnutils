package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/britojr/lkbn/vars"
	"github.com/britojr/utl/conv"
	"github.com/britojr/utl/ioutl"
)

var randSource = rand.New(rand.NewSource(time.Now().UnixNano()))

func main() {
	var inpFile, outFile string
	var num int
	var query bool

	flag.StringVar(&inpFile, "i", "", "input file name in bif format")
	flag.StringVar(&outFile, "o", "", "output file name in .q or .ev format")
	flag.IntVar(&num, "n", 1, "number of queries/evidences to generate")
	flag.Parse()

	if len(inpFile) == 0 || len(outFile) == 0 {
		fmt.Printf("\n error: missing input/output file name\n\n")
		flag.PrintDefaults()
		os.Exit(1)
	}
	switch path.Ext(outFile) {
	case ".q":
		query = true
	case ".ev":
		query = false
	default:
		fmt.Printf("\n error: output file must be '.q' or '.ev'\n\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	b := parseBIFstruct(inpFile)
	f := ioutl.CreateFile(outFile)
	defer f.Close()
	if query {
		sampleQuery(b, f, num)
	} else {
		sampleEvid(b, f, num)
	}
}

func sampleVar(vs vars.VarList) *vars.Var {
	return vs[randSource.Intn(len(vs))]
}

func sampleState(v *vars.Var) int {
	return randSource.Intn(v.NState())
}

func writeLine(w io.Writer, line []string) {
	for i := range line {
		if line[i] == "" {
			line[i] = "*"
		}
	}
	fmt.Fprintf(w, "%s\n", strings.Join(line, ","))
}

func sampleQuery(b *bifStruct, w io.Writer, num int) {
	for i := 0; i < num; i++ {
		v := sampleVar(b.Roots())
		state := sampleState(v)
		line := make([]string, len(b.vs))
		line[v.ID()] = strconv.Itoa(state)
		writeLine(w, line)
	}
}

func sampleEvid(b *bifStruct, w io.Writer, num int) {
	for i := 0; i < num; i++ {
		line := make([]string, len(b.vs))
		for _, v := range b.Leafs() {
			line[v.ID()] = strconv.Itoa(sampleState(v))
		}
		writeLine(w, line)
	}
}

type bifStruct struct {
	vs                vars.VarList
	parents, children map[string]vars.VarList
}

func newBifStruct() *bifStruct {
	b := new(bifStruct)
	b.vs = vars.VarList{}
	b.parents = make(map[string]vars.VarList)
	b.children = make(map[string]vars.VarList)
	return b
}

func (b *bifStruct) String() string {
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

func (b *bifStruct) Leafs() (ls vars.VarList) {
	for _, v := range b.vs {
		if len(b.children[v.Name()]) == 0 {
			ls.Add(v)
		}
	}
	return
}

func (b *bifStruct) Roots() (rs vars.VarList) {
	for _, v := range b.vs {
		if len(b.parents[v.Name()]) == 0 {
			rs.Add(v)
		}
	}
	return
}

func parseBIFstruct(fname string) *bifStruct {
	f := ioutl.OpenFile(fname)
	defer f.Close()
	b := newBifStruct()
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
