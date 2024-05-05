package parser

import (
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

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statement[0] is not an expression statement, got %T", program.Statements[0])
	}
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

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	program := parseAndTestCommonStep(t, input, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("the statement is not an expression. Got %T", program.Statements[0])
	}

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
