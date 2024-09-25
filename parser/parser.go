package parser

import (
	"rowanlovejoy/monkey/ast"
	"rowanlovejoy/monkey/lexer"
	"rowanlovejoy/monkey/token"
)

type Parser struct {
	l *lexer.Lexer
	// Analogous to Lexer's position and readPosition but store tokens instead of chars
	currToken token.Token // Current token under examination
	peekToken token.Token // Next token in the sequence, can give context to current token when parsing
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

	// Read two tokens so that currToken and peekToken are both initialised
	p.nextToken() // Initialises peekToken
	p.nextToken() // Initialises currToken with value of peekToken and updates peekToken

	return p
}

// Advances the parser through the token sequence
func (p *Parser) nextToken() {
	p.currToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{
		Statements: []ast.Statement{},
	}

	for p.currToken.Type != token.EOF {
		if statement := p.parseStatement(); statement != nil {
			program.Statements = append(program.Statements, statement)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currToken.Type {
	case token.LET:
		return p.parseLetStatement()
	default:
		return nil
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	statement := &ast.LetStatement{
		Token: p.currToken,
	}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	statement.Name = &ast.Identifier{
		Token: p.currToken,
		Value: p.currToken.Literal,
	}

	// TODO: Skip over expressions for now
	for !p.currTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return statement
}

// Compare type of current token to expected
func (p *Parser) currTokenIs(t token.TokenType) bool {
	return p.currToken.Type == t
}

// Compare type of next token to expected
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// Asserts the next token and advances the parser if true
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		return false
	}
}
