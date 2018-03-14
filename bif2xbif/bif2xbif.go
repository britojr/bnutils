package main

import (
	"flag"
	"fmt"
)

func main() {
	var inpFile, outFile string
	flag.StringVar(&inpFile, "i", "", "input file name in bif format")
	flag.StringVar(&outFile, "o", "", "output file name in xbif format")
	flag.Parse()

	if len(inpFile) == 0 || len(outFile) == 0 {
		flag.PrintDefaults()
	}

	fmt.Println(inpFile, outFile)
}
