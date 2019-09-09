package main

import (
	"fmt"
	"time"

	"github.com/cbergoon/pipes/pkg/pipeline"
)

// PluginTypeName defines the "TypeName" used to identify the correct processes component
// to be instantiated in th pipeline.
const PluginTypeName string = "Logger"

// LoggerProcess defines the process component which will log messages passed to inputs.
type LoggerProcess struct {
	pipeline.Process
	*pipeline.FlowProcess

	ProcessName string

	InitialState map[string]string

	Inputs  []string
	Outputs []string
}

// TypeName returns PluginTypeName.
func TypeName() string {
	return PluginTypeName
}

// New returns a new instance of the logger process component. Plugins must define a New
// function matching this description.
func New(processName string, inputs, outputs []string, state map[string]string) pipeline.Process {
	return &LoggerProcess{
		Inputs:       inputs,
		Outputs:      outputs,
		InitialState: state,
		FlowProcess:  pipeline.NewFlowProcess(PluginTypeName),
		ProcessName:  processName,
	}
}

// Run executes the logger process component printing input messages allong with a timestamp.
func (c *LoggerProcess) Run() {
	for s := range c.GetInputChannelByName("In") {
		if c.PipelineRef.IsProcessStateChangedEnabled() {
			inputs := make(map[string]string)
			inputs["In"] = s
			c.ProcessStateChanged(&pipeline.ProcessState{
				ProcessName:     c.ProcessName,
				ProcessTypeName: c.TypeName,
				Inputs:          inputs,
			})
		}
		fmt.Println(time.Now(), " ", s)
	}
}

// Initialize handles any required setup for the logger process component and creates
// connections between processes.
func (c *LoggerProcess) Initialize() {
	for _, v := range c.Inputs {
		c.FlowProcess.AddInput(v)
	}
	for _, v := range c.Outputs {
		c.FlowProcess.AddOutput(v)
	}
}

// GetFlowProcess returns core flow process reference.
func (c *LoggerProcess) GetFlowProcess() *pipeline.FlowProcess {
	return c.FlowProcess
}
