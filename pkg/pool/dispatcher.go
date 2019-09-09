package pool

import (
	"log"
	"plugin"
	"sync"

	"github.com/cbergoon/pipes/pkg/pipeline"
)

// Dispatcher handles routing jobs to workers
type Dispatcher struct {
	Workers    []*Worker
	WorkerPool chan chan *Job
	maxWorkers int

	plugins   map[string]*plugin.Plugin
	pluginDir string

	activeWorkers sync.WaitGroup

	JobQueue chan *Job

	quit chan bool

	LivePipelineState bool
}

// NewDispatcher creates a new Dispatcher
func NewDispatcher(maxW int, livePipelineState bool) *Dispatcher {
	return &Dispatcher{
		JobQueue:          make(chan *Job),
		WorkerPool:        make(chan chan *Job, maxW),
		maxWorkers:        maxW,
		pluginDir:         "",
		plugins:           make(map[string]*plugin.Plugin),
		quit:              make(chan bool),
		LivePipelineState: livePipelineState,
	}
}

// NewDispatcherWithPlugins creates a new Dispatcher with plugin search information
func NewDispatcherWithPlugins(maxW int, livePipelineState bool, pluginDir string) (*Dispatcher, error) {
	pluginMap, err := pipeline.LoadProcessPluginNameMap(pluginDir)
	if err != nil {
		return nil, err
	}
	return &Dispatcher{
		JobQueue:          make(chan *Job),
		WorkerPool:        make(chan chan *Job, maxW),
		maxWorkers:        maxW,
		pluginDir:         pluginDir,
		plugins:           pluginMap,
		quit:              make(chan bool),
		LivePipelineState: livePipelineState,
	}, err
}

// Run creates workers and builds pool
func (d *Dispatcher) Run() {
	for i := 0; i < d.maxWorkers; i++ {
		worker := NewWorker(i, d.LivePipelineState, d.WorkerPool, d)
		d.Workers = append(d.Workers, &worker)
		err := worker.initialize()
		if err != nil {
			log.Fatalf("failed to initialize workerpool: %+v", err)
		}
		worker.Start()
	}
	go d.dispatch()
}

// Run creates workers and builds pool
func (d *Dispatcher) RunBlocking() {
	for i := 0; i < d.maxWorkers; i++ {
		worker := NewWorker(i, d.LivePipelineState, d.WorkerPool, d)
		d.Workers = append(d.Workers, &worker)
		err := worker.initialize()
		if err != nil {
			log.Fatalf("failed to initialize workerpool: %+v", err)
		}
		worker.Start()
	}
	d.dispatch()
}

// Dispatcher.dispatch sends worker to execute a Job from the queue
func (d *Dispatcher) dispatch() {
	for {
		select {
		case <-d.quit:
			for _, w := range d.Workers {
				w.Stop()
			}
			d.activeWorkers.Wait()
			return
		case job := <-d.JobQueue:
			go func(job *Job) {
				jobChannel := <-d.WorkerPool
				jobChannel <- job
			}(job)

		}
	}
}

func (d *Dispatcher) Stop() {
	d.quit <- true
}
