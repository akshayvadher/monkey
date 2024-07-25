package evaluator

import (
	"monkey/ast"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"
)

func TestDefineMacros(t *testing.T) {
	input := `
		let number = 1;
		let function = fn(x, y) { x + y };
		let myMacro = macro(x, y) { x + y };
	`
	env := object.NewEnvironment()
	program := testParseProgram(input)

	DefineMacros(program, env)

	if len(program.Statements) != 2 {
		t.Fatalf("Wrong number of statemts. Got %d", len(program.Statements))
	}
	_, ok := env.Get("number")
	if ok {
		t.Fatalf("Number should not be defined")
	}
	_, ok = env.Get("function")
	if ok {
		t.Fatalf("'function' should not be defined")
	}
	obj, ok := env.Get("myMacro")
	if !ok {
		t.Fatalf("macro myMacro not in env")
	}
	macro, ok := obj.(*object.Macro)
	if !ok {
		t.Fatalf("object is not Macro. Got %T (%+v)", obj, obj)
	}
	if len(macro.Parameters) != 2 {
		t.Fatalf("Wrong number of macro parameters. Got %d", len(macro.Parameters))
	}
	if macro.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. Got %q", macro.Parameters[0])
	}
	if macro.Parameters[1].String() != "y" {
		t.Fatalf("parameter is not 'y'. Got %q", macro.Parameters[1])
	}
	expectedBody := "(x + y)"
	if macro.Body.String() != expectedBody {
		t.Fatalf("body is not %q. Got %q", expectedBody, macro.Body.String())
	}
}

func testParseProgram(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}

func TestExpandMacro(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`let infixExpression = macro () { quote(1 + 2); }; infixExpression();`, `(1 + 2)`},
		{`
				let reverse = macro(a, b) { quote(unquote(b) - unquote(a)); };
				reverse(2 + 2, 10 - 5);`,
			`(10 - 5) - (2 + 2)`,
		},
		{
			`
			let unless = macro(condition, consequence, alternative) {
				quote(if (!(unquote(condition))) {
					unquote(consequence);
				} else {
					unquote(alternative);
				});
			};
			unless(10 > 5, puts("not greater"), puts("greater"));`,
			`if(!(10>5)){puts("not greater")} else { puts("greater") }`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			expected := testParseProgram(tt.expected)
			input := testParseProgram(tt.input)

			env := object.NewEnvironment()
			DefineMacros(input, env)
			expanded := ExpandMacros(input, env)

			if expanded.String() != expected.String() {
				t.Errorf("expanded macros are not equal. Got %q, want %q", expanded.String(), expected.String())
			}
		})
	}
}
