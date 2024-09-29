package parser

import (
	"rowanlovejoy/monkey/ast"
	"rowanlovejoy/monkey/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
		let x = 5;
		let y = 10;
		let foobar = 838383;
	`
	expectedNumStatements := 3

	parser := New(lexer.New(input))

	program := parser.ParseProgram()
	checkParserErrors(t, parser)

	if numStatements := len(program.Statements); numStatements != expectedNumStatements {
		t.Fatalf("Unexpected statement count. Expected %d statements. Got %d", expectedNumStatements, numStatements)
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, test := range tests {
		if !testLetStatement(t, program.Statements[i], test.expectedIdentifier) {
			return
		}
	}
}

func testLetStatement(t *testing.T, statement ast.Statement, identifier string) bool {
	if statement.TokenLiteral() != "let" {
		t.Errorf("Unexpected token literal. Expected \"let\". Got %q", statement.TokenLiteral())
		return false
	}

	letStatement, ok := statement.(*ast.LetStatement)
	if !ok {
		t.Errorf("Unexpected statement type. Expected *ast.LetStatement. Got %T", statement)
		return false
	}

	if name := letStatement.Name.Value; name != identifier {
		t.Errorf("Unexpected let statement name. Expected %q. Got %q", identifier, name)
		return false
	}

	if literal := letStatement.Name.TokenLiteral(); literal != identifier {
		t.Errorf("Unexpected let statement token literal. Expected %q. Got %q", identifier, literal)
		return false
	}

	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()

	if len(errors) == 0 {
		return
	}

	t.Errorf("Parser has %d error(s)", len(errors))

	for _, message := range errors {
		t.Errorf("Parser error: %q", message)
	}

	t.FailNow()
}
