package pipes

import (
	"fmt"
)

type GeneratorProcess struct {
	Process
	*FlowProcess

	ProcessName string

	InitialState map[string]string

	Inputs      []string
	Outputs     []string
}

func NewGeneratorProcess(processName string, inputs, outputs []string, state map[string]string) *GeneratorProcess {
	return &GeneratorProcess{
		Inputs: inputs,
		Outputs: outputs,
		FlowProcess: NewFlowProcess("Generator"),
		ProcessName: processName,
	}
}

func (c *GeneratorProcess) Run() {
	for i := 1; i <= 5; i++ {
		c.FlowProcess.Outputs["Out1"] <- fmt.Sprintf("Hi for the %d'th time!", i)
		//c.FlowProcess.Outputs["Out2"] <- fmt.Sprintf("Cameron")
	}

	for _, v := range c.FlowProcess.Outputs {
		close(v)
	}
}

func (c *GeneratorProcess) Initialize() {
	//Build output ports; Inputs are implicit based on connections
	for _, v := range c.Inputs {
		c.FlowProcess.AddInput(v)
	}
	for _, v := range c.Outputs {
		c.FlowProcess.AddOutput(v)
	}
}

func (c *GeneratorProcess) GetFlowProcess() *FlowProcess {
	return c.FlowProcess
}
