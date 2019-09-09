package dl

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/cbergoon/pipes/pkg/pipeline"
)

func GenerateDLFromPipelineDefinition(definition *pipeline.PipelineDefinition) (string, error) {
	if definition.Pipeline == nil {
		return "", errors.Errorf("failed to create pipeline DL; pipeline is nil")
	}
	dl := generatePipelineInfoDL(*definition.Pipeline)
	for _, process := range definition.Processes {
		dl = fmt.Sprint(dl, generateProcessDL(*process))
	}
	for _, connection := range definition.Connections {
		dl = fmt.Sprint(dl, generateConnectDL(*connection))
	}
	return dl, nil
}

func generatePipelineInfoDL(pipelineInfo pipeline.PipelineInfoDefinition) string {
	return fmt.Sprint("CREATE PIPELINE \"", pipelineInfo.Name, "\";\n")
}

func generateProcessDL(process pipeline.ProcessDefinition) string {
	p := ""
	if process.Sink {
		p = fmt.Sprint("\tADD SINK \"", process.ProcessName, "\" OF \"", process.TypeName, "\"")
	} else {
		p = fmt.Sprint("\tADD \"", process.ProcessName, "\" OF \"", process.TypeName, "\"")
	}
	if len(process.Outputs) > 0 {
		p = fmt.Sprint(p, " \n\t\tOUTPUTS = (")
		for i, o := range process.Outputs {
			if i == 0 {
				p = fmt.Sprint(p, quoteEnclose(o))
			} else {
				p = fmt.Sprint(p, ", ", quoteEnclose(o))
			}
		}
		p = fmt.Sprint(p, ")")
	}
	if len(process.Inputs) > 0 {
		p = fmt.Sprint(p, " \n\t\tINPUTS = (")
		for i, in := range process.Inputs {
			if i == 0 {
				p = fmt.Sprint(p, quoteEnclose(in))
			} else {
				p = fmt.Sprint(p, ", ", quoteEnclose(in))
			}
		}
		p = fmt.Sprint(p, ")")
	}
	if len(process.State) > 0 {
		p = fmt.Sprint(p, " \n\t\tSET ")
		count := 0
		for k, v := range process.State {
			if count == 0 {
				p = fmt.Sprint(p, quoteEnclose(k), " = ", quoteEnclose(v))
			} else {
				p = fmt.Sprint(p, ", ", quoteEnclose(k), " = ", quoteEnclose(v))
			}
			count++
		}
	}
	p = fmt.Sprint(p, ";\n")
	return p
}

func generateConnectDL(connection pipeline.ConnectionDefinition) string {
	p := fmt.Sprint("\tCONNECT ", quoteEnclose(connection.OriginProcessName), ":", quoteEnclose(connection.OriginPortName), " TO ", quoteEnclose(connection.DestinationProcessName), ":", quoteEnclose(connection.DestinationPortName), ";\n")
	return p
}

func GenerateDLFromPipelineDefinitionJSON(definition []byte) (string, error) {
	pd, err := pipeline.NewPipelineDefinitionFromJson(definition)
	if err != nil {
		return "", err
	}
	dl, err := GenerateDLFromPipelineDefinition(pd)
	if err != nil {
		return "", err
	}
	return dl, nil
}

func quoteEnclose(unq string) string {
	return fmt.Sprint("\"", unq, "\"")
}
