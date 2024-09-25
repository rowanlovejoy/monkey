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

	parser := New(lexer.New(input))

	program := parser.ParseProgram()
	if program == nil {
		t.Fatal("Unexpected value. Expected ast.Program. Got nil")
	}

	expectedNumStatements := 3
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
		t.Errorf("Unexpected token literal. Expected 'let'. Got %s", statement.TokenLiteral())
		return false
	}

	letStatement, ok := statement.(*ast.LetStatement)
	if !ok {
		t.Errorf("Unexpected statement type. Expected *ast.LetStatement. Got %T", statement)
		return false
	}

	if name := letStatement.Name.Value; name != identifier {
		t.Errorf("Unexpected let statement name. Expected %s. Got %s", identifier, name)
		return false
	}

	if literal := letStatement.Name.TokenLiteral(); literal != identifier {
		t.Errorf("Unexpected let statement token literal. Expected %s. Got %s", identifier, literal)
		return false
	}

	return true
}
