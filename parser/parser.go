package parser

import (
	"errors"
	"fmt"

	"github.com/lycheng/monkey-go/ast"
	"github.com/lycheng/monkey-go/lexer"
	"github.com/lycheng/monkey-go/token"
)

type (
	prefixParseFn func() (ast.Expression, error)
	// accept left side operator
	infixParseFn func(ast.Expression) (ast.Expression, error)
)

// Parser uses Lexer to parse tokens
type Parser struct {
	l *lexer.Lexer

	errors []string

	currToken token.Token
	peekToken token.Token

	prefixParseFns map[token.Type]prefixParseFn
	infixParseFns  map[token.Type]infixParseFn
}

// New return new Parser with Lexer instance
func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: make([]string, 0)}

	p.registerParseFuncs()

	// Read two times to set curr and peek token
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.currToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) parseStatement() (ast.Statement, error) {
	switch p.currToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() (*ast.LetStatement, error) {
	stmt := &ast.LetStatement{Token: p.currToken}
	if !p.expectPeek(token.IDENT) {
		return nil, errors.New("let statement has no indent token")
	}
	stmt.Name = &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
	if !p.expectPeek(token.ASSIGN) {
		return nil, errors.New("let statement has no assign token")
	}

	p.nextToken()
	val, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, errors.New("can not parse assign expression")
	}
	stmt.Value = val

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt, nil
}

func (p *Parser) parseReturnStatement() (*ast.ReturnStatement, error) {
	stmt := &ast.ReturnStatement{Token: p.currToken}
	p.nextToken()

	expr, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, errors.New("can not parse return expression")
	}
	stmt.ReturnValue = expr

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt, nil
}

func (p *Parser) peekError(t token.Type) {
	msg := fmt.Sprintf(
		"expect next token to be %s, but got %s",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) currTokenIs(t token.Type) bool {
	return p.currToken.Type == t

}
func (p *Parser) peekTokenIs(t token.Type) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.Type) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

// ParseProgram returns AST Program
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = make([]ast.Statement, 0)

	for p.currToken.Type != token.EOF {
		if stmt, err := p.parseStatement(); err == nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

// Errors returns the errors durning parsing
func (p *Parser) Errors() []string {
	return p.errors
}
