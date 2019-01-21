package pipes

import (
	"encoding/json"
	"github.com/pkg/errors"
)

type PipelineInfoDefinition struct {
	Name string `json:"name"`
}

type ProcessDefinition struct {
	TypeName    string            `json:"typeName"`
	ProcessName string            `json:"processName"`
	Sink        bool              `json:"sink"`
	Inputs      []string          `json:"inputs"`
	Outputs     []string          `json:"outputs"`
	State       map[string]string `json:"state"`
}

type ConnectionDefinition struct {
	OriginProcessName      string `json:"originProcessName"`
	OriginPortName         string `json:"originPortName"`
	DestinationProcessName string `json:"destinationProcessName"`
	DestinationPortName    string `json:"destinationPortName"`
}

type PipelineDefinition struct {
	Pipeline    *PipelineInfoDefinition `json:"pipeline"`
	Processes   []*ProcessDefinition    `json:"processes"`
	Connections []*ConnectionDefinition `json:"connections"`
}

func NewPipelineDefinitionFromJson(definition []byte) (*PipelineDefinition, error) {
	pipelineDefinition := &PipelineDefinition{}
	err := json.Unmarshal(definition, pipelineDefinition)
	if err != nil {
		return nil, errors.Wrap(err, "could not unmarshal pipeline definition")
	}
	return pipelineDefinition, nil
}
