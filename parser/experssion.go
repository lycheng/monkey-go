package parser

import (
	"errors"

	"github.com/lycheng/monkey-go/ast"
	"github.com/lycheng/monkey-go/token"
)

func (p *Parser) parseExpressionStatement() (*ast.ExpressionStatement, error) {
	stmt := &ast.ExpressionStatement{Token: p.currToken}
	exp, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}
	stmt.Expression = exp
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt, nil
}

func (p *Parser) parseExpression(precedence int) (ast.Expression, error) {
	prefix, ok := p.prefixParseFns[p.currToken.Type]
	if !ok {
		msg := "parse func for " + string(p.currToken.Type) + " not found"
		return nil, errors.New(msg)
	}
	exp, err := prefix()
	if err != nil {
		return nil, err
	}

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix, ok := p.infixParseFns[p.peekToken.Type]
		if !ok {
			return exp, nil
		}

		p.nextToken()
		infixExp, err := infix(exp)
		if err != nil {
			return nil, err
		}
		exp = infixExp
	}
	return exp, nil
}

func (p *Parser) parsePrefixExpression() (ast.Expression, error) {
	e := &ast.PrefixExpression{
		Token:    p.currToken,
		Operator: p.currToken.Literal,
	}
	p.nextToken()
	exp, err := p.parseExpression(PREFIX)
	if err != nil {
		return nil, err
	}
	e.Right = exp
	return e, nil
}
