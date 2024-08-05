package lexer

import "rowanlovejoy/monkey/token"

type Lexer struct {
	input        string
	position     int  // Position of last read character
	readPosition int  // Position of next character to read
	ch           byte // Current char under examination (pointed to by position)
}

// Create and initialise a new Lexer instance with first input char already read
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// Return the token corresponding to the current char and then advance the lexer
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	switch l.ch {
	case '=':
		tok = token.New(token.ASSIGN, l.ch)
	case ';':
		tok = token.New(token.SEMICOLON, l.ch)
	case ',':
		tok = token.New(token.COMMA, l.ch)
	case '+':
		tok = token.New(token.PLUS, l.ch)
	case '(':
		tok = token.New(token.LPAREN, l.ch)
	case ')':
		tok = token.New(token.RPAREN, l.ch)
	case '{':
		tok = token.New(token.LBRACE, l.ch)
	case '}':
		tok = token.New(token.RBRACE, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	}

	l.readChar()
	return tok
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // ASCII code for the "NUL" char, represents EOF
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}
