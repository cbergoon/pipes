package dl

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENT  = "IDENT"
	INT    = "INT"
	STRING = "STRING"

	ASSIGN   = "="
	ASTERISK = "*"
	SLASH    = "/"

	COMMA     = ","
	SEMICOLON = ";"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
	LBRAKET   = "["
	RBRAKET   = "]"
	COLON     = ":"

	CREATE   = "CREATE"
	PIPELINE = "PIPELINE"
	ADD      = "ADD"
	SINK     = "SINK"
	CONNECT  = "CONNECT"
	INPUTS   = "INPUTS"
	OUTPUTS  = "OUTPUTS"
	SET      = "SET"
	TO       = "TO"
	OF       = "OF"
)

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

var keywords = map[string]TokenType{
	"CREATE":   CREATE,
	"PIPELINE": PIPELINE,
	"ADD":      ADD,
	"SINK":     SINK,
	"CONNECT":  CONNECT,
	"INPUTS":   INPUTS,
	"OUTPUTS":  OUTPUTS,
	"SET":      SET,
	"TO":       TO,
	"OF":       OF,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}

	return IDENT
}
