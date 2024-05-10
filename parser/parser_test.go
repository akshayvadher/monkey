package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"testing"
)

func TestLetStatement(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Test let statement for %s", tt.input), func(t *testing.T) {
			program := parseAndTestCommonStep(t, tt.input, 1)
			stmt := program.Statements[0]
			if !testLetStatement(t, stmt, tt.expectedIdentifier) {
				return
			}
			val := stmt.(*ast.LetStatement).Value
			if !testLiteralExpression(t, val, tt.expectedValue) {
				return
			}
		})
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 5;", 5},
		{"return true;", true},
		{"return foobar;", "foobar"},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("Test return statement for %s", tt.input), func(t *testing.T) {
			program := parseAndTestCommonStep(t, tt.input, 1)
			returnStmt, ok := program.Statements[0].(*ast.ReturnStatement)
			if !ok {
				t.Errorf("stmt is not a *ast.ReturnStatement. Got %T", program.Statements[0])
				return
			}
			if returnStmt.TokenLiteral() != "return" {
				t.Errorf("returnStmt.TokenLiteral() is not 'return'. Got %q", returnStmt.TokenLiteral())
			}
			if !testLiteralExpression(t, returnStmt.ReturnValue, tt.expectedValue) {
				return
			}
		})
	}

}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	program := parseAndTestCommonStep(t, input, 1)

	stmt := parseAndTestExpressionStatement(t, program)
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

func parseAndTestExpressionStatement(t *testing.T, program *ast.Program) *ast.ExpressionStatement {
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statement[0] is not an expression statement, got %T", program.Statements[0])
	}
	return stmt
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	program := parseAndTestCommonStep(t, input, 1)

	stmt := parseAndTestExpressionStatement(t, program)

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
			stmt := parseAndTestExpressionStatement(t, program)
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
			stmt := parseAndTestExpressionStatement(t, program)
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
			stmt := parseAndTestExpressionStatement(t, program)
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
		{"a + add(b * c) + d", "((a + add((b * c)) ) + d)", 1},
		{"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))", "add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)) ) ", 1},
		{"add(a + b + c * d / f + g)", "add((((a + b) + ((c * d) / f)) + g)) ", 1},
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

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	program := parseAndTestCommonStep(t, input, 1)
	expStmt := parseAndTestExpressionStatement(t, program)
	ifExp, ok := expStmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("statement expression is not an if expression. Got %T", expStmt.Expression)
	}
	if !testInfixExpression(t, ifExp.Condition, "x", "<", "y") {
		return
	}
	if len(ifExp.Consequence.Statements) != 1 {
		t.Errorf("Consequence is not 1 statement. Got %d", len(ifExp.Consequence.Statements))
	}
	consequence, ok := ifExp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Consequence[0] statement is not an expression. Got %T", ifExp.Consequence.Statements[0])
	}
	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}
	if ifExp.Alternative != nil {
		t.Errorf("if expression alternative was not nil. Got %+v", ifExp.Alternative)
	}
}
func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	program := parseAndTestCommonStep(t, input, 1)
	expStmt := parseAndTestExpressionStatement(t, program)
	ifExp, ok := expStmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("statement expression is not an if expression. Got %T", expStmt.Expression)
	}
	if !testInfixExpression(t, ifExp.Condition, "x", "<", "y") {
		return
	}
	if len(ifExp.Consequence.Statements) != 1 {
		t.Errorf("Consequence is not 1 statement. Got %d", len(ifExp.Consequence.Statements))
	}
	consequence, ok := ifExp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Consequence[0] statement is not an expression. Got %T", ifExp.Consequence.Statements[0])
	}
	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}
	if len(ifExp.Alternative.Statements) != 1 {
		t.Errorf("Alternative is not 1 statement. Got %d", len(ifExp.Alternative.Statements))
	}
	altConsequence, ok := ifExp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Alternative[0] statement is not an expression. Got %T", ifExp.Alternative.Statements[0])
	}
	if !testIdentifier(t, altConsequence.Expression, "y") {
		return
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x, y) { x + y; }`
	program := parseAndTestCommonStep(t, input, 1)
	stmt := parseAndTestExpressionStatement(t, program)

	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("statement expression is not a functiona literal. Got %T", stmt.Expression)
	}
	if len(function.Parameters) != 2 {
		t.Fatalf("function literal is not 2. Got %d", len(function.Parameters))
	}
	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")
	if len(function.Body.Statements) != 1 {
		t.Fatalf("Function body does not have 1 statement. Got %d", len(function.Body.Statements))
	}
	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function body is not an expression statement. Got %T", function.Body.Statements[0])
	}
	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{"fn() {};", []string{}},
		{"fn(x) {};", []string{"x"}},
		{"fn(x, y, z) {};", []string{"x", "y", "z"}},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("Test function parameter parsing for %s", tt.input), func(t *testing.T) {
			program := parseAndTestCommonStep(t, tt.input, 1)
			function := program.Statements[0].(*ast.ExpressionStatement).Expression.(*ast.FunctionLiteral)

			if len(function.Parameters) != len(tt.expectedParams) {
				t.Errorf("parameters lenth wrong. Want %d. Got %d", len(tt.expectedParams), len(function.Parameters))
			}
			for i, ident := range tt.expectedParams {
				testLiteralExpression(t, function.Parameters[i], ident)
			}
		})
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := `add(1, 2 * 3, 4 + 5);`
	program := parseAndTestCommonStep(t, input, 1)
	stmt := parseAndTestExpressionStatement(t, program)
	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("statement expression is not a call expression. Got %T", stmt.Expression)
	}
	if !testIdentifier(t, exp.Function, "add") {
		return
	}
	if len(exp.Arguments) != 3 {
		t.Fatalf("Wrong length of arguments. Got %d", len(exp.Arguments))
	}
	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world"`
	program := parseAndTestCommonStep(t, input, 1)
	stmt := parseAndTestExpressionStatement(t, program)
	s, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("statement expression is not a string literal. Got %T", stmt.Expression)
	}
	if s.Value != "hello world" {
		t.Errorf("string literal value is not %q. Got %q", "hello world", s.Value)
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
