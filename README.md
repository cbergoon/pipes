<h1 align="center">Pipes</h1>
<p align="center">
<a href="https://goreportcard.com/report/github.com/cbergoon/pipes"><img src="https://goreportcard.com/badge/github.com/cbergoon/pipes-dl?1=1" alt="Report"></a>
<a href="https://godoc.org/github.com/cbergoon/pipes"><img src="https://img.shields.io/badge/godoc-reference-brightgreen.svg" alt="Docs"></a>
<a href="#"><img src="https://img.shields.io/badge/version-0.1.0-brightgreen.svg" alt="Version"></a>
</p>

Pipes provides the ability to rapidly define an application using prebuilt components (processes) that are dynamically
defined. For now Pipes is just an experiment.

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

NEW: A definition language for pipes called [pipes-dl](github.com/cbergoon/pipes-dl) is also available. This provides a simple
DSL that can be used to define a pipeline.

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

###### JSON File Process

###### HTTP Process

###### DynamicJs Process

##### Connections

##### Pipelines

##### JSON Pipeline Definition

#### Example Usage

```go
    package main

    import "fmt"

    func main() {
        fmt.Println("Hello, Pipes")
    }
```

#### Contributions

All contributions are welcome.

#### License

This project is licensed under the MIT License.