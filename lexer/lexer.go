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

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if literal, ok := l.makeTwoCharLiteral("=="); ok {
			tok = token.Token{Type: token.EQ, Literal: literal}
		} else {
			tok = token.New(token.ASSIGN, l.ch)
		}
	case '+':
		tok = token.New(token.PLUS, l.ch)
	case '-':
		tok = token.New(token.MINUS, l.ch)
	case '!':
		if literal, ok := l.makeTwoCharLiteral("!="); ok {
			tok = token.Token{Type: token.NOTEQ, Literal: literal}
		} else {
			tok = token.New(token.BANG, l.ch)
		}
	case '/':
		tok = token.New(token.SLASH, l.ch)
	case '*':
		tok = token.New(token.ASTERISK, l.ch)
	case '<':
		tok = token.New(token.LT, l.ch)
	case '>':
		tok = token.New(token.GT, l.ch)
	case ',':
		tok = token.New(token.COMMA, l.ch)
	case ';':
		tok = token.New(token.SEMICOLON, l.ch)
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
	default:
		// Letter and digit branches exit early due to having already advanced the lexer
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = token.New(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

// Attempt to construct the specified two char literal from the current and next char, advancing lexer if successful
func (l *Lexer) makeTwoCharLiteral(expected string) (string, bool) {
	nextCh := l.peekChar()
	literal := string(l.ch) + string(nextCh)
	if literal == expected {
		l.readChar()
		return expected, true
	} else {
		return "", false
	}
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
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

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// Return the next char to be read without advancing the lexer
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
