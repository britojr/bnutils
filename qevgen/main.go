package main

import (
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/britojr/bnutils/bif"
	"github.com/britojr/lkbn/vars"
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

	b := bif.ParseStruct(inpFile)
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

func sampleQuery(b *bif.Struct, w io.Writer, num int) {
	for i := 0; i < num; i++ {
		v := sampleVar(b.Roots())
		state := sampleState(v)
		line := make([]string, len(b.Variables()))
		line[v.ID()] = strconv.Itoa(state)
		writeLine(w, line)
	}
}

func sampleEvid(b *bif.Struct, w io.Writer, num int) {
	for i := 0; i < num; i++ {
		line := make([]string, len(b.Variables()))
		for _, v := range b.Leafs() {
			line[v.ID()] = strconv.Itoa(sampleState(v))
		}
		writeLine(w, line)
	}
}
