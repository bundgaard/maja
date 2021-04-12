package eval

import (
	"maja/pkg/ast"
	"testing"
)

func TestAppend(t *testing.T) {

	env := StandardEnvironment()

	_ = env

}

func TestSubtract(t *testing.T) {

	tests := []struct {
		input    ast.ConsList
		expected ast.Cons
	}{
		{ast.ConsList{ast.NewNumber(1), ast.NewNumber(2)}, ast.NewNumber(-1)},
		{ast.ConsList{ast.NewNumber(5), ast.NewNumber(4)}, ast.NewNumber(1)},
	}

	for idx, test := range tests {
		got := subtract(test.input)

		if test.expected.Number.Cmp(got.Number) != 0 {

			t.Errorf("test[%02d] expected=%s. got=%s", idx, test.expected, got)
		}
	}
}

func TestAdd(t *testing.T) {
	tests := []struct {
		input    ast.ConsList
		expected ast.Cons
	}{
		{ast.ConsList{ast.NewNumber(1), ast.NewNumber(1), ast.NewNumber(1), ast.NewNumber(1), ast.NewNumber(1)}, ast.NewNumber(5)},
	}

	for idx, test := range tests {
		got := add(test.input)
		if test.expected.Number.Cmp(got.Number) != 0 {
			t.Errorf("test[%02d] expected=%s. got=%s", idx, test.expected.Number, got.Number)
		}
	}
}
