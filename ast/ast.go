package ast

import "rowanlovejoy/monkey/token"

type Node interface {
	TokenLiteral() string
}

type Statement interface {
	Node
	// Dummy method to help with type checking
	statementNode()
}

type Expression interface {
	Node
	// Dummy method to help with type checking
	expressionNode()
}

type Program struct {
	Statements []Statement
}

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

func (ls *LetStatement) statementNode() {}

func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

type Identifier struct {
	Token token.Token // token.IDENT
	Value string      // The identifier literal
}

func (i *Identifier) expressionNode() {} // Identifiers are expressions because in some cases they *can* produce values, e.g., when binding one variable to another, i.e., let second_identifier = first_identifier;

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}
