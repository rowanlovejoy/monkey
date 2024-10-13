package parser

import (
	"fmt"
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

	checkParserErrors(t, parser)
	checkStatementCount(t, program, 3)

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

func TestReturnStatements(t *testing.T) {
	input := `
		return 5;
		return 10;
		return 993322;
	`
	parser := New(lexer.New(input))

	program := parser.ParseProgram()

	checkParserErrors(t, parser)
	checkStatementCount(t, program, 3)

	for _, statement := range program.Statements {
		returnStatement, ok := statement.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("Unexpected statement type. Expected *ast.ReturnStatement; got %T", returnStatement)
			continue
		}
		if literal := returnStatement.TokenLiteral(); literal != "return" {
			t.Errorf("Unexpected return statement token literal. Expected \"return\". Got %q", literal)
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := `
		foobar;
	`

	parser := New(lexer.New(input))
	program := parser.ParseProgram()

	checkParserErrors(t, parser)
	checkStatementCount(t, program, 1)

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Unexpected statement type. Expected *ast.ExpressionStatement; got %T", program.Statements[0])
	}

	identifier, ok := statement.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("Unexpected expression type. Expected *ast.Identifier; got %T", statement.Expression)
	}

	expectedIdentifier := "foobar"
	if value := identifier.Value; value != expectedIdentifier {
		t.Errorf("Unexpected identifier value. Expected %q; got %q", expectedIdentifier, value)
	}

	expectedTokenLiteral := "foobar"
	if literal := identifier.TokenLiteral(); literal != expectedTokenLiteral {
		t.Errorf("Unexpected token literal. Expected %q; got %q", expectedTokenLiteral, literal)
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := `
		5;
	`
	parser := New(lexer.New(input))
	program := parser.ParseProgram()

	checkParserErrors(t, parser)
	checkStatementCount(t, program, 1)

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Unexpected statement type. Expected *ast.ExpressionStatement; got %T", program.Statements[0])
	}

	integerLiteral, ok := statement.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("Unexpected expression type. Expected *ast.IntegerLiteral; got %T", statement.Expression)
	}

	expectedValue := 5
	if value := integerLiteral.Value; value != int64(expectedValue) {
		t.Errorf("Unexpected literal value. Expected %d; got %d", expectedValue, value)
	}

	expectedTokenLiteral := "5"
	if literal := integerLiteral.TokenLiteral(); literal != expectedTokenLiteral {
		t.Errorf("Unexpected token literal. Expected %q; got %q", expectedTokenLiteral, literal)
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
	}

	for _, test := range prefixTests {
		parser := New(lexer.New(test.input))
		program := parser.ParseProgram()

		checkParserErrors(t, parser)
		checkStatementCount(t, program, 1)

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Unexpected statement type. Expected *ast.ExpressionStatement; got %T", program.Statements[0])
		}

		prefixExpression, ok := statement.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("Unexpected expression type. Expected *ast.PrefixExpression; got %T", statement.Expression)
		}

		if operator := prefixExpression.Operator; operator != test.operator {
			t.Fatalf("Unexpected operator. Expected %q; got %q", test.operator, operator)
		}

		if !testIntegerLiteral(t, prefixExpression.Right, test.integerValue) {
			return
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"5 + 5", 5, "+", 5},
		{"5 - 5", 5, "-", 5},
		{"5 * 5", 5, "*", 5},
		{"5 / 5", 5, "/", 5},
		{"5 > 5", 5, ">", 5},
		{"5 < 5", 5, "<", 5},
		{"5 == 5", 5, "==", 5},
		{"5 != 5", 5, "!=", 5},
	}

	for _, test := range infixTests {
		parser := New(lexer.New(test.input))
		program := parser.ParseProgram()

		checkParserErrors(t, parser)
		checkStatementCount(t, program, 1)

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Unexpected statement type. Expected *ast.ExpressionStatement; got %T", program.Statements[0])
		}

		infixExpression, ok := statement.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("Unexpected expression type. Expected *ast.InfixExpression; got %T", statement.Expression)
		}

		if !testIntegerLiteral(t, infixExpression.Left, test.leftValue) {
			return
		}

		if infixExpression.Operator != test.operator {
			t.Fatalf("Unexpected infix expression operator. Expected %q; got %q", test.operator, infixExpression.Operator)
		}

		if !testIntegerLiteral(t, infixExpression.Right, test.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e -f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
	}

	for _, test := range tests {
		parser := New(lexer.New(test.input))
		program := parser.ParseProgram()
		checkParserErrors(t, parser)

		if actual := program.String(); actual != test.expected {
			t.Errorf("Unexpected string output. Expected %q; got %q", test.expected, actual)
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

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integerLiteral, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("Unexpected expression type. Expected *ast.IntegerLiteral; got %T", il)
		return false
	}

	if literalValue := integerLiteral.Value; literalValue != value {
		t.Errorf("Unexpected literal value. Expected %d; got %d", value, literalValue)
		return false
	}

	if tokenLiteral := integerLiteral.TokenLiteral(); tokenLiteral != fmt.Sprintf("%d", value) {
		t.Errorf("Unexpected token literal. Expected %d; got %q", value, tokenLiteral)
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

func checkStatementCount(t *testing.T, program *ast.Program, expectedCount int) {
	if numStatements := len(program.Statements); numStatements != expectedCount {
		t.Fatalf("Unexpected statement count. Expected %d statement(s); got %d", expectedCount, numStatements)
	}
}
