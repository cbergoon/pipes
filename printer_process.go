package pipes

import (
	"fmt"
)

type PrinterProcess struct {
	Process
	*FlowProcess

	ProcessName string

	InitialState map[string]string

	Inputs  []string
	Outputs []string
}

func NewPrinterProcess(processName string, inputs, outputs []string, state map[string]string) *PrinterProcess {
	return &PrinterProcess{
		Inputs:       inputs,
		Outputs:      outputs,
		InitialState: state,
		FlowProcess:  NewFlowProcess("Printer"),
		ProcessName:  processName,
	}
}

func (c *PrinterProcess) Run() {
	for s := range c.GetInputChannelByName("In") {
		if c.PipelineRef.IsProcessStateChangedEnabled() {
			inputs := make(map[string]string)
			inputs["In"] = s
			c.ProcessStateChanged(&ProcessState{
				ProcessName:     c.ProcessName,
				ProcessTypeName: c.TypeName,
				Inputs:          inputs,
			})
		}
		fmt.Println(s)
	}
}

func (c *PrinterProcess) Initialize() {
	for _, v := range c.Inputs {
		c.FlowProcess.AddInput(v)
	}
	for _, v := range c.Outputs {
		c.FlowProcess.AddOutput(v)
	}
}

func (c *PrinterProcess) GetFlowProcess() *FlowProcess {
	return c.FlowProcess
}
