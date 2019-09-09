<h1 align="center">Pipes Pool</h1>
<p align="center">
<a href="https://goreportcard.com/report/github.com/cbergoon/pipes-pool"><img src="https://goreportcard.com/badge/github.com/cbergoon/pipes-pool?1=1" alt="Report"></a>
<a href="https://godoc.org/github.com/cbergoon/pipes-pool"><img src="https://img.shields.io/badge/godoc-reference-brightgreen.svg" alt="Docs"></a>
<a href="#"><img src="https://img.shields.io/badge/version-0.1.0-brightgreen.svg" alt="Version"></a>
</p>

Pipes Pool provides a worker pool implementation to the Pipes library enabling concurrent execution of pipelines as well as the ability to control troughput and resource utilization.

For now Pipes Pool is a proof of concept and should not be used in production yet.

#### Features

* Concurrent execution of entire pipelines. 
* Access to Pipeline during execution.
* Ability to limit level of concurrency. 
* Graceful shutdown of processes.
* Worker state and Job information.
* Ability to register callback functions to be executed by workers. 

#### Installation

Get the source with ```go get```:

```bash
$ go get github.com/cbergoon/pipes-pool
```

Then import the library in your project:

```go
import "github.com/cbergoon/pipes"
import "github.com/cbergoon/pipes-dl"
import "github.com/cbergoon/pipes-pool"
```

#### Documentation

A Pipes Pool ...

#### Example Usage

```go
package main

import (
	"fmt"
	"log"
	"time"

	pipesdl "github.com/cbergoon/pipes-dl"
	pipespool "github.com/cbergoon/pipes-pool"
)

func main() {
	// MaxWorker controls number of workers available
	MaxWorker := 10

	// create and start Dispatcher
	dispatcher := pipespool.NewDispatcher(MaxWorker, true)
	dispatcher.Run()

	var jobs []*pipespool.Job
	for i := 0; i < 15; i++ {
		source := `CREATE PIPELINE "MyPipeline";

		ADD "Alfa" OF "Generator" OUTPUTS = ("Out1");
		ADD "Beta" OF "DynamicJS"
		  INPUTS = ("In1")
		  OUTPUTS = ("Out")
		  SET "src" = 'Out = In1;',
			"gg" = "kk";
		ADD SINK "Charlie" OF "Printer" INPUTS = ("In");
		
		CONNECT "Alfa":"Out1" TO "Beta":"In1";
		CONNECT "Beta":"Out" TO "Charlie":"In";`

		l := pipesdl.NewLexer(source)
		p := pipesdl.NewParser(l)

		pd, err := p.ParseProgram()
		if err != nil {
			log.Fatal(err)
		}

		work := &pipespool.Job{
			ID:                 int64(i),
			InitialState:       make(map[string]string),
			PipelineDefinition: pd,
		}
		jobs = append(jobs, work)
		go func() {
			for {
				if jobs[0].FlowPipeline != nil {
					fmt.Println(jobs[0].FlowPipeline.GetPipelineState())
					time.Sleep(time.Millisecond * 1000)
				}
			}
		}()
		dispatcher.JobQueue <- work
	}

	time.Sleep(time.Second * 10)
}
```

#### Contributions

All contributions are welcome.

#### License

This project is licensed under the MIT License.