package main

import (
	"fmt"
	"math/big"
	"testing"
)

func TestQuote(t *testing.T) {
	tests := []struct {
		input string

		expected Expr
	}{
		{"(1 2 3)", Expr{
			List: ConsList{
				NewNumber(big.NewInt(1)),
				NewNumber(big.NewInt(2)),
				NewNumber(big.NewInt(3)),
			},
		},
		}, // end of (1 2 3)

		{`(if (< 1 2) "true" "false")`, Expr{List: ConsList{
			NewSymbol("if"),
			NewList(ConsList{NewSymbol("<"), NewNumber(big.NewInt(1)), NewNumber(big.NewInt((2)))}),
			NewString("true"),
			NewString("false"),
		}}}, // end of (if (< 1 2) "true" "false")

		{`(define r 10)`, Expr{
			List: ConsList{
				NewSymbol("define"), NewSymbol("r"), NewNumber(big.NewInt(10)),
			}}}, // end of (define r 10)

		{`(begin (define r 10) (+ r r))`, Expr{List: ConsList{
			NewSymbol("begin"),
			NewList(ConsList{NewSymbol("define"), NewSymbol("r"), NewNumber(big.NewInt(10))}), // (define r 10)
			NewList(ConsList{NewSymbol("+"), NewSymbol("r"), NewSymbol("r")}),                 // (+ r r)
		}}}, // end of (begin (define r 10) (+ r r))
	}

	for idx, test := range tests {
		l := NewLexer(test.input)
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
		l := NewLexer(test.input)
		p := NewParser(l)

		got := p.current
		if got != test.expected {
			t.Errorf("test [%03d] - expected %q. got=%q", tidx, test.expected, got)
		}
	}

}
