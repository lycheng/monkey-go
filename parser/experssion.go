package parser

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/lycheng/monkey-go/ast"
	"github.com/lycheng/monkey-go/token"
)

// precedence
const (
	_ int = iota
	LOWEST
	EQUALS
	LESSGRAEATER
	SUM
	PRODUCT
	PREFIX
	CALL
	INDEX
)

var precedences = map[token.Type]int{
	token.EQ:       EQUALS,
	token.NOTEQ:    EQUALS,
	token.LT:       LESSGRAEATER,
	token.GT:       LESSGRAEATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
	token.LBRACKET: INDEX,
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) currPrecedence() int {
	if p, ok := precedences[p.currToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) registerParseFuncs() {
	p.prefixParseFns = make(map[token.Type]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)

	p.infixParseFns = make(map[token.Type]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOTEQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)
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
	pfn, ok := p.prefixParseFns[p.currToken.Type]
	if !ok {
		msg := "parse func for " + string(p.currToken.Type) + " not found"
		return nil, errors.New(msg)
	}
	exp, err := pfn()
	if err != nil {
		return nil, err
	}

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		ifn, ok := p.infixParseFns[p.peekToken.Type]
		if !ok {
			return exp, nil
		}

		p.nextToken()
		iExp, err := ifn(exp)
		if err != nil {
			return nil, err
		}
		exp = iExp
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

func (p *Parser) parseBoolean() (ast.Expression, error) {
	return &ast.Boolean{Token: p.currToken, Value: p.currTokenIs(token.TRUE)}, nil
}

func (p *Parser) parseStringLiteral() (ast.Expression, error) {
	return &ast.StringLiteral{Token: p.currToken, Value: p.currToken.Literal}, nil
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

func (p *Parser) parseCallExpression(function ast.Expression) (ast.Expression, error) {
	exp := &ast.CallExpression{Token: p.currToken, Function: function}
	args, err := p.parseExpressionList(token.RPAREN)
	if err != nil {
		return nil, err
	}
	exp.Arguments = args
	return exp, nil
}

func (p *Parser) parseGroupedExpression() (ast.Expression, error) {
	p.nextToken()
	exp, err := p.parseExpression(LOWEST)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as expression", p.currToken.Literal)
		p.errors = append(p.errors, msg)
		return nil, errors.New(msg)
	}

	if !p.expectPeek(token.RPAREN) {
		msg := "could not match the right parenthesis"
		p.errors = append(p.errors, msg)
		return nil, errors.New(msg)
	}
	return exp, nil
}

func (p *Parser) parseBlockStatement() (*ast.BlockStatement, error) {
	block := &ast.BlockStatement{Token: p.currToken}
	block.Statements = []ast.Statement{}
	p.nextToken()
	for !p.currTokenIs(token.RBRACE) && !p.currTokenIs(token.EOF) {
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		block.Statements = append(block.Statements, stmt)
		p.nextToken()
	}
	return block, nil
}

func (p *Parser) parseIfExpression() (ast.Expression, error) {
	expr := &ast.IfExpression{Token: p.currToken}
	if !p.expectPeek(token.LPAREN) {
		return nil, errors.New("token ( not found for if expression")
	}

	p.nextToken()
	ce, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}
	expr.Condition = ce

	if !p.expectPeek(token.RPAREN) {
		return nil, errors.New("token ) not found for if expression")
	}

	if !p.expectPeek(token.LBRACE) {
		return nil, errors.New("token { not found for if expression")
	}
	bs, err := p.parseBlockStatement()
	if err != nil {
		return nil, err
	}
	expr.Consequence = bs

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		if !p.expectPeek(token.LBRACE) {
			return nil, errors.New("token { for found for else block")
		}

		al, err := p.parseBlockStatement()
		if err != nil {
			return nil, err
		}
		expr.Alternative = al
	}
	return expr, nil
}

func (p *Parser) parseFunctionParameters() ([]*ast.Identifier, error) {
	identifiers := []*ast.Identifier{}
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers, nil
	}
	p.nextToken()
	ident := &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
	identifiers = append(identifiers, ident)
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		ident := &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
		identifiers = append(identifiers, ident)
	}
	if !p.expectPeek(token.RPAREN) {
		return nil, errors.New("no ) token for function parameters")
	}
	return identifiers, nil
}

func (p *Parser) parseFunctionLiteral() (ast.Expression, error) {
	fn := &ast.FunctionLiteral{Token: p.currToken}
	if !p.expectPeek(token.LPAREN) {
		return nil, errors.New("no ( token for function definition")
	}
	params, err := p.parseFunctionParameters()
	if err != nil {
		return nil, err
	}
	fn.Parameters = params
	if !p.expectPeek(token.LBRACE) {
		return nil, errors.New("no ) token for function definition")
	}
	body, err := p.parseBlockStatement()
	if err != nil {
		return nil, err
	}
	fn.Body = body
	return fn, nil
}

func (p *Parser) parseArrayLiteral() (ast.Expression, error) {
	array := &ast.ArrayLiteral{Token: p.currToken}
	elements, err := p.parseExpressionList(token.RBRACKET)
	if err != nil {
		return nil, err
	}
	array.Elements = elements
	return array, nil
}

func (p *Parser) parseExpressionList(end token.Type) ([]ast.Expression, error) {
	list := []ast.Expression{}
	if p.peekTokenIs(end) {
		p.nextToken()
		return list, nil
	}
	p.nextToken()
	expr, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}
	list = append(list, expr)
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		expr, err = p.parseExpression(LOWEST)
		if err != nil {
			return nil, err
		}
		list = append(list, expr)
	}
	if !p.expectPeek(end) {
		return nil, fmt.Errorf("Expect to get %s but get %s", end, p.peekToken.Literal)
	}
	return list, nil
}

func (p *Parser) parseIndexExpression(left ast.Expression) (ast.Expression, error) {
	exp := &ast.IndexExpression{Token: p.currToken, Left: left}
	p.nextToken()
	expr, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}
	exp.Index = expr
	if !p.expectPeek(token.RBRACKET) {
		return nil, fmt.Errorf("Expect to get ] but get %s", p.peekToken.Literal)
	}
	return exp, nil
}
