package ast

import (
	"bytes"
	"rowanlovejoy/monkey/token"
)

// String returned when calling TokenLiteral on a nil receiver
const NIL_TOKEN_LITERAL = "<nil>"

type Node interface {
	TokenLiteral() string
	String() string
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
} // Satisfies Node interface

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
} // Satisfies Node interface

type LetStatement struct {
	Token token.Token // token.LET
	Name  *Identifier // Identifier being bound to
	Value Expression  // Expression returning the value to be bound
}

func (ls *LetStatement) statementNode() {} // Satisfies Statement interface
func (ls *LetStatement) TokenLiteral() string {
	if ls == nil {
		return NIL_TOKEN_LITERAL
	}
	return ls.Token.Literal
} // Satisfies Node interface

func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")

	return out.String()
} // Satisfies Node interface

type ReturnStatement struct {
	Token       token.Token // token.RETURN
	ReturnValue Expression  // Expression returning the value to return
}

func (rs *ReturnStatement) statementNode() {} // Satisfies Statement interface
func (rs *ReturnStatement) TokenLiteral() string {
	if rs == nil {
		return NIL_TOKEN_LITERAL
	}
	return rs.Token.Literal
} // Satisfies Node interface

func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
} // Satisfies Node interface

type ExpressionStatement struct {
	Token      token.Token // First token in the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {} // Satisfies Statement interface
func (es *ExpressionStatement) TokenLiteral() string {
	if es == nil {
		return NIL_TOKEN_LITERAL
	}
	return es.Token.Literal
} // Satisfies Node interface

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// A name to which a value has been bound
type Identifier struct {
	Token token.Token // token.IDENT
	Value string      // The identifier literal
}

func (i *Identifier) expressionNode() {} // Satisfies Expression interface. Identifiers are expressions because in some cases they *can* produce values, e.g., when binding one variable to another, i.e., let second_identifier = first_identifier;
func (i *Identifier) TokenLiteral() string {
	if i == nil {
		return NIL_TOKEN_LITERAL
	}
	return i.Token.Literal
} // Satisfies Node interface

func (i *Identifier) String() string {
	return i.Value
} // Satisfies Node interface
