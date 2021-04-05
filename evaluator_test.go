package main

import (
	"fmt"
	"math/big"
	"testing"
)

func TestEvaluatorAst(t *testing.T) {

	tests := []struct {
		ast      Cons
		expected Cons
	}{
		{ast: NewList(ConsList{
			NewSymbol("begin"),
			NewList(ConsList{NewSymbol("define"), NewSymbol("r"), NewNumber(big.NewInt(10))}), // (define r 10)
			NewList(ConsList{NewSymbol("+"), NewSymbol("r"), NewSymbol("r")}),                 // (+ r r)
		}), expected: NewNumber(big.NewInt(20))},

		{NewList(ConsList{}), NewList(ConsList{})},
		{NewNumber(big.NewInt(10)), NewNumber(big.NewInt(10))},

		// (if (= 10 20) "foo" "bar")
		// ConsList
		{
			ast: NewList(ConsList{
				NewSymbol("if"),
				NewList(ConsList{NewSymbol("="), NewNumber(big.NewInt(10)), NewNumber(big.NewInt(20))}),
				NewString("foo"), // #t
				NewString("bar"), // #f
			}),
			expected: NewString("bar")},
	}

	env := standardEnvironment()
	for idx, test := range tests {
		got, err := evaluate(test.ast, env)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Printf("test[%02d] got %s %s\n", idx, got.Type.String(), got.String())
		fmt.Printf("test[%02d] ast %s %s", idx, test.ast.Type.String(), test.ast.String())
		switch got.Type {
		case Number:
			if got.Number.Cmp(test.expected.Number) != 0 {
				t.Errorf("test[%02d] Number -- expected=%+v. got=%+v", idx, got, test.expected)
			}
		case String, Symbol:
			if got.Value != test.expected.Value {
				t.Errorf("test[%02d] String -- expected=%+v. got=%+v", idx, got.Value, test.expected.Value)
			}
		case Pair:
			if len(got.List) != len(test.expected.List) {
				t.Errorf("test[%02d] Pair -- expected=%+v. got=%+v", idx, len(got.List), len(test.expected.List))
			}
		default:
			t.Fatalf("received %s %+v\n", got.Type, got)
		}

	}
}
