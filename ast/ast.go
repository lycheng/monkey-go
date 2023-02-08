package ast

import (
	"bytes"

	"github.com/lycheng/monkey-go/token"
)

// Node for AST node interface
type Node interface {
	TokenLiteral() string
	String() string
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

// String returns all statements' string value
func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// Identifier for the ident in AST
type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (i *Identifier) expressionNode() {}

// TokenLiteral returns identifier's literal value
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

// String returns the identifier
func (i *Identifier) String() string { return i.Value }

// IntegerLiteral for the int value in AST
type IntegerLiteral struct {
	Token token.Token // the token.IDENT token
	Value int64
}

func (il *IntegerLiteral) expressionNode() {}

// TokenLiteral returns identifier's literal value
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }

// String returns the identifier
func (il *IntegerLiteral) String() string { return il.Token.Literal }

// PrefixExpression structure
type PrefixExpression struct {
	Token    token.Token // The prefix token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}

// TokenLiteral return the operator literal
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

// InfixExpression structure
type InfixExpression struct {
	Token    token.Token // The operator token, e.g. +
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode() {}

// TokenLiteral return the left expression literal
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")
	return out.String()
}
