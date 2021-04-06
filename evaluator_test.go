package main

import (
	"fmt"
	"maja/pkg/ast"
	"math/big"
	"testing"
)

func TestEvaluatorAst(t *testing.T) {

	tests := []struct {
		ast      ast.Cons
		expected ast.Cons
	}{
		{ast: ast.NewList(ast.ConsList{
			ast.NewSymbol("begin"),
			ast.NewList(ast.ConsList{ast.NewSymbol("define"), ast.NewSymbol("r"), ast.NewNumber(big.NewInt(10))}), // (define r 10)
			ast.NewList(ast.ConsList{ast.NewSymbol("+"), ast.NewSymbol("r"), ast.NewSymbol("r")}),                 // (+ r r)
		}), expected: ast.NewNumber(big.NewInt(20))},

		{ast.NewList(ast.ConsList{}), ast.NewList(ast.ConsList{})},
		{ast.NewNumber(big.NewInt(10)), ast.NewNumber(big.NewInt(10))},

		// (if (= 10 20) "foo" "bar")
		// ConsList
		{
			ast: ast.NewList(ast.ConsList{
				ast.NewSymbol("if"),
				ast.NewList(ast.ConsList{ast.NewSymbol("="), ast.NewNumber(big.NewInt(10)), ast.NewNumber(big.NewInt(20))}),
				ast.NewString("foo"), // #t
				ast.NewString("bar"), // #f
			}),
			expected: ast.NewString("bar")},
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
		case ast.Number:
			if got.Number.Cmp(test.expected.Number) != 0 {
				t.Errorf("test[%02d] Number -- expected=%+v. got=%+v", idx, got, test.expected)
			}
		case ast.String, ast.Symbol:
			if got.Value != test.expected.Value {
				t.Errorf("test[%02d] String -- expected=%+v. got=%+v", idx, got.Value, test.expected.Value)
			}
		case ast.Pair:
			if len(got.List) != len(test.expected.List) {
				t.Errorf("test[%02d] Pair -- expected=%+v. got=%+v", idx, len(got.List), len(test.expected.List))
			}
		default:
			t.Fatalf("received %s %+v\n", got.Type, got)
		}

	}
}
