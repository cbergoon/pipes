package pipes

type Process interface {
	Initialize()
	GetFlowProcess() *FlowProcess
	Run()
}

type FlowProcess struct {
	TypeName    string
	Inputs      map[string]chan string
	Outputs     map[string]chan string
	PipelineRef Pipeline
}

func NewFlowProcess(typeName string) *FlowProcess {
	f := &FlowProcess{}
	f.TypeName = typeName
	f.Outputs = make(map[string]chan string)
	f.Inputs = make(map[string]chan string)
	return f
}

func (f *FlowProcess) GetTypeName() string {
	return f.TypeName
}

func (f *FlowProcess) SetTypeName(typeName string) {
	f.TypeName = typeName
}

func (f *FlowProcess) AddInput(name string) {
	f.Inputs[name] = nil
}

func (f *FlowProcess) AddOutput(name string) {
	f.Outputs[name] = make(chan string, 1)
}

func (f *FlowProcess) GetInputChannelMap() map[string]chan string {
	return f.Inputs
}

func (f *FlowProcess) GetInputChannelByName(name string) chan string {
	return f.Inputs[name]
}

func (f *FlowProcess) SetInputChannelByName(name string, port chan string) {
	f.Inputs[name] = port
}

func (f *FlowProcess) AddInputChannel(name string, input chan string) {
	f.Inputs[name] = input
}

func (f *FlowProcess) GetOutputChannelMap() map[string]chan string {
	return f.Outputs
}

func (f *FlowProcess) GetOutputChannelByName(name string) chan string {
	return f.Outputs[name]
}

func (f *FlowProcess) SetOutputChannelByName(name string, port chan string) {
	f.Outputs[name] = port
}

func (f *FlowProcess) AddOutputChannel(name string, output chan string) {
	f.Outputs[name] = output
}

func (f *FlowProcess) AddError(err *PipelineError) {
	f.PipelineRef.AddError(err)
}

func (f *FlowProcess) ProcessStateChanged(state *ProcessState) {
	f.PipelineRef.ProcessStateChanged(state)
}

func NewProcessByTypeName(typeName, processName string, inputs, outputs []string, state map[string]string) (Process, error) {
	//TODO: Build Factory
	return nil, nil
}
