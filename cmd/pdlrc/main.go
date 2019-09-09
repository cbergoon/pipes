package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/cbergoon/pipes/pkg/dl"
	"github.com/pkg/errors"
)

func main() {

	flag.Usage = func() {
		fmt.Printf("[USAGE]: pdlrc [-out-file] <in-file-name>\n")
		flag.PrintDefaults()
	}

	var optOutputFile string
	flag.StringVar(&optOutputFile, "out-file", "", "name of output file with extension")

	flag.Parse()

	args := flag.Args()

	if len(args) != 1 {
		flag.Usage()
		os.Exit(-1)
	}

	inFileName := args[0]

	source, err := ioutil.ReadFile(inFileName)
	if err != nil {
		log.Fatal(errors.Wrapf(err, "reverse compile failed: could not read input file %s", inFileName))
	}

	definition, err := dl.GenerateDLFromPipelineDefinitionJSON(source)
	if err != nil {
		log.Fatal(errors.Wrap(err, "reverse compile failed: failed to generate DL"))
	}

	outFileName := ""
	if optOutputFile != "" {
		outFileName = optOutputFile
	} else {
		outFileName = strings.TrimSuffix(inFileName, filepath.Ext(inFileName)) + ".pdl"
	}

	err = ioutil.WriteFile(outFileName, []byte(definition), 0644)
	if err != nil {
		log.Fatal(errors.Wrapf(err, "reverse compile failed: could not write output file %s", outFileName))
	}

}
