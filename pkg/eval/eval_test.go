package eval

import (
	"fmt"
	"maja/pkg/ast"
	"maja/pkg/parser"
	"maja/pkg/scanner"
	"math/big"
	"testing"
)

func TestEvaluation(t *testing.T) {

	tests := []struct {
		input    string
		expected ast.Cons
	}{
		{`'(1 2 3 4)'`, ast.NewList(ast.ConsList{
			ast.NewNumber(1),
			ast.NewNumber(2),
			ast.NewNumber(3),
			ast.NewNumber(4)}),
		}, // test 00
	}

	for idx, test := range tests {

		l := scanner.NewLexer(test.input)
		p := parser.NewParser(l)

		program := p.Parse()

		output, err := Evalautor(program, StandardEnvironment())
		if err != nil {
			t.Error(err)
		}

		if !deepCompare(t, test.expected, output) {
			t.Errorf("test[%02d] deepcompare failed", idx)
			t.Errorf("test[%02d] program %#v", idx, program)

		}
		t.Logf("test[%02d] END", idx)
	}
}

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

	env := StandardEnvironment()
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

func TestInputToAst(t *testing.T) {

	tests := []struct {
		input        string
		cons         ast.Cons
		expectedType ast.Type
	}{

		{`+`, ast.Cons{Type: ast.Proc, Proc: add}, ast.Proc},
		{`(+ 1 2 3)`, ast.NewNumber(6), ast.Number},
		// {`(+)`, ast.Cons{Type: ast.Error, Value: "error: nothing "}, ast.Error},
		{`'(1 2 3 4)`, ast.Cons{Type: ast.Pair, List: ast.ConsList{ast.NewNumber(big.NewInt(1)), ast.NewNumber(big.NewInt(2)), ast.NewNumber(big.NewInt(3)), ast.NewNumber(big.NewInt(4))}}, ast.Pair},
		{`(quote (1 2 3 4))`, ast.Cons{Type: ast.Pair, List: ast.ConsList{ast.NewNumber(big.NewInt(1)), ast.NewNumber(big.NewInt(2)), ast.NewNumber(big.NewInt(3)), ast.NewNumber(big.NewInt(4))}}, ast.Pair},
		{`(lambda (n) (* 2 n))`, ast.Cons{Type: ast.Closure, List: ast.ConsList{ast.NewSymbol("lambda"), ast.NewList(ast.ConsList{ast.NewSymbol("n")}), ast.NewList(ast.ConsList{ast.NewSymbol("*"), ast.NewNumber(big.NewInt(2)), ast.NewSymbol("n")})}}, ast.Closure},
		{`((lambda (n) (* 2 n)) 2)`, ast.Cons{Type: ast.Number, Number: big.NewInt(4)}, ast.Number},

		{`(begin (define r 10) (+ r r))`, ast.Cons{Type: ast.Number, Number: big.NewInt(20)}, ast.Number},
	}

	for idx, test := range tests {
		l := scanner.NewLexer(test.input)
		p := parser.NewParser(l)
		program := p.Parse()
		output, err := Evalautor(program, StandardEnvironment())
		if err != nil {
			t.Fatalf("%v", err)
		}
		if output.Type != test.expectedType {
			t.Fatalf("test[%02d] - Type expected=%q. got=%q", idx, test.expectedType.String(), output.Type.String())
		}

		if !deepCompare(t, test.cons, output) {
			return
		}
	}
}

func deepCompare(t *testing.T, expected ast.Cons, actual ast.Cons) bool {
	switch actual.Type {

	case ast.String, ast.Error:
		if actual.Value != expected.Value {
			t.Errorf("String expected=%q. got=%q", expected.Value, actual.Value)
			return false
		}
	case ast.Symbol:
		if actual.Value != expected.Value {
			t.Errorf("Symbol expected=%q. got=%q", expected.Value, actual.Value)
			return false
		}
	case ast.Number:
		if actual.Number.Cmp(expected.Number) != 0 {
			t.Errorf("Number expected=%d. got=%d", expected.Number, actual.Number)
			return false
		}
	case ast.Pair, ast.Closure:
		if len(actual.List) != len(expected.List) {
			t.Errorf("Pair expected=%d. got=%d", len(expected.List), len(actual.List))
			return false
		}
		for i := range actual.List {
			expectedCons := expected.List[i]
			actualCons := actual.List[i]
			if !deepCompare(t, expectedCons, actualCons) {
				t.Fail()
			}
		}
		// compare each value to be the same
	case ast.Proc:
		// don't know how to compare Virtual Functions.
		t.Log("skipping Proc, I have no idea how to compare function pointer")

	default:
		t.Logf("not handled; Type=%s", actual.Type)
		return false
	}
	return true
}
