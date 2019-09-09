package dl

import (
	"github.com/cbergoon/pipes/pkg/pipeline"
	"github.com/pkg/errors"
)

type Parser struct {
	l      *Lexer
	errors []string

	curToken  Token
	peekToken Token
}

func NewParser(l *Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.Initialize()

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) Initialize() {
	p.nextToken()
	p.nextToken()
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() (*pipeline.PipelineDefinition, error) {
	definition := &pipeline.PipelineDefinition{}
	definition.Processes = []*pipeline.ProcessDefinition{}
	definition.Connections = []*pipeline.ConnectionDefinition{}

	for p.curToken.Type != EOF {
		if p.curToken.Type == CREATE {
			def, err := p.parseCreatePipelineStatement()
			if err != nil {
				return nil, errors.Wrap(err, "invalid syntax: failed to parse pipeline statement")
			}
			definition.Pipeline = def
		} else if p.curToken.Type == ADD {
			def, err := p.parseCreateProcessStatement()
			if err != nil {
				return nil, errors.Wrap(err, "invalid syntax: failed to parse add process statement")
			}
			definition.Processes = append(definition.Processes, def)
		} else if p.curToken.Type == CONNECT {
			def, err := p.parseCreateConnectionStatement()
			if err != nil {
				return nil, errors.Wrap(err, "invalid syntax: failed to parse create connection statement")
			}
			definition.Connections = append(definition.Connections, def)
		}
		p.nextToken()
	}

	return definition, nil
}

func (p *Parser) parseCreatePipelineStatement() (*pipeline.PipelineInfoDefinition, error) {
	pipelineInfo := &pipeline.PipelineInfoDefinition{}
	if p.curToken.Type == CREATE {
		p.nextToken()
		if p.curToken.Type == PIPELINE {
			p.nextToken()
			if p.curToken.Type == STRING {
				pipelineInfo.Name = p.curToken.Literal
			} else {
				return nil, errors.Errorf("syntax error: expected string literal of <pipeline name> after %s %s; found %s of type %s", p.curToken.Literal, p.curToken.Type)
			}
		} else {
			return nil, errors.Errorf("syntax error: expected PIPELINE after CREATE; found %s of type %s", p.curToken.Literal, p.curToken.Type)
		}
	} else {
		return nil, errors.Errorf("syntax error: expected CREATE; found %s of type %s", p.curToken.Literal, p.curToken.Type)
	}
	return pipelineInfo, nil
}

func (p *Parser) parseCreateProcessStatement() (*pipeline.ProcessDefinition, error) {
	process := &pipeline.ProcessDefinition{}
	if p.curToken.Type == ADD {
		p.nextToken()
		if p.curToken.Type == SINK {
			process.Sink = true
			p.nextToken()
		}
		if p.curToken.Type == STRING {
			process.ProcessName = p.curToken.Literal
			p.nextToken()
			if p.curToken.Type == OF {
				p.nextToken()
				if p.curToken.Type == STRING {
					process.TypeName = p.curToken.Literal
					p.nextToken()
					for p.curToken.Type != SEMICOLON {
						switch p.curToken.Type {
						case OUTPUTS:
							outs, err := p.parseOutputs()
							if err != nil {
								return nil, err
							}
							process.Outputs = outs
						case INPUTS:
							ins, err := p.parseInputs()
							if err != nil {
								return nil, err
							}
							process.Inputs = ins
						case SET:
							sets, err := p.parseSetState()
							if err != nil {
								return nil, err
							}
							process.State = sets
						default:

						}
						p.nextToken()
					}
					if process.Inputs == nil && process.Outputs == nil && process.State == nil {
						return nil, errors.Errorf("definition error: no INPUTS, OUTPUTS, or STATE defined for %s", process.ProcessName)
					}
				} else {
					return nil, errors.Errorf("syntax error: expected string literal of <process type name> after OF; found %s, %s", p.curToken.Literal, p.curToken.Type)
				}
			} else {
				return nil, errors.Errorf("syntax error: expected OF after string literal of <process name>; found %s of type %s", p.curToken.Literal, p.curToken.Type)
			}
		} else {
			return nil, errors.Errorf("syntax error: expected string literal of <process name> after ADD [SINK]; found %s of type %s", p.curToken.Literal, p.curToken.Type)
		}
	} else {
		return nil, errors.Errorf("syntax error: expected ADD for process create statement; found %s of type %s", p.curToken.Literal, p.curToken.Type)
	}
	return process, nil
}

func (p *Parser) parseCreateConnectionStatement() (*pipeline.ConnectionDefinition, error) {
	connection := &pipeline.ConnectionDefinition{}
	if p.curToken.Type == CONNECT {
		p.nextToken()
		if p.curToken.Type == STRING {
			connection.OriginProcessName = p.curToken.Literal
			p.nextToken()
			if p.curToken.Type == COLON {
				p.nextToken()
				if p.curToken.Type == STRING {
					connection.OriginPortName = p.curToken.Literal
					p.nextToken()
					if p.curToken.Type == TO {
						p.nextToken()
						if p.curToken.Type == STRING {
							connection.DestinationProcessName = p.curToken.Literal
							p.nextToken()
							if p.curToken.Type == COLON {
								p.nextToken()
								if p.curToken.Type == STRING {
									connection.DestinationPortName = p.curToken.Literal
									p.nextToken()
								} else {
									return nil, errors.Errorf("syntax error: expected string literal of <destination port name> after <destination process name>:; found %s, %s", p.curToken.Literal, p.curToken.Type)
								}
							} else {
								return nil, errors.Errorf("syntax error: expected colon after <destination process name>; found %s, %s", p.curToken.Literal, p.curToken.Type)
							}
						} else {
							return nil, errors.Errorf("syntax error: expected <destination process name> after TO; found %s, %s", p.curToken.Literal, p.curToken.Type)
						}
					} else {
						return nil, errors.Errorf("syntax error: expected TO after <origin process name>:<origin port name>; found %s, %s", p.curToken.Literal, p.curToken.Type)
					}
				} else {
					return nil, errors.Errorf("syntax error: expected string literal of <origin port name> after <origin process name>:; found %s, %s", p.curToken.Literal, p.curToken.Type)
				}
			} else {
				return nil, errors.Errorf("syntax error: expected colon after <origin process name>; found %s, %s", p.curToken.Literal, p.curToken.Type)
			}
		} else {
			return nil, errors.Errorf("syntax error: expected <origin process name> after CONNECT; found %s, %s", p.curToken.Literal, p.curToken.Type)
		}
	}
	return connection, nil
}

func (p *Parser) parseOutputs() ([]string, error) {
	var outputs []string
	p.nextToken()
	if p.curToken.Type == ASSIGN {
		p.nextToken()
		if p.curToken.Type == LPAREN {
			p.nextToken()
			for p.curToken.Type != RPAREN {
				if p.curToken.Type == COMMA {
					p.nextToken()
				}
				if p.curToken.Type == STRING {
					outputs = append(outputs, p.curToken.Literal)
				} else {
					return nil, errors.Errorf("syntax error: expected string literal of <output name> after OUTPUTS = (; found %s, %s", p.curToken.Literal, p.curToken.Type)
				}
				p.nextToken()
			}
		} else if p.curToken.Type == STRING {
			outputs = append(outputs, p.curToken.Literal)
		} else {
			return nil, errors.Errorf("syntax error: expected string literal of <output name> or RPAREN after OUTPUTS = ; found %s, %s", p.curToken.Literal, p.curToken.Type)
		}
	}
	return outputs, nil
}

func (p *Parser) parseInputs() ([]string, error) {
	var inputs []string
	p.nextToken()
	if p.curToken.Type == ASSIGN {
		p.nextToken()
		if p.curToken.Type == LPAREN {
			p.nextToken()
			for p.curToken.Type != RPAREN {
				if p.curToken.Type == COMMA {
					p.nextToken()
				}
				if p.curToken.Type == STRING {
					inputs = append(inputs, p.curToken.Literal)
				} else {
					return nil, errors.Errorf("syntax error: expected string literal of <input name> after INPUTS = (; found %s of type %s", p.curToken.Literal, p.curToken.Type)
				}
				p.nextToken()
			}
		} else if p.curToken.Type == STRING {
			inputs = append(inputs, p.curToken.Literal)
		} else {
			return nil, errors.Errorf("syntax error: expected string literal of <inout name> or RPAREN after INPUTS = ; found %s of type %s", p.curToken.Literal, p.curToken.Type)
		}
	}
	return inputs, nil
}

func (p *Parser) parseSetState() (map[string]string, error) {
	state := make(map[string]string)
	for p.curToken.Type != SEMICOLON && p.peekToken.Type != SEMICOLON {
		p.nextToken()
		if p.curToken.Type == COMMA {
			p.nextToken()
		}
		if p.curToken.Type == STRING {
			key := p.curToken.Literal
			p.nextToken()
			if p.curToken.Type == ASSIGN {
				p.nextToken()
				if p.curToken.Type == STRING {
					state[key] = p.curToken.Literal
				} else {
					return nil, errors.Errorf("syntax error: expected string literal of <state value> after =; found %s of type %s", p.curToken.Literal, p.curToken.Type)
				}
			} else {
				return nil, errors.Errorf("syntax error: expected = after string literal; found %s of type %s", p.curToken.Literal, p.curToken.Type)
			}
		} else {
			return nil, errors.Errorf("syntax error: expected string literal of <state key> after SET; found %s of type %s", p.curToken.Literal, p.curToken.Type)
		}
	}
	return state, nil
}
