package main

import (
	"fmt"
	"math/big"
	"testing"
)

func TestEvaluator(t *testing.T) {

	tests := []struct {
		ast      Cons
		expected Cons
	}{
		{ast: Cons{List: ConsList{
			NewSymbol("begin"),
			NewList(ConsList{NewSymbol("define"), NewSymbol("r"), NewNumber(big.NewInt(10))}), // (define r 10)
			NewList(ConsList{NewSymbol("+"), NewSymbol("r"), NewSymbol("r")}),                 // (+ r r)
		}}, expected: NewNumber(big.NewInt(21))},
	}

	env := standardEnvironment()
	for idx, test := range tests {
		got, err := evaluate(test.ast, env)
		if err != nil {
			t.Fatal(err)
		}

		fmt.Printf("FISK    FISK \n \n \n %T\n%+v\n%#v", got, got, got)

		switch got.Type {
		case Number:
			if got.Number != test.expected.Number {
				t.Errorf("test[%02d] -- expected=%+v. got=%+v", idx, got, test.expected)
			}
		case String, Symbol:
			if got.Value != test.expected.Value {
				t.Errorf("test[%02d] -- expected=%+v. got=%+v", idx, got.Value, test.expected.Value)
			}
		default:
			t.Fatalf("received %s %+v\n", got.Type, got)
		}

	}
}
