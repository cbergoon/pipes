package main

import (
	"fmt"
	"log"
	"time"

	"github.com/cbergoon/pipes/pkg/dl"
	"github.com/cbergoon/pipes/pkg/pool"
)

func main() {
	// MaxWorker controls number of workers available
	MaxWorker := 10

	// create and start Dispatcher
	dispatcher := pool.NewDispatcher(MaxWorker, true)

	go func() {
		var jobs []*pool.Job
		for i := 0; i < 1; i++ {
			source := `CREATE PIPELINE "MyPipeline";
	
			ADD "Alfa" OF "Generator" OUTPUTS = ("Out1");
			ADD "Beta" OF "DynamicJS"
			  INPUTS = ("In1")
			  OUTPUTS = ("Out")
			  SET "src" = 'Out = In1;',
				"gg" = "kk";
			ADD SINK "Charlie" OF "Logger" INPUTS = ("In");
			
			CONNECT "Alfa":"Out1" TO "Beta":"In1";
			CONNECT "Beta":"Out" TO "Charlie":"In";`

			l := dl.NewLexer(source)
			p := dl.NewParser(l)

			pd, err := p.ParseProgram()
			if err != nil {
				log.Fatal(err)
			}

			work := &pool.Job{
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
	}()

	dispatcher.RunBlocking()

}
