package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let foobar = 838383;
`
	program := parseAndTestCommonStep(t, input, 3)

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
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
	program := parseAndTestCommonStep(t, input, 3)

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt is not a *ast.ReturnStatement. Got %T", stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral() is not 'return'. Got %q", returnStmt.TokenLiteral())
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	program := parseAndTestCommonStep(t, input, 1)

	stmt := testAndParseToExpressionStatement(t, program)
	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp is not ast.Identifier, got %T", ident)
	}
	if ident.Value != "foobar" {
		t.Errorf("ident value is not %s. got %s", "foobar", ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("token literal is not %s, got %s", "foobar", ident.TokenLiteral())
	}
}

func testAndParseToExpressionStatement(t *testing.T, program *ast.Program) *ast.ExpressionStatement {
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statement[0] is not an expression statement, got %T", program.Statements[0])
	}
	return stmt
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	program := parseAndTestCommonStep(t, input, 1)

	stmt := testAndParseToExpressionStatement(t, program)

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("expression is not integer literal. Got %T", stmt.Expression)
	}
	if literal.Value != 5 {
		t.Errorf("Integer literal's value is not 5. Got %d", literal.Value)
	}
	if literal.TokenLiteral() != "5" {
		t.Errorf("token literal is not %s. Got %s", "5", literal.TokenLiteral())
	}
}

func TestParsingPrefixExpression(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
	}
	for _, tt := range prefixTests {
		t.Run(fmt.Sprintf("Parsing prefix expression for input %s, operator %s", tt.input, tt.operator), func(t *testing.T) {
			program := parseAndTestCommonStep(t, tt.input, 1)
			stmt := testAndParseToExpressionStatement(t, program)
			exp, ok := stmt.Expression.(*ast.PrefixExpression)
			if !ok {
				t.Fatalf("expression is not prefix expression. Got %T", stmt.Expression)
			}
			if exp.Operator != tt.operator {
				t.Fatalf("expression operator is not %s. Got %s", tt.operator, exp.Operator)
			}
			if !testIntegerLiteral(t, exp.Right, tt.integerValue) {
				return
			}
		})
	}
}

func TestParsingInfixExpression(t *testing.T) {
	infixTests := []struct {
		input    string
		left     int64
		operator string
		right    int64
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

	for _, tt := range infixTests {
		t.Run(fmt.Sprintf("Infix parsing for %s", tt.input), func(t *testing.T) {
			program := parseAndTestCommonStep(t, tt.input, 1)
			stmt := testAndParseToExpressionStatement(t, program)
			exp, ok := stmt.Expression.(*ast.InfixExpression)
			if !ok {
				t.Fatalf("exp is not an infix expression. Got %T", stmt.Expression)
			}
			if !testIntegerLiteral(t, exp.Left, tt.left) {
				return
			}
			if exp.Operator != tt.operator {
				t.Fatalf("expression operator is not %s. Got %s", tt.operator, exp.Operator)
			}
			if !testIntegerLiteral(t, exp.Right, tt.right) {
				return
			}
		})
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input              string
		expected           string
		expectedStatements int
	}{
		{"-a * b", "((-a) * b)", 1},
		{"!-a", "(!(-a))", 1},
		{"a + b + c", "((a + b) + c)", 1},
		{"a + b - c", "((a + b) - c)", 1},
		{"a * b * c", "((a * b) * c)", 1},
		{"a * b / c", "((a * b) / c)", 1},
		{"a + b / c", "(a + (b / c))", 1},
		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)", 1},
		{"3 + 4; -f * 5", "(3 + 4)((-f) * 5)", 2},
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))", 1},
		{"5 < 4 != 3 > 4", "((5 < 4) != (3 > 4))", 1},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))", 1},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("Operator precedence for %q", tt.input), func(t *testing.T) {
			program := parseAndTestCommonStep(t, tt.input, tt.expectedStatements)

			actual := program.String()
			if actual != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, actual)
			}
		})
	}
}

func testIntegerLiteral(t *testing.T, expression ast.Expression, value int64) bool {
	integ, ok := expression.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("righ expression is not an integer literal. Got %T", expression)
		return false
	}
	if integ.Value != value {
		t.Errorf("integ.Value is not %d. Got %d", value, integ.Value)
		return false
	}
	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral is not %d. got=%s", value, integ.TokenLiteral())
		return false
	}
	return true
}

func parseAndTestCommonStep(t *testing.T, input string, expectedStatements int) *ast.Program {
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParsedErrors(t, p)

	if program == nil {
		t.Fatalf("ParsedProgram() returned Nil")
	}
	if len(program.Statements) != expectedStatements {
		t.Fatalf("Program does not contain %d statement. Got %d", expectedStatements, len(program.Statements))
	}
	return program
}

func checkParsedErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}
	t.Errorf("parser has %d errors", len(errors))
	for _, e := range errors {
		t.Errorf("parser error %q", e)
	}
	t.FailNow()
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral() not 'let', got %q", s.TokenLiteral())
		return false
	}
	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. Got %T", s)
		return false
	}
	if letStmt.Name.Value != name {
		t.Errorf("let statement's name value not %s. Got %s", name, letStmt.Name.Value)
		return false
	}
	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("let statement's name value not %s. Got %s", name, letStmt.Name.Token.Literal)
		return false
	}
	return true
}
