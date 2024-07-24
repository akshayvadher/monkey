package evaluator

import (
	"fmt"
	"monkey/object"
	"testing"
)

func TestQuote(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`quote(5)`, `5`},
		{`quote(5 + 8)`, `(5 + 8)`},
		{`quote(foobar)`, `foobar`},
		{`quote(foobar + barfoo)`, `(foobar + barfoo)`},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s", tt.input), func(t *testing.T) {
			eval := testEval(tt.input)
			quote, ok := eval.(*object.Quote)
			if !ok {
				t.Fatalf("expected *object.Quote. Got %T (%+v)", eval, eval)
			}
			if quote.Node == nil {
				t.Fatalf("quote.Node is nil")
			}
			if quote.Node.String() != tt.expected {
				t.Errorf("not equal. got %q, want %q", quote.Node.String(), tt.expected)
			}
		})
	}
}

func TestQuoteUnquote(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`quote(unquote(4))`, `4`},
		{`quote(unquote(4 + 4))`, `8`},
		{`quote(8 + unquote(4 + 4))`, `(8 + 8)`},
		{`quote(unquote(4 + 4) + 9)`, `(8 + 9)`},
		{`let foobar = 8;
				quote(foobar)`,
			`foobar`,
		},
		{`let foobar = 8;
				quote(unquote(foobar))`,
			`8`,
		},
		{`quote(unquote(true))`, `true`},
		{`quote(unquote(false))`, `false`},
		{`quote(unquote(quote(4 + 4)))`, `(4 + 4)`},
		{`let quotedInfixExpression = quote(4 + 4);
				quote(unquote(4 + 4) + unquote(quotedInfixExpression))`,
			`(8 + (4 + 4))`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			eval := testEval(tt.input)
			quote, ok := eval.(*object.Quote)
			if !ok {
				t.Fatalf("expected *object.Quote. Got %T (%+v)", eval, eval)
			}
			if quote.Node == nil {
				t.Fatalf("quote.Node is nil")
			}
			if quote.Node.String() != tt.expected {
				t.Errorf("not equal. got %q, want %q", quote.Node.String(), tt.expected)
			}
		})
	}
}
