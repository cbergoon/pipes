<h1 align="center">Pipes</h1>
<p align="center">
<a href="https://goreportcard.com/report/github.com/cbergoon/pipes"><img src="https://goreportcard.com/badge/github.com/cbergoon/pipes?1=1" alt="Report"></a>
<a href="https://godoc.org/github.com/cbergoon/pipes"><img src="https://img.shields.io/badge/godoc-reference-brightgreen.svg" alt="Docs"></a>
<a href="#"><img src="https://img.shields.io/badge/version-0.1.0-brightgreen.svg" alt="Version"></a>
</p>

Pipes provides the ability to rapidly define an application using prebuilt components (processes) that are dynamically
defined.

For now Pipes is a proof of concept and should not be used in production yet.

#### Features

* Concurrent execution of pipeline paths.
* Dynamic Javascript process.
* Prebuilt start shapes: HTTP, JSON FILE, Static Generator.
* Pipeline definable with JSON.
* State tracking of process blocks.
* Customizable state changed handler.
* Process level error reporting.
* FUTURE: External plugin system for process blocks running as their own process in any language.
* FUTURE: Built-in database processes.

NEW: A definition language for pipes called [pipes-dl](https://github.com/cbergoon/pipes-dl) is also available. This provides a simple
DSL that can be used to define a pipeline.

![flow](/docs/images/pipes-diagram.png)

#### Installation

Get the source with ```go get```:

```bash
$ go get github.com/cbergoon/pipes
```

Then import the library in your project:

```go
import "github.com/cbergoon/pipes"
```

#### Documentation

A Pipes pipeline consists of two main concepts: processes and connections. As you might have guessed a processes are
are linked and communicate via connections to form a pipeline. These connections also define the process graph which
define the flow of messages and execution through the pipeline.

##### Processes

Processes are the main parts of a process. These are similar to functions in a traditional program and define the logic
of the pipeline. Processes consist of a name, type, input ports, output ports, and a state.

The type of the process specifies which of the built in types the process should use. An example of type is an HTTP
process which make HTTP requests.

Inputs and outputs are named "ports" that the processes use to communicate.

The state of a process is set of definable initial data which is specifically defined per process instance.

There are currently four built-in process types: HTTP, JSON, DYNAMICJS, and GENERATOR.

##### Connections

Connections define the flow of the pipeline. A complete pipeline's connections will form a subset of a p-graph where only
one start and end vertex exists. Connections pass JSON data.

##### Pipelines

Pipelines represent the entire flow through the application.

#### Example Usage

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, Pipes")
}
```

#### Example Pipes Definition Language

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

#### Contributions

All contributions are welcome.

#### License

This project is licensed under the MIT License.