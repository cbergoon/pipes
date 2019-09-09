package dl

type Lexer struct {
	input        string
	position     int
	readPosition int
	linePosition int
	ch           byte
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		tok = newToken(ASSIGN, l.ch)
	case '/':
		tok = newToken(SLASH, l.ch)
	case '*':
		tok = newToken(ASTERISK, l.ch)
	case ';':
		tok = newToken(SEMICOLON, l.ch)
	case '(':
		tok = newToken(LPAREN, l.ch)
	case ')':
		tok = newToken(RPAREN, l.ch)
	case ',':
		tok = newToken(COMMA, l.ch)
	case '{':
		tok = newToken(LBRACE, l.ch)
	case '}':
		tok = newToken(RBRACE, l.ch)
	case '[':
		tok = newToken(LBRAKET, l.ch)
	case ']':
		tok = newToken(RBRAKET, l.ch)
	case ':':
		tok = newToken(COLON, l.ch)
	case '"':
		tok.Type = STRING
		tok.Literal = l.readString('"')
	case '`':
		tok.Type = STRING
		tok.Literal = l.readString('`')
	case '\'':
		tok.Type = STRING
		tok.Literal = l.readString('\'')
	case 0:
		tok.Literal = ""
		tok.Type = EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}

	return l.input[l.readPosition]
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}

	return l.spliceInput(position, l.position)
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.spliceInput(position, l.position)
}

func (l *Lexer) skipWhitespace() {
	for isWhitespace(l.ch) {
		l.readChar()
	}
}

func (l *Lexer) readString(delim byte) string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == delim {
			break
		}
	}

	return l.spliceInput(position, l.position)
}

func (l *Lexer) spliceInput(start, end int) string {
	return spliceString(l.input, start, end)
}

func newToken(tokenType TokenType, ch byte) Token {
	return Token{
		Type:    tokenType,
		Literal: string(ch),
	}
}

func spliceString(in string, start, end int) string {
	return in[start:end]
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}
