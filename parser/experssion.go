package parser

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/lycheng/monkey-go/ast"
	"github.com/lycheng/monkey-go/token"
)

func (p *Parser) registerParseFuncs() {
	p.prefixParseFns = make(map[token.Type]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)

	p.infixParseFns = make(map[token.Type]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOTEQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
}

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

func (p *Parser) parseIdentifier() (ast.Expression, error) {
	return &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}, nil
}

func (p *Parser) parseIntegerLiteral() (ast.Expression, error) {
	il := &ast.IntegerLiteral{Token: p.currToken}
	val, err := strconv.ParseInt(p.currToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.currToken.Literal)
		p.errors = append(p.errors, msg)
		return nil, errors.New(msg)
	}

	il.Value = val
	return il, nil
}

func (p *Parser) parseInfixExpression(left ast.Expression) (ast.Expression, error) {
	expression := &ast.InfixExpression{
		Token:    p.currToken,
		Operator: p.currToken.Literal,
		Left:     left,
	}
	precedence := p.currPrecedence()
	p.nextToken()
	exp, err := p.parseExpression(precedence)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as infix expression", p.currToken.Literal)
		p.errors = append(p.errors, msg)
		return nil, errors.New(msg)
	}
	expression.Right = exp
	return expression, nil
}
