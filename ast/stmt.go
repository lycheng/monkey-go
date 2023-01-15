package ast

import (
	"bytes"

	"github.com/lycheng/monkey-go/token"
)

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

// String returns let statement string value
func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

// ReturnStatement for statement return ...
type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}

// TokenLiteral returns the let token literal value
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }

// String return expression string value
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

// ExpressionStatement is a wrapper of statement
type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}

// TokenLiteral returns the first token of the expression
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }

// String returns expression string value
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}
