package evaluator

import (
	"fmt"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * ((3 * 3)) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s", tt.input), func(t *testing.T) {
			evaluated := testEval(tt.input)
			testIntegerObject(t, evaluated, tt.expected)
		})
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true != false", true},
		{"false != true", true},
		{"false != false", false},
		{"true != true", false},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s", tt.input), func(t *testing.T) {
			evaluated := testEval(tt.input)
			testBooleanObject(t, evaluated, tt.expected)
		})
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s", tt.input), func(t *testing.T) {
			e := testEval(tt.input)
			integer, ok := tt.expected.(int)
			if ok {
				testIntegerObject(t, e, int64(integer))
			} else {
				testNullObjects(t, e)
			}
		})
	}
}

func TestReturnValue(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{`if (10 > 1) {
					if (10 > 1) {
						return 10;	
					}
					return 1;
				}
`, 10},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s", tt.input), func(t *testing.T) {
			evaluated := testEval(tt.input)
			testIntegerObject(t, evaluated, tt.expected)
		})
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{"5 + true", "type mismatch: INTEGER + BOOLEAN"},
		{"5 + true; 5;", "type mismatch: INTEGER + BOOLEAN"},
		{"-true", "unknown operator: -BOOLEAN"},
		{"true + false", "unknown operator: BOOLEAN + BOOLEAN"},
		{"5;true + false;5;", "unknown operator: BOOLEAN + BOOLEAN"},
		{"if (10 > 1) { true + false; }", "unknown operator: BOOLEAN + BOOLEAN"},
		{"if (10 > 1) { if (10 > 1) { return true + false; } return 1 } ", "unknown operator: BOOLEAN + BOOLEAN"},
		{"foobar", "identifier not found: foobar"},
		{`"a" - "c"`, "unknown operator: STRING - STRING"},
		{`{"name": "monkey"}[fn(x) { x }];`, "unusable as hash key: FUNCTION"},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s", tt.input), func(t *testing.T) {
			evaluated := testEval(tt.input)
			err, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("no error object returned. Got %T (%+v)", evaluated, evaluated)
				return
			}
			if err.Message != tt.expectedMessage {
				t.Errorf("Wrong error message. Expected %q, got %q", tt.expectedMessage, err.Message)
			}
		})
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s", tt.input), func(t *testing.T) {
			testIntegerObject(t, testEval(tt.input), tt.expected)
		})
	}
}

func testNullObjects(t *testing.T, e object.Object) bool {
	if e != NULL {
		t.Errorf("object is not NULL. Got %T (%+v)", e, e)
		return false
	}
	return true
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s", tt.input), func(t *testing.T) {
			evaluated := testEval(tt.input)
			testBooleanObject(t, evaluated, tt.expected)
		})
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2 };"

	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("Object is not a function. Got %T (%+v)", evaluated, evaluated)
	}
	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters %+v", fn.Parameters)
	}
	if fn.Parameters[0].String() != "x" {
		t.Fatalf("patameter is not 'x'. Got %q", fn.Parameters[0])
	}
	expectedBody := "(x + 2)"
	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. Got %q", expectedBody, fn.Body.String())
	}
}

func TestFunctionCall(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let identity = fn(x) { return x; }; identity(5);", 5},
		{"let double = fn(x) { x * 2; }; double(5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5, 6);", 11},
		{"let add = fn(x, y) { x + y; }; add(5, 6 + 1);", 12},
		{"let add = fn(x, y) { x + y; }; add(5, add(4, 7));", 16},
		{"fn(x) { x; }(5)", 5},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s", tt.input), func(t *testing.T) {
			testIntegerObject(t, testEval(tt.input), tt.expected)
		})
	}
}

func TestClosures(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{`
let newAdder = fn(x) {
	fn(y) { x + y }
};
let addTwo = newAdder(2);
addTwo(2)
`, 4},
		{`
let add = fn(a, b) { a + b };
let sub = fn(a, b) { a - b };
let apply = fn(a, b, func) { func(a, b) };
apply(2, apply(3, 4, sub), add);
`, 1},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s", tt.input), func(t *testing.T) {
			testIntegerObject(t, testEval(tt.input), tt.expected)
		})
	}
}

func TestStringLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"Yo !"`, "Yo !"},
		{`"Yo" + " " +  "Ho"`, "Yo Ho"},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s", tt.input), func(t *testing.T) {
			eval := testEval(tt.input)
			s, ok := eval.(*object.String)
			if !ok {
				t.Fatalf("Object is not a string. Got %T (%+v)", eval, eval)
			}
			if s.Value != tt.expected {
				t.Errorf("String has wrong value. Got %q", s.Value)
			}
		})
	}
}

func TestBuiltinFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		// String
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("four two")`, 8},
		{`len(1)`, "argument to `len` not supported, got INTEGER"},
		{`len("1", "2")`, "wrong number of arguments. got=2, want=1"},
		// Array
		{`len([])`, 0},
		{`len([33])`, 1},
		{`first([33])`, 33},
		{`first([33, 34])`, 33},
		//{`first([])`, nil},
		{`first([], [])`, "wrong number of arguments. got=2, want=1"},
		{`first()`, "wrong number of arguments. got=0, want=1"},
		{`first("")`, "argument to `first` not supported, got STRING"},
		{`last([], [])`, "wrong number of arguments. got=2, want=1"},
		{`last()`, "wrong number of arguments. got=0, want=1"},
		{`last("")`, "argument to `last` not supported, got STRING"},
		{`last([33])`, 33},
		{`last([33, 34])`, 34},
		//{`last([])`, nil},
		{`rest([], [])`, "wrong number of arguments. got=2, want=1"},
		{`rest()`, "wrong number of arguments. got=0, want=1"},
		{`rest("")`, "argument to `rest` not supported, got STRING"},
		//{`rest([33])`, 33}, // TODO
		//{`rest([33, 34])`, [34]},
		//{`rest([33, 34, 35])`, [34, 35]},
		//{`last([])`, nil},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s", tt.input), func(t *testing.T) {
			eval := testEval(tt.input)
			switch expected := tt.expected.(type) {
			case int:
				testIntegerObject(t, eval, int64(expected))
			case string:
				err, ok := eval.(*object.Error)
				if !ok {
					t.Errorf("object is not Error. Got %T (%+v)", eval, eval)
					return
				}
				if err.Message != expected {
					t.Errorf("Wrong error message. Expected %q, got %q", expected, err.Message)
				}
			}
		})
	}
}

func TestArrayLiteral(t *testing.T) {
	input := "[1, 2 * 3, 4 + 5]"

	evaluated := testEval(input)
	a, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("Object is not an array. Got %T (%+v)", evaluated, evaluated)
	}
	if len(a.Elements) != 3 {
		t.Fatalf("array has wrong elements %+v", a.Elements)
	}
	testIntegerObject(t, a.Elements[0], 1)
	testIntegerObject(t, a.Elements[1], 6)
	testIntegerObject(t, a.Elements[2], 9)
}

func TestArrayIndexExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"[1, 2, 3][0]", 1},
		{"[1, 2, 3][1]", 2},
		{"[1, 2, 3][2]", 3},
		{"[1, 2, 3][1 + 1]", 3},
		{"let i = 0; [1, 2, 3][i]", 1},
		{"let myArray = [1, 2, 3]; myArray[2]", 3},
		{"let myArray = [1, 2, 3]; let i = myArray[0]; myArray[i]", 2},
		{"[1, 2, 3][3]", nil},
		{"[1, 2, 3][-1]", nil},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s", tt.input), func(t *testing.T) {
			eval := testEval(tt.input)
			integer, ok := tt.expected.(int)
			if ok {
				testIntegerObject(t, eval, int64(integer))
			} else {
				testNullObjects(t, eval)
			}
		})
	}
}

func TestHashLiterals(t *testing.T) {
	input := `
let two = "two";
{
	"one": 10 - 9,
	two: 1 + 1,
	"thr" + "ee": 6 / 2,
	4: 4,
	true: 5,
	false: 6
}
`

	eval := testEval(input)
	result, ok := eval.(*object.Hash)
	if !ok {
		t.Fatalf("Eval didn't return Hash. Got %T (%+v)", eval, eval)
	}
	expected := map[object.HashKey]int64{
		(&object.String{Value: "one"}).HashKey():   1,
		(&object.String{Value: "two"}).HashKey():   2,
		(&object.String{Value: "three"}).HashKey(): 3,
		(&object.Integer{Value: 4}).HashKey():      4,
		TRUE.HashKey():                             5,
		FALSE.HashKey():                            6,
	}
	if len(expected) != len(result.Pairs) {
		t.Fatalf("Hash has wrong number of pairs. Got %d", len(result.Pairs))
	}
	for expectedKey, expectedValue := range expected {
		pair, ok := result.Pairs[expectedKey]
		if !ok {
			t.Errorf("no pair for given key in Pairs %d", expectedKey.Value)
		}
		testIntegerObject(t, pair.Value, expectedValue)
	}

}

func TestHashIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`{"foo": 5}["foo"]`, 5},
		{`{"foo": 5}["bar"]`, nil},
		{`let key = "foo"; {"foo": 5}[key]`, 5},
		{`{}["bar"]`, nil},
		{`{5: 5}[5]`, 5},
		{`{true: 5}[true]`, 5},
		{`{false: 5}[false]`, 5},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s", tt.input), func(t *testing.T) {
			eval := testEval(tt.input)
			integer, ok := tt.expected.(int)
			if ok {
				testIntegerObject(t, eval, int64(integer))
			} else {
				testNullObjects(t, eval)
			}
		})
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
}

func testIntegerObject(t *testing.T, evaluated object.Object, expected int64) bool {
	result, ok := evaluated.(*object.Integer)
	if !ok {
		t.Errorf("Evaluated is not an Integer. Got %T (%+v)", evaluated, evaluated)
		return false
	}
	if result.Value != expected {
		t.Errorf("Evaluated has wrong value. Got %d, want %d", result.Value, expected)
		return false
	}
	return true
}
func testBooleanObject(t *testing.T, evaluated object.Object, expected bool) bool {
	result, ok := evaluated.(*object.Boolean)
	if !ok {
		t.Errorf("Evaluated is not an Boolean. Got %T (%+v)", evaluated, evaluated)
		return false
	}
	if result.Value != expected {
		t.Errorf("Evaluated has wrong value. Got %t, want %t", result.Value, expected)
		return false
	}
	return true
}
