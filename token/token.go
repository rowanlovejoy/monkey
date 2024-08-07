package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	// Special
	ILLEGAL = "ILLEGAL" // Unsupported token
	EOF     = "EOF"     // End of file

	// Identifiers and literals
	IDENT = "IDENT" // E.g., add, foobar, x, y
	INT   = "INT"   // E.g., 3, 5

	// Operators
	ASSIGN   = "ASSIGN"   // =
	PLUS     = "PLUS"     // +
	MINUS    = "MINUS"    // -
	BANG     = "BANG"     // !
	ASTERISK = "ASTERISK" // *
	SLASH    = "SLASH"    // /
	LT       = "LT"       // AKA less than, <
	GT       = "GT"       // AKA greater than, >

	// Delimiters
	COMMA     = "COMMA"     // ,
	SEMICOLON = "SEMICOLON" // ;
	LPAREN    = "LPAREN"    // (
	RPAREN    = "RPAREN"    // )
	LBRACE    = "LBRACE"    // {
	RBRACE    = "RBRACE"    // }

	// Keywords
	FUNCTION = "FUNCTION" // fn
	LET      = "LET"      // let
)

var keywords = map[string]TokenType{
	"fn":  FUNCTION,
	"let": LET,
}

func New(tokenType TokenType, ch byte) Token {
	return Token{Type: tokenType, Literal: string(ch)}
}

// Get the token type matching the specified identifier
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
