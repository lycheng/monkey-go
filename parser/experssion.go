package parser

import (
	"errors"

	"github.com/lycheng/monkey-go/ast"
	"github.com/lycheng/monkey-go/token"
)

func (p *Parser) registerPrefix(tokenType token.Type, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.Type, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

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
	return prefix()
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
