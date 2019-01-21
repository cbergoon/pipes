package pipes

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"
)

const (
	S_TYPE_OBJECT = "OBJECT"
	S_TYPE_ARRAY  = "ARRAY"

	P_TYPE_SPLIT    = "SPLIT"
	P_TYPE_NO_SPLIT = "NOSPLIT"
)

type JSONFileReaderProcess struct {
	Process
	*FlowProcess

	ProcessName string

	InitialState map[string]string

	Inputs  []string
	Outputs []string
}

func NewJSONFileReaderProcess(processName string, inputs, outputs []string, state map[string]string) *JSONFileReaderProcess {
	return &JSONFileReaderProcess{
		Inputs:       inputs,
		Outputs:      outputs,
		InitialState: state,
		ProcessName:  processName,
		FlowProcess:  NewFlowProcess("JSONFileReader"),
	}
}

func (c *JSONFileReaderProcess) Run() {
	sType, ok := c.InitialState["sType"]
	if !ok {

	}
	pType, ok := c.InitialState["pType"]
	if !ok {

	}

	dir, ok := c.InitialState["directory"]
	if !ok {

	}
	fileFilter, ok := c.InitialState["fileFilter"]
	if !ok {

	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		matched, _ := filepath.Match(fileFilter, f.Name())
		if matched {
			content, err := ioutil.ReadFile(f.Name())
			if err != nil {
				//add error and continue
			}
			if sType == S_TYPE_OBJECT {
				var objobj *json.RawMessage
				err := json.Unmarshal([]byte(content), &objobj)
				if err != nil {
					//add error and continue
				}
				contentItem, err := objobj.MarshalJSON()
				if err != nil {
					//add error and continue
				}
				c.FlowProcess.Outputs[c.Outputs[0]] <- string(contentItem)
			} else if sType == S_TYPE_ARRAY {
				if pType == P_TYPE_NO_SPLIT {
					var objobj *json.RawMessage
					err := json.Unmarshal([]byte(content), &objobj)
					if err != nil {
						//add error and continue
					}
					contentItem, err := objobj.MarshalJSON()
					if err != nil {
						//add error and continue
					}
					c.FlowProcess.Outputs[c.Outputs[0]] <- string(contentItem)
				} else if pType == P_TYPE_SPLIT {
					var objarr []*json.RawMessage
					err := json.Unmarshal([]byte(content), &objarr)
					if err != nil {
						//add error and continue
					}
					for _, j := range objarr {
						contentItem, err := j.MarshalJSON()
						if err != nil {
							//add error and continue
						}
						c.FlowProcess.Outputs[c.Outputs[0]] <- string(contentItem)
					}
				}
			}
		}
	}

	for _, v := range c.FlowProcess.Outputs {
		close(v)
	}
}

func (c *JSONFileReaderProcess) Initialize() {
	for _, v := range c.Outputs {
		c.FlowProcess.AddOutput(v)
	}
}

func (c *JSONFileReaderProcess) GetFlowProcess() *FlowProcess {
	return c.FlowProcess
}
