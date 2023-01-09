package ast

import "github.com/lycheng/monkey-go/token"

// Node for AST node interface
type Node interface {
	TokenLiteral() string
}

// Statement for AST statement interface
type Statement interface {
	Node
	statementNode()
}

// Expression for AST expression intearface
type Expression interface {
	Node
	expressionNode()
}

// Program contains statements
// It's a rot node of AST
type Program struct {
	Statements []Statement
}

// TokenLiteral returns the first statement's literal
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

// Identifier for the ident in AST
type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (i *Identifier) expressionNode() {}

// TokenLiteral returns identifier's literal value
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

// LetStatement for the statement let x = ...
type LetStatement struct {
	// Token for Let token
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode() {}

// TokenLiteral returns the let token literal value
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
