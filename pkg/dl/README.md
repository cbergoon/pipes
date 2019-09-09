<h1 align="center">Pipes DL</h1>
<p align="center">
<a href="https://goreportcard.com/report/github.com/cbergoon/pipes-dl"><img src="https://goreportcard.com/badge/github.com/cbergoon/pipes-dl?1=1" alt="Report"></a>
<a href="https://godoc.org/github.com/cbergoon/pipes-dl"><img src="https://img.shields.io/badge/godoc-reference-brightgreen.svg" alt="Docs"></a>
<a href="#"><img src="https://img.shields.io/badge/version-0.1.0-brightgreen.svg" alt="Version"></a>
</p>

Pipes DL defines a definition language for pipelines built using [pipes](https://github.com/cbergoon/pipes). It provides a human
readable and feature complete way to define pipelines and individual processes.

For now Pipes DL is a proof of concept and should not be used in production yet.

#### Features

* Feature complete pipeline definition language for [pipes](https://github.com/cbergoon/pipes)
* Ability to define processes.
* Ability to define connections.
* Compiler tool to compile to JSON.
* DL generator tool to decompile from JSON to PDL.
* FUTURE: Syntax definition for editors.
* FUTURE: Improved JS and JSON definitions.

#### Installation

Get the source with ```go get```:

```bash
$ go get github.com/cbergoon/pipes-dl
```

Then import the library in your project:

```go
import "github.com/cbergoon/pipes-dl"
```

#### Usage

To compile a PDL definition to a JSON definition readable by [pipes](https://github.com/cbergoon/pipes). Use the `pdlc` tool.

```[USAGE]: pdlrc [-out-file] <in-file-name>```

```bash
$ pdlc example.pdl
```

To decompile a CPDL (compiled PDL) json definition. Use the `pdlrc` tool.

```[USAGE]: pdlc [-out-file|-minify-output] <in-file-name>```

```bash
$ pdlrc example.cpdl
```

#### Documentation

The same 3 concepts in pipes are present in pipes-dl: pipelines, processes, and connections.

##### CREATE PIPELINE Statement

The `CREATE PIPELINE` statement begins the definition and names the pipeline.

##### ADD (PROCESS) Statement

The `ADD` statement is used to add a process to the pipeline. The process type and name are specified followed by
`INPUT` and `OUTPUT` definitions. Any required initial state can be included as key-value pairs after the `SET` keyword.

The final process in the pipeline must be declared as `SINK`.

##### CONNECT Statement

The `CONNECT` statement links the outputs of one process to the inputs of another. The format for specifying a port is
`<PROCESS>:<PORT_NAME>.

#### Example Definition

```pdl
CREATE PIPELINE "MyPipeline";

ADD "Alfa" OF "Generator" OUTPUTS = ("Out1", "Out2");
ADD "Beta" OF "DynamicJs"
    INPUTS = ("In1", "In2")
    OUTPUTS = ("Out")
    SET "src" = 'o = {
        "MyVal": In1 + "hello" + In2
    };
    console.log("hellofrom js");
    Out = JSON.stringify(o);',
    "gg" = "kk";
ADD SINK "Charlie" OF "Printer" INPUTS = ("In");

CONNECT "Alfa":"Out1" TO "Beta":"In1";
CONNECT "Alfa":"Out2" TO "Beta":"In2";
CONNECT "Beta":"Out" TO "Charlie":"In";
```

#### Example Parser Usage

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/cbergoon/pipes-dl"
	"github.com/pkg/errors"
)

func main() {
	source := `CREATE PIPELINE "MyPipeline";

	ADD "Alfa" OF "Generator" OUTPUTS = ("Out1", "Out2");
	ADD "Beta" OF "DynamicJs"
	INPUTS = ("In1", "In2")
	OUTPUTS = ("Out")
	SET "src" = 'o = {
		"MyVal": In1 + "hello" + In2
	};
	console.log("hellofrom js");
	Out = JSON.stringify(o);',
	"gg" = "kk";
	ADD SINK "Charlie" OF "Printer" INPUTS = ("In");

	CONNECT "Alfa":"Out1" TO "Beta":"In1";
	CONNECT "Alfa":"Out2" TO "Beta":"In2";
	CONNECT "Beta":"Out" TO "Charlie":"In";`

	l := pipesdl.NewLexer(source)
	p := pipesdl.NewParser(l)

	pd, err := p.ParseProgram()
	if err != nil {
		log.Fatal(errors.Wrap(err, "could not parse definition"))
	}

	var definition []byte
	definition, err = json.MarshalIndent(pd, "", "  ")
	if err != nil {
		log.Fatal(errors.Wrap(err, "could not marshal definition"))
	}

	fmt.Println(string(definition))
}
```

#### Contributions

All contributions are welcome.

#### License

This project is licensed under the MIT License.