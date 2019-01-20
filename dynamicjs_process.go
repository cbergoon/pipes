package pipes

import (
	"github.com/robertkrimen/otto"
)

type DynamicJsProcess struct {
	Process
	*FlowProcess

	ProcessName string

	InitialState map[string]string

	Inputs  []string
	Outputs []string
}

func NewDynamicJsProcess(processName string, inputs, outputs []string, state map[string]string) *DynamicJsProcess {
	return &DynamicJsProcess{
		Inputs:       inputs,
		Outputs:      outputs,
		InitialState: state,
		ProcessName:  processName,
		FlowProcess:  NewFlowProcess("DynamicJS"),
	}
}

func (c *DynamicJsProcess) Run() {
	var inputData map[string]string
	src := c.InitialState["src"]
	for {
		closedCount := 0
		inputData = make(map[string]string)
		for k, v := range c.FlowProcess.Inputs {
			tmp, ok := <-v
			if !ok {
				closedCount++
			}
			inputData[k] = tmp
		}
		if closedCount == len(c.Inputs) {
			break
		}
		vm := otto.New()
		for k, v := range inputData {
			vm.Set(k, v)
		}
		vm.Run(src)
		for _, v := range c.Outputs {
			val, _ := vm.Get(v)
			c.FlowProcess.Outputs[v] <- val.String()
		}
	}

	for _, v := range c.FlowProcess.Outputs {
		close(v)
	}
}

func (c *DynamicJsProcess) Initialize() {
	//Build output ports; Inputs are implicit based on connections
	for _, v := range c.Inputs {
		c.FlowProcess.AddInput(v)
	}
	for _, v := range c.Outputs {
		c.FlowProcess.AddOutput(v)
	}
}

func (c *DynamicJsProcess) GetFlowProcess() *FlowProcess {
	return c.FlowProcess
}
