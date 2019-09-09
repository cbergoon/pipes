package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/cbergoon/pipes/pkg/dl"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {

	flag.Usage = func() {
		fmt.Printf("[USAGE]: pdlc [-out-file|-minify-output] <in-file-name>\n")
		flag.PrintDefaults()
	}

	var optOutputFile string
	flag.StringVar(&optOutputFile, "out-file", "", "name of output file with extension")

	var optMinifyOutput bool
	flag.BoolVar(&optMinifyOutput, "minify-output", true, "minifies output file json contents")

	flag.Parse()

	args := flag.Args()

	if len(args) != 1 {
		flag.Usage()
		os.Exit(-1)
	}


	inFileName := args[0]

	source, err := ioutil.ReadFile(inFileName)
	if err != nil {
		log.Fatal(errors.Wrapf(err, "compile failed: could not read input file %s", inFileName))
	}

	l := dl.NewLexer(string(source))
	p := dl.NewParser(l)

	pd, err := p.ParseProgram()
	if err != nil {
		log.Fatal(err)
	}

	var definition []byte
	if optMinifyOutput {
		definition, err = json.Marshal(pd)
		if err != nil {
			log.Fatal(errors.Wrap(err, "compile failed: could not marshal definition"))
		}
	} else {
		definition, err = json.MarshalIndent(pd, "", "  ")
		if err != nil {
			log.Fatal(errors.Wrap(err, "compile failed: could not marshal definition"))
		}
	}

	outFileName := ""
	if optOutputFile != "" {
		outFileName = optOutputFile
	}else{
		outFileName = strings.TrimSuffix(inFileName, filepath.Ext(inFileName)) + ".cpdl"
	}

	err = ioutil.WriteFile(outFileName, definition, 0644)
	if err != nil {
		log.Fatal(errors.Wrapf(err, "compile failed: could not write output file %s", outFileName))
	}

}
