package pool

import "github.com/cbergoon/pipes/pkg/pipeline"

// Job defines the task to be executed by worker
type Job struct {
	ID                 int64
	InitialState       map[string]string
	PipelineDefinition *pipeline.PipelineDefinition
	FlowPipeline       *pipeline.FlowPipeline
	PipelineCallback   func(state pipeline.PipelineState)
}
