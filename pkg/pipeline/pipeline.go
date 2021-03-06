package pipeline

import (
	"fmt"
	"os"
	"plugin"
	"sync"
	"time"

	"github.com/pkg/errors"
)

//TODO: Process with no output could potentially identified as SINK automatically eliminating explicit definition as such

type PipelineState struct {
	PipelineName  string
	Errors        []*PipelineError
	ProcessStates map[string]*ProcessState
}

type ProcessState struct {
	ExecutionId     string
	ProcessName     string
	ProcessTypeName string
	InitialState    map[string]string
	Inputs          map[string]string
	Outputs         map[string]string
}

type PipelineError struct {
	ProcessName     string
	ProcessTypeName string

	Error error

	ErrorTime    time.Time
	ErrorMessage string
	Content      string
	Inputs       map[string]string
}

func NewPipelineError(pName, pTypeName string, err error, errTime time.Time, errMessage, content string, inputs map[string]string) PipelineError {
	return PipelineError{
		ProcessName:     pName,
		ProcessTypeName: pTypeName,
		Error:           err,
		ErrorTime:       errTime,
		ErrorMessage:    errMessage,
		Content:         content,
		Inputs:          inputs,
	}
}

type Pipeline interface {
	Initialize()
	AddProcess(name string, process Process, sink bool) error
	Connect(originProcess, originPort string, destinationProcess, destinationPort string) error
	Execute()
	AddError(err *PipelineError)
	ProcessStateChanged(state *ProcessState)
	IsProcessStateChangedEnabled() bool
	GetPipelineState() *PipelineState
}

type FlowPipeline struct {
	Name      string
	Sink      string
	Processes map[string]Process

	Errors []*PipelineError

	StateChangeEnable           bool
	StateChangedCallbacksEnable bool
	StateChangedCallbackFn      func(state PipelineState)
	stateMux                    sync.RWMutex
	State                       *PipelineState
}

func (f *FlowPipeline) Initialize() {
	for _, proc := range f.Processes {
		proc.Initialize()
	}
}

func NewPipeline(name string, stateChangedEnabled bool, stateChangedCallbackEnabled bool, stateChangedCallbackFn func(state PipelineState)) *FlowPipeline {
	state := &PipelineState{
		PipelineName:  name,
		ProcessStates: make(map[string]*ProcessState),
	}
	return &FlowPipeline{
		Name:                        name,
		Processes:                   make(map[string]Process),
		Sink:                        "",
		StateChangeEnable:           stateChangedEnabled,
		StateChangedCallbacksEnable: stateChangedCallbackEnabled,
		StateChangedCallbackFn:      stateChangedCallbackFn,
		State:                       state,
	}
}

func NewPipelineFromPipelineDefinition(definition *PipelineDefinition, pluginMap map[string]*plugin.Plugin, stateChangedEnabled bool, stateChangedCallbackEnabled bool, stateChangedCallbackFn func(state PipelineState)) (*FlowPipeline, error) {
	pipeline := NewPipeline(definition.Pipeline.Name, stateChangedEnabled, stateChangedCallbackEnabled, stateChangedCallbackFn)
	for _, proc := range definition.Processes {
		if _, ok := pluginMap[proc.TypeName]; ok {
			plug := pluginMap[proc.TypeName]

			symLogger, err := plug.Lookup("New")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			plugedNewProcess, ok := symLogger.(func(processName string, inputs, outputs []string, state map[string]string) Process)
			if !ok {
				fmt.Println("unexpected type from module symbol")
				os.Exit(1)
			}

			err = pipeline.AddProcess(proc.ProcessName, plugedNewProcess(proc.ProcessName, proc.Inputs, proc.Outputs, proc.State), proc.Sink)
			if err != nil {
				return nil, errors.Wrapf(err, "could not create pipeline; specified process %s of type %s could not be created", proc.ProcessName, proc.TypeName)
			}
		} else if proc.TypeName == "Http" {
			err := pipeline.AddProcess(proc.ProcessName, NewHttpProcess(proc.ProcessName, proc.Inputs, proc.Outputs, proc.State), proc.Sink)
			if err != nil {
				return nil, errors.Wrapf(err, "could not create pipeline; specified process %s of type %s could not be created", proc.ProcessName, proc.TypeName)
			}
		} else if proc.TypeName == "JSONFileReader" {
			err := pipeline.AddProcess(proc.ProcessName, NewJSONFileReaderProcess(proc.ProcessName, proc.Inputs, proc.Outputs, proc.State), proc.Sink)
			if err != nil {
				return nil, errors.Wrapf(err, "could not create pipeline; specified process %s of type %s could not be created", proc.ProcessName, proc.TypeName)
			}
		} else if proc.TypeName == "Printer" {
			err := pipeline.AddProcess(proc.ProcessName, NewPrinterProcess(proc.ProcessName, proc.Inputs, proc.Outputs, proc.State), proc.Sink)
			if err != nil {
				return nil, errors.Wrapf(err, "could not create pipeline; specified process %s of type %s could not be created", proc.ProcessName, proc.TypeName)
			}
		} else if proc.TypeName == "Generator" {
			err := pipeline.AddProcess(proc.ProcessName, NewGeneratorProcess(proc.ProcessName, proc.Inputs, proc.Outputs, proc.State), proc.Sink)
			if err != nil {
				return nil, errors.Wrapf(err, "could not create pipeline; specified process %s of type %s could not be created", proc.ProcessName, proc.TypeName)
			}
		} else if proc.TypeName == "DynamicJS" {
			err := pipeline.AddProcess(proc.ProcessName, NewDynamicJsProcess(proc.ProcessName, proc.Inputs, proc.Outputs, proc.State), proc.Sink)
			if err != nil {
				return nil, errors.Wrapf(err, "could not create pipeline; specified process %s of type %s could not be created", proc.ProcessName, proc.TypeName)
			}
		} else {
			return nil, errors.Errorf("could not create pipeline; specified process of type %s does not exist", proc.TypeName)
		}
	}
	pipeline.Initialize()
	for _, conn := range definition.Connections {
		err := pipeline.Connect(conn.OriginProcessName, conn.OriginPortName, conn.DestinationProcessName, conn.DestinationPortName)
		if err != nil {
			return nil, errors.Wrapf(err, "could not create pipeline; connection from %s:%s to %s:%s could not be created", conn.OriginProcessName, conn.OriginPortName, conn.DestinationProcessName, conn.DestinationPortName)
		}
	}
	return pipeline, nil
}

func (f *FlowPipeline) AddProcess(name string, process Process, sink bool) error {
	process.GetFlowProcess().PipelineRef = f
	if _, ok := f.Processes[name]; ok {
		return errors.Errorf("could not add process due to duplicate name %s", name)
	}
	f.Processes[name] = process
	if sink {
		f.Sink = name
	}
	return nil
}

func (f *FlowPipeline) Connect(originProcess, originPort string, destinationProcess, destinationPort string) error {
	//Get Origin Process
	op, ok := f.Processes[originProcess]
	if !ok {
		return errors.Errorf("could not connect origin process %s not found", originProcess)
	}
	//Get Destination Process
	dp, ok := f.Processes[destinationProcess]
	if !ok {
		return errors.Errorf("could not connect destination process %s not found", destinationProcess)
	}
	//Set the Input Channel of the Destination Process to the Output Channel of the Origin Process
	dp.GetFlowProcess().SetInputChannelByName(destinationPort, op.GetFlowProcess().GetOutputChannelByName(originPort))

	//Single Line Equivalent w/o Error Handling
	//f.Processes[destinationProcess].GetFlowProcess().SetInputChannelByName(destinationPort, f.Processes[originProcess].GetFlowProcess().GetOutputChannelByName(originPort))

	return nil
}

func (f *FlowPipeline) Execute() {
	for procName, proc := range f.Processes {
		if f.Sink != procName {
			go proc.Run()
		}
	}
	sp := f.Processes[f.Sink]
	sp.Run()
}

func (f *FlowPipeline) AddError(err *PipelineError) {
	f.Errors = append(f.Errors, err)
}

func (f *FlowPipeline) ProcessStateChanged(state *ProcessState) {
	if f.StateChangeEnable {
		f.stateMux.Lock()
		f.State.ProcessStates[state.ProcessName] = state
		f.State.Errors = f.Errors
		f.stateMux.Unlock()
		if f.StateChangedCallbackFn != nil {
			f.stateMux.RLock()
			f.StateChangedCallbackFn(*f.State)
			f.stateMux.RUnlock()
		}
	}
}

func (f *FlowPipeline) IsProcessStateChangedEnabled() bool {
	return f.StateChangedCallbacksEnable
}

func (f *FlowPipeline) GetPipelineState() *PipelineState {
	f.stateMux.RLock()
	defer f.stateMux.RUnlock()
	return f.State
}
