package pool

import (
	"fmt"
	"time"

	"github.com/cbergoon/pipes/pkg/pipeline"
)

// Worker represents the worker that executes the job
type Worker struct {
	ID int

	WorkerPool chan chan *Job
	JobChannel chan *Job
	quit       chan bool

	LivePipelineState bool

	DispatcherRef *Dispatcher
}

// WorkerError represents error encountered by worker.
type WorkerError struct {
	WorkerID int

	Error error

	ErrorTime    time.Time
	ErrorMessage string
	Content      string
}

// WorkerState Represents state information of reporting worker.
type WorkerState struct {
	WorkerID  int
	IsRunning bool
	StartTime time.Time

	IsRunningJob bool      //TODO (cbergoon): Implement
	JobStartTime time.Time //TODO (cbergoon): Implement
	JobEndTime   time.Time //TODO (cbergoon): Implement

	CurrentJob *Job

	Errors []*WorkerError
}

// NewWorker creates a new worker using WorkerPool from Dispatcher
func NewWorker(workerID int, livePipelineState bool, workerPool chan chan *Job, dispatcher *Dispatcher) Worker {
	return Worker{
		ID:                workerID,
		WorkerPool:        workerPool,
		JobChannel:        make(chan *Job),
		quit:              make(chan bool),
		DispatcherRef:     dispatcher,
		LivePipelineState: livePipelineState,
	}
}

func (w *Worker) initialize() error {
	return nil
}

// Start method starts the run loop of the worker
func (w *Worker) Start() {
	// register new state
	go func() {
		w.DispatcherRef.activeWorkers.Add(1)
		for {
			// register the current worker into the worker queue
			w.WorkerPool <- w.JobChannel

			select {
			case <-w.quit:
				fmt.Println("workerquit")
				w.DispatcherRef.activeWorkers.Done()
				return
			case job := <-w.JobChannel:
				var err error
				job.FlowPipeline, err = pipeline.NewPipelineFromPipelineDefinition(job.PipelineDefinition, w.DispatcherRef.plugins, w.LivePipelineState, w.LivePipelineState, job.PipelineCallback)
				if err != nil {
					fmt.Println(err)
					//TODO (cbergoon): Need to Add to Worker Pool Errors - Maybe Need New Type. Should this be Worker Pool Error or Job Error
				}
				job.FlowPipeline.Execute()
			}
		}
	}()
}

// Stop exits the run loop of the worker
func (w *Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}
