package compiler

import (
	"fmt"
	"monkey/ast"
	"monkey/code"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"
)

type compilerTestCase struct {
	input                string
	expectedConstants    []any
	expectedInstructions []code.Instructions
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []compilerTestCase{
		{input: "1 + 2",
			expectedConstants: []any{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpContstant, 0),
				code.Make(code.OpContstant, 1),
				code.Make(code.OpAdd),
			},
		},
	}

	runCompilerTests(t, tests)
}

func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}

func runCompilerTests(t *testing.T, tests []compilerTestCase) {
	t.Helper()

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			program := parse(tt.input)

			compiler := New()
			err := compiler.Compile(program)
			if err != nil {
				t.Fatalf("Compiler error: %s", err)
			}

			bytecode := compiler.Bytecode()

			err = testInstructions(tt.expectedInstructions, bytecode.Instructions)
			if err != nil {
				t.Fatalf("Test instructions failed: %s", err)
			}

			err = testConstants(tt.expectedConstants, bytecode.Constants)
			if err != nil {
				t.Fatalf("Test constants failed: %s", err)
			}
		})
	}
}

func testInstructions(expected []code.Instructions, actual code.Instructions) error {
	concatted := concatInstructions(expected)

	if len(actual) != len(concatted) {
		return fmt.Errorf("wrong instructions length. want=%q, got=%q", concatted, actual)
	}
	for i, ins := range concatted {
		if actual[i] != ins {
			return fmt.Errorf("wrong instruction at %d. want=%q, got=%q", i, concatted, actual)
		}
	}
	return nil
}

func testConstants(expected []any, actual []object.Object) error {
	if len(expected) != len(actual) {
		return fmt.Errorf("wrong number of constants. got=%d, want=%d", len(actual), len(expected))
	}
	for i, constant := range expected {
		switch constant := constant.(type) {
		case int:
			err := testIntegerObject(int64(constant), actual[i])
			if err != nil {
				return fmt.Errorf("constant %d - testIntegerObject failed: %s", i, err)
			}
		}
	}
	return nil
}

func testIntegerObject(expected int64, actual object.Object) error {
	result, ok := actual.(*object.Integer)
	if !ok {
		return fmt.Errorf("object is not integer. got=%T (%+v)", actual, actual)
	}
	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
	}
	return nil
}

func concatInstructions(s []code.Instructions) code.Instructions {
	out := code.Instructions{}
	for _, ins := range s {
		out = append(out, ins...)
	}
	return out
}
