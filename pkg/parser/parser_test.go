package parser

import (
	"fmt"
	"maja/pkg/ast"
	"maja/pkg/scanner"
	"math/big"
	"testing"
)

func TestQuote(t *testing.T) {
	tests := []struct {
		input string

		expected ast.Cons
	}{
		{"(1 2 3)", ast.Cons{
			List: ast.ConsList{
				ast.NewNumber(big.NewInt(1)),
				ast.NewNumber(big.NewInt(2)),
				ast.NewNumber(big.NewInt(3)),
			},
		},
		}, // end of (1 2 3)

		{`(if (< 1 2) "true" "false")`, ast.Cons{List: ast.ConsList{
			ast.NewSymbol("if"),
			ast.NewList(ast.ConsList{ast.NewSymbol("<"), ast.NewNumber(big.NewInt(1)), ast.NewNumber(big.NewInt((2)))}),
			ast.NewString("true"),
			ast.NewString("false"),
		}}}, // end of (if (< 1 2) "true" "false")

		{`(define r 10)`, ast.Cons{
			List: ast.ConsList{
				ast.NewSymbol("define"), ast.NewSymbol("r"), ast.NewNumber(big.NewInt(10)),
			}}}, // end of (define r 10)

		{`(begin (define r 10) (+ r r))`,
			ast.NewList(ast.ConsList{
				ast.NewSymbol("begin"),
				ast.NewList(ast.ConsList{ast.NewSymbol("define"), ast.NewSymbol("r"), ast.NewNumber(big.NewInt(10))}), // (define r 10)
				ast.NewList(ast.ConsList{ast.NewSymbol("+"), ast.NewSymbol("r"), ast.NewSymbol("r")})})},              // end of (begin (define r 10) (+ r r))

		// '(1 2 3)
		{`'(1 2 3)`, ast.NewList(ast.ConsList{
			ast.NewSymbol("'"),
			ast.NewList(ast.ConsList{
				ast.NewNumber(big.NewInt(1)),
				ast.NewNumber(big.NewInt(2)),
				ast.NewNumber(big.NewInt(3)),
			})})},

		// (quote (1 2 3 4))

		{`(quote (1 2 3 4))`, ast.NewList(ast.ConsList{
			ast.NewSymbol("quote"),
			ast.NewList(ast.ConsList{
				ast.NewNumber(big.NewInt(1)),
				ast.NewNumber(big.NewInt(2)),
				ast.NewNumber(big.NewInt(3)),
			})})},
	}

	for idx, test := range tests {
		l := scanner.NewLexer(test.input)
		p := NewParser(l)

		got := p.Parse()
		fmt.Printf("%#v\n", got)
		if len(got.List) != len(test.expected.List) {
			t.Errorf("test[%02d] len -- expected=%d. got=%d", idx, len(test.expected.List), len(got.List))
		}

		if got.List[0].Type != test.expected.List[0].Type {
			t.Errorf("test[%02d] -- expected=%q. got=%q", idx, test.expected.Type, got.Type)
		}

	}
}

func TestParseWithSetBang(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"x", "x"},
		{"set!", "set!"},
		{"❤", "❤"},
		{"#t", "#t"},
		{"#f", "#f"},
		{"modulo)", "modulo"},
		{"hej-med", "hej-med"},
		{"is-number?", "is-number?"},
	}

	for tidx, test := range tests {
		l := scanner.NewLexer(test.input)
		p := NewParser(l)

		got := p.current
		if got != test.expected {
			t.Errorf("test [%03d] - expected %q. got=%q", tidx, test.expected, got)
		}
	}

}
