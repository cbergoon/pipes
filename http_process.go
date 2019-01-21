package pipes

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

type HttpProcess struct {
	Process
	*FlowProcess

	ProcessName string

	InitialState map[string]string

	Inputs  []string
	Outputs []string
}

type HttpProcessConfig struct {
	Url     string
	Method  string
	Headers map[string]string
	Body    string
}

func NewHttpProcess(processName string, inputs, outputs []string, state map[string]string) *HttpProcess {
	return &HttpProcess{
		Inputs:       inputs,
		Outputs:      outputs,
		InitialState: state,
		ProcessName:  processName,
		FlowProcess:  NewFlowProcess("Http"),
	}
}

func (c *HttpProcess) Run() {
	if len(c.Inputs) > 0 {
		for configJson := range c.GetInputChannelByName(c.Inputs[0]) {
			config := &HttpProcessConfig{}
			err := json.Unmarshal([]byte(configJson), config)
			if err != nil {
				//TODO: Error
			}

			if config.Method == "GET" || config.Method == "POST" || config.Method == "PUT" || config.Method == "PATCH" {
				url := config.Url

				var contentBuffer *bytes.Buffer
				if config.Method == "GET" {
					contentBuffer = nil
				} else {
					contentBuffer = bytes.NewBuffer([]byte(config.Body))
				}

				req, err := http.NewRequest(config.Method, url, contentBuffer)
				for k, v := range config.Headers {
					req.Header.Set(k, v)
				}

				client := &http.Client{}
				client.Timeout = time.Second * 60
				resp, err := client.Do(req)
				if err != nil {
					//TODO: Error
				}

				body, _ := ioutil.ReadAll(resp.Body)
				c.FlowProcess.Outputs[c.Outputs[0]] <- string(body)

				err = resp.Body.Close()
				if err != nil {
					//TODO: Error
				}
			}
		}
	} else {
		//Todo: Implement errors
		//c.AddError(NewPipelineError(c.ProcessName, c.TypeName, errors.New("Testing this error shit out."), time.Now(), "Camerons process error message", "Content String", make(map[string]string)))
		//for _, v := range c.FlowProcess.Outputs {
		//	close(v)
		//}
		//if true {
		//	return
		//}
		configJson, ok := c.InitialState["config"]
		if !ok {
			//TODO: Error
		}
		config := &HttpProcessConfig{}
		err := json.Unmarshal([]byte(configJson), config)
		if err != nil {
			//TODO: Error
		}

		if config.Method == "GET" || config.Method == "POST" || config.Method == "PUT" || config.Method == "PATCH" {
			url := config.Url

			var req *http.Request
			var err error
			if config.Method == "GET" {
				req, err = http.NewRequest(config.Method, url, nil)
			} else {
				var contentBuffer *bytes.Buffer
				contentBuffer = bytes.NewBuffer([]byte(config.Body))
				req, err = http.NewRequest(config.Method, url, contentBuffer)
			}

			for k, v := range config.Headers {
				req.Header.Set(k, v)
			}

			client := &http.Client{}
			client.Timeout = time.Second * 60
			resp, err := client.Do(req)
			if err != nil {
				//TODO: Error
			}

			body, _ := ioutil.ReadAll(resp.Body)
			c.FlowProcess.Outputs[c.Outputs[0]] <- string(body)

			err = resp.Body.Close()
			if err != nil {
				//TODO: Error
			}
		}
	}

	for _, v := range c.FlowProcess.Outputs {
		close(v)
	}
}

func (c *HttpProcess) Initialize() {
	//Build output ports; Inputs are implicit based on connections
	for _, v := range c.Inputs {
		c.FlowProcess.AddInput(v)
	}
	for _, v := range c.Outputs {
		c.FlowProcess.AddOutput(v)
	}
}

func (c *HttpProcess) GetFlowProcess() *FlowProcess {
	return c.FlowProcess
}

func (c *HttpProcess) CloseOutputs() {
	for _, v := range c.FlowProcess.Outputs {
		close(v)
	}
}
