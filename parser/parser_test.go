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
	expression := stmt.Expression
	testIdentifier(t, expression, "foobar")
}

func testIdentifier(t *testing.T, expression ast.Expression, value string) bool {
	ident, ok := expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp is not ast.Identifier, got %T", ident)
	}
	if ident.Value != value {
		t.Errorf("ident value is not %s. got %s", value, ident.Value)
	}
	if ident.TokenLiteral() != value {
		t.Errorf("token literal is not %s, got %s", value, ident.TokenLiteral())
	}
	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled. Got %T", exp)
	return false
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, b bool) bool {
	boolExp, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("righ expression is not an integer literal. Got %T", exp)
		return false
	}
	if boolExp.Value != b {
		t.Errorf("boolExp.Value is not %t. Got %t", b, boolExp.Value)
		return false
	}
	if boolExp.TokenLiteral() != fmt.Sprintf("%t", b) {
		t.Errorf("boolExp.TokenLiteral is not %t. got=%s", b, boolExp.TokenLiteral())
		return false
	}
	return true
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

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input           string
		expectedBoolean bool
	}{
		{"true", true},
		{"true;", true},
		{"false", false},
		{"false;", false},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("Boolean parsing for %s", tt.input), func(t *testing.T) {
			program := parseAndTestCommonStep(t, tt.input, 1)
			stmt := testAndParseToExpressionStatement(t, program)
			testLiteralExpression(t, stmt.Expression, tt.expectedBoolean)
		})
	}
}

func TestParsingPrefixExpression(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!true;", "!", true},
		{"!false;", "!", false},
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
			if !testLiteralExpression(t, exp.Right, tt.value) {
				return
			}
		})
	}
}

func TestParsingInfixExpression(t *testing.T) {
	infixTests := []struct {
		input    string
		left     interface{}
		operator string
		right    interface{}
	}{
		{"5 + 5", 5, "+", 5},
		{"5 - 5", 5, "-", 5},
		{"5 * 5", 5, "*", 5},
		{"5 / 5", 5, "/", 5},
		{"5 > 5", 5, ">", 5},
		{"5 < 5", 5, "<", 5},
		{"5 == 5", 5, "==", 5},
		{"5 != 5", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range infixTests {
		t.Run(fmt.Sprintf("Infix parsing for %s", tt.input), func(t *testing.T) {
			program := parseAndTestCommonStep(t, tt.input, 1)
			stmt := testAndParseToExpressionStatement(t, program)
			stmtExp := stmt.Expression
			if !testInfixExpression(t, stmtExp, tt.left, tt.operator, tt.right) {
				return
			}
		})
	}
}

func testInfixExpression(t *testing.T, stmtExp ast.Expression, left interface{}, operator string, right interface{}) bool {
	exp, ok := stmtExp.(*ast.InfixExpression)
	if !ok {
		t.Fatalf("exp is not an infix expression. Got %T", stmtExp)
	}
	if !testLiteralExpression(t, exp.Left, left) {
		return false
	}
	if exp.Operator != operator {
		t.Fatalf("expression operator is not %s. Got %s", operator, exp.Operator)
	}
	if !testLiteralExpression(t, exp.Right, right) {
		return false
	}
	return true
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
		{"true", "true", 1},
		{"false", "false", 1},
		{"3 > 5 == false", "((3 > 5) == false)", 1},
		{"3 < 5 == true", "((3 < 5) == true)", 1},
		{"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4)", 1},
		{"(5 + 5) * 2", "((5 + 5) * 2)", 1},
		{"2 / (5 + 5)", "(2 / (5 + 5))", 1},
		{"-(5 + 5)", "(-(5 + 5))", 1},
		{"!(true == true)", "(!(true == true))", 1},
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
