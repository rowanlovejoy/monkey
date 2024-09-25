package ast

import "rowanlovejoy/monkey/token"

type Node interface {
	TokenLiteral() string
}

// Represents a unit of code that doesn't produce value, e.g., 'let x = 5';
type Statement interface {
	Node
	// Dummy method to help with type checking
	statementNode()
}

// Represents a unit of code that produces a value, e.g., '5', 'add(1, 2)', 'fn (...) {...}',
type Expression interface {
	Node
	// Dummy method to help with type checking
	expressionNode()
}

// The root node of every AST produced by the parser
type Program struct {
	Statements []Statement
}

// Print the token literal associated with the program's first statement
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

type LetStatement struct {
	Token token.Token // token.LET
	Name  *Identifier // Identifier being bound to
	Value Expression  // Expression returning the value to be bound
}

func (ls *LetStatement) statementNode() {} // Satisfies Statement interface
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
} // Satisfies Node interface

// A name to which a value has been bound
type Identifier struct {
	Token token.Token // token.IDENT
	Value string      // The identifier literal
}

func (i *Identifier) expressionNode() {} // Satisfies Expression interface. Identifiers are expressions because in some cases they *can* produce values, e.g., when binding one variable to another, i.e., let second_identifier = first_identifier;
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
} // Satisfies Node interface
