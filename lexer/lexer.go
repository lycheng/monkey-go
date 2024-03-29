package lexer

import (
	"github.com/lycheng/monkey-go/token"
)

// Lexer for monkey
type Lexer struct {
	input   string
	currPos int // current position of input
	nextPos int // next position of input
	ch      byte
}

// New return new Lexer object with input string
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// NextToken returns next token from the input
func (l *Lexer) NextToken() token.Token {
	l.skipWhitespace()
	tk := token.Token{Literal: string(l.ch)}
	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			tk.Type = token.EQ
			tk.Literal = "=="
		} else {
			tk.Type = token.ASSIGN
		}
	case ';':
		tk.Type = token.SEMICOLON
	case ':':
		tk.Type = token.COLON
	case '(':
		tk.Type = token.LPAREN
	case ')':
		tk.Type = token.RPAREN
	case '{':
		tk.Type = token.LBRACE
	case '}':
		tk.Type = token.RBRACE
	case '[':
		tk.Type = token.LBRACKET
	case ']':
		tk.Type = token.RBRACKET
	case ',':
		tk.Type = token.COMMA
	case '+':
		tk.Type = token.PLUS
	case '-':
		tk.Type = token.MINUS
	case '*':
		tk.Type = token.ASTERISK
	case '/':
		tk.Type = token.SLASH
	case '<':
		tk.Type = token.LT
	case '>':
		tk.Type = token.GT
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			tk.Type = token.NOTEQ
			tk.Literal = "!="
		} else {
			tk.Type = token.BANG
		}
	case '"':
		tk.Type = token.STRING
		tk.Literal = l.readString()
	case 0:
		tk.Literal = ""
		tk.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tk.Literal = l.readIdentifier()
			tk.Type = token.LookupIdent(tk.Literal)
			return tk
		} else if isDigit(l.ch) {
			tk.Type = token.INT
			tk.Literal = l.readNumber()
			return tk
		} else {
			tk.Type = token.ILLEGAL
		}
	}
	l.readChar()
	return tk
}

func (l *Lexer) readString() string {
	i := l.currPos + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[i:l.currPos]
}

func (l *Lexer) readChar() {
	if l.nextPos >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.nextPos]
	}
	l.currPos = l.nextPos
	l.nextPos++
}

func (l *Lexer) peekChar() byte {
	if l.nextPos >= len(l.input) {
		return 0
	}
	return l.input[l.nextPos]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func isLetter(b byte) bool {
	return ('a' <= b && b <= 'z') || ('A' <= b && b <= 'Z') || (b == '_')
}

func isDigit(b byte) bool {
	return '0' <= b && b <= '9'
}

func (l *Lexer) readIdentifier() string {
	pos := l.currPos
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[pos:l.currPos]
}

func (l *Lexer) readNumber() string {
	pos := l.currPos
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[pos:l.currPos]
}
