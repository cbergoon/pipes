<h1 align="center">Pipes</h1>
<p align="center">
<a href="https://goreportcard.com/report/github.com/cbergoon/pipes"><img src="https://goreportcard.com/badge/github.com/cbergoon/pipes?1=1" alt="Report"></a>
<a href="https://godoc.org/github.com/cbergoon/pipes"><img src="https://img.shields.io/badge/godoc-reference-brightgreen.svg" alt="Docs"></a>
<a href="#"><img src="https://img.shields.io/badge/version-0.1.0-brightgreen.svg" alt="Version"></a>
</p>

Pipes provides the ability to rapidly define an application using prebuilt components (processes) that are linked to form pipelines. Processes and Pipelines can 
be defined using a definition language called pipesdl. Pipes is a "pluggable" system meaning you can create your own process components to accomplish specialized tasks.

The goal of pipes is to provide a black box programming model that enables developers and non-devlopers to efficiently build processes tailored to a specific use case. 

For now Pipes is a proof of concept and should not be used in production.

### Features

* Concurrent execution of pipeline paths.
* Dynamic Javascript process.
* Prebuilt entrypoint components: HTTP, JSON FILE, CSV (TODO), FILE (TODO), DIRECTORY (TODO), ...
* Pipeline definable with JSON.
* State tracking of process blocks.
* Customizable state changed handler.
* Process level error reporting.
* Worker pool to enable concurrent execution of pipelines.
* FUTURE: Plugin system to create custom process components.
* FUTURE: Robust daemon with CLI tooling. 
* FUTURE: External plugin system for process blocks running as their own process in any language.
* FUTURE: Server application providing visual composition of processes and pipelines, execution management, scheduling, deployment, etc. 
* FUTURE: Connections that provide control   structures. 

![flow](/docs/images/pipes-diagram.png)

### Installation

Get the source with ```go get```:

```bash
$ go get github.com/cbergoon/pipes
```

Then import the library in your project:

```go
import "github.com/cbergoon/pipes/pkg/pipeline"
import "github.com/cbergoon/pipes/pkg/dl"
import "github.com/cbergoon/pipes/pkg/pool"
```

---

### Documentation

A Pipes pipeline consists of two main concepts: processes and connections. As you might have guessed a processes are
are linked and communicate via connections to form a pipeline. These connections also define the process graph which
define the flow of messages and execution through the pipeline.

Also see the [godocs](https://godoc.org/github.com/cbergoon/pipes) and the [docs](/docs/README.md) directory.

### Pipelines

Pipelines represent the entire flow through the application formed by a group of Processes and Connections. Pipelines always
start with a generator type (a process with only outputs) and always end with a 'sink' type (a process with only inputs). 
Process components that make up the middle have both inputs and outputs with varying degrees of branching. 

Pipelines may be completely sequential, or branch off into concurrently executed components. All pipelines have at most one 
'generator' and one 'sink'. 

### Processes

In pipes, processes are components that handle specific operations or tasks. These building blocks are composable and 
customizable. Processes have similar characteristics to a function and can be connected and invoked by other components 
in the pipeline. Processes are the main piece of a pipeline. Processes consist of a name, type, input ports, output ports, 
and a state.

The type of the process specifies which of the built in types the process should use. An example of type is an HTTP
process which make HTTP requests. The process type is also used to identify process plugins. 

Inputs and outputs in pipes process components are often refered to as "ports". Processes components use ports to communicate 
and indicate stage completion. 

The state of a process is set of definable initial data which is specifically defined per process instance.

There are currently four built-in process types: HTTP, JSON, DYNAMICJS, and GENERATOR. Additional official plugins can be 
found in the /plugins directory at the root of the project. In the future, we hope to have a community plugin repository. 

### Connections

Connections define the flow of the pipeline. A complete pipeline's connections will form a subset of a p-graph where only
one start and end vertex exists. Connections pass JSON data amongst process blocks. 

Connections are defined in pipesdl by enumerating which input port the output port of a preceeding process should connect to. 

Connections have a secondary purpose to communication, they define the pipeline flow and in the future will provide control-flow
primitives. 

Process ports are blocking which means that if a process has 4 input ports all 4 MUST receive a message in order for the process 
to continue.

### Plugins

Custom processes can be built by creating a plugin that implements the process API. More information about using and creating custom 
processes can be found [here](/docs/plugins.md).

### Definition Language

The definition language allows entire pipelines to be scripted. This provides a human readable representation and a way to preserve 
and replicate pipelines. A readable language also simplifies reuse of process components. 

### Worker Pool

The worker pool provides a safe way to execute many pipelines at once with varying workloads in a way that preserves and extends the 
existing pipeline API and state. 

---

### Example Usage

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, Pipes")
}
```

### Example Pipes Definition Language

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

### Contributions

All contributions are welcome.

### License

This project is licensed under the MIT License.