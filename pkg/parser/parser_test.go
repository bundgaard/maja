package parser

import (
	"maja/pkg/ast"
	"maja/pkg/scanner"
	"math/big"
	"testing"
)

func TestAll(t *testing.T) {
	t.Run("Quotation", testQuotation)
	t.Run("Identifiers", testIdentifiers)
	t.Run("Inputs", testInput)
	t.Run("Numbers", testNumbers)
}

func testQuotation(t *testing.T) {
	tests := []struct {
		input    string
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
			ast.NewList(ast.ConsList{ast.NewSymbol("<"), ast.NewNumber(big.NewInt(1)), ast.NewNumber(big.NewInt(2))}),
			ast.NewString(`"true"`),
			ast.NewString(`"false"`),
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

		// Test 04 - '(1 2 3)
		{`'(1 2 3)`, ast.NewList(ast.ConsList{
			ast.NewSymbol("quote"),
			ast.NewNumber(big.NewInt(1)),
			ast.NewNumber(big.NewInt(2)),
			ast.NewNumber(big.NewInt(3))})},

		// (quote (1 2 3 4))

		{`(quote (1 2 3 4))`,
			ast.NewList(ast.ConsList{ // (
				ast.NewSymbol("quote"), // quote
				ast.NewList(ast.ConsList{ // (
					ast.NewNumber(1), // 1
					ast.NewNumber(2),
					ast.NewNumber(3),
					ast.NewNumber(4),
				}), // )
			}), // )
		}, // test 05

	}

	for idx, test := range tests {
		l := scanner.NewLexer(test.input)
		p := NewParser(l)

		got := p.Parse()
		if got.List[0].Type != test.expected.List[0].Type {
			t.Logf("test[%02d] structure\n%#v\n", idx, got)
			t.Errorf("test[%02d] type expected=%s. got=%s", idx, test.expected.List[0].Type, got.List[0].Type)
		}
		if !compareDeep(t, test.expected, got) {
			t.Errorf("test[%02d] failed", idx)
			return
		}

	}
}

func testIdentifiers(t *testing.T) {
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
func compareDeep(t *testing.T, expected ast.Cons, actual ast.Cons) bool {
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
			if !compareDeep(t, expectedCons, actualCons) {
				t.Fail()
			}
		}
		// compare each value to be the same
	default:
		t.Logf("not handled; Type=%s", actual.Type)
		return false
	}
	return true
}

func testInput(t *testing.T) {
	tests := []struct {
		input    string
		expected ast.Cons
	}{
		{`(eval '(1 2 3 4))`, ast.NewList(ast.ConsList{ast.NewSymbol("eval"), ast.NewSymbol("'"), ast.NewList(ast.ConsList{
			ast.NewNumber(big.NewInt(1)),
			ast.NewNumber(big.NewInt(2)),
			ast.NewNumber(big.NewInt(3)),
			ast.NewNumber(big.NewInt(4))})}),
		},
		{`"foo"`, ast.NewString(`"foo"`)},
		{`bar`, ast.NewSymbol("bar")},
	}

	for idx, test := range tests {
		l := scanner.NewLexer(test.input)
		p := NewParser(l)
		actual := p.Parse()
		if test.expected.Type != actual.Type {
			t.Fatalf("test[%02d] - Type expected=%q. got=%q", idx, test.expected.Type, actual.Type)
		}

		if test.expected.Value != actual.Value {
			t.Fatalf("test[%02d] - Value expected=%q. got=%q", idx, test.expected.Value, actual.Value)
		}
	}
}

func createBigInt(t *testing.T, value interface{}) *big.Int {
	zero := big.NewInt(0)
	switch x := value.(type) {
	case string:
		zero.SetString(x, 10)
	default:
		t.Fatalf("error: createBigInt unhandled type %T", x)
	}
	return zero
}

func testNumbers(t *testing.T) {
	tests := []struct {
		input    string
		expected ast.Cons
	}{
		{"12", ast.NewNumber(big.NewInt(12))},
		{"123123123123123123123", ast.NewNumber(createBigInt(t, "123123123123123123123"))},
	}

	for idx, test := range tests {
		l := scanner.NewLexer(test.input)
		p := NewParser(l)

		program := p.Parse()

		if test.expected.Type != program.Type {
			t.Errorf("test[%02d] - Type expected=%q. got=%q", idx, test.expected.Type, program.Type)
		}
	}
}
