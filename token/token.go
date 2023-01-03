package token

// Token types
const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identitfiers + literals
	IDENT = "IDENT"
	INT   = "INT"

	// Operators
	ASSIGN = "="
	PLUS   = "+"

	// Delimiters
	COMMA     = ","
	SEMICOLON = "."

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
)

// Type for monkey's token type
type Type string

// Token for monkey's token
type Token struct {
	Type    Type
	Literal string
}

var keywords = map[string]Type{
	"fn":  FUNCTION,
	"let": LET,
}

// LookupIdent returns ident's type
func LookupIdent(ident string) Type {
	if tk, ok := keywords[ident]; ok {
		return tk
	}
	return IDENT
}
