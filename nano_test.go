package main

import "testing"

func TestVerifyParenthesis(t *testing.T) {
	tests := []struct {
		input string
	}{
		{"()"},
		{"(())"},
		{"(()()())"},
	}

	for _, test := range tests {
		actual, err := verifyParenthesis(test.input)
		if err != nil {
			t.Fail()
		}
		if actual != test.input {
			t.Fail()
		}
	}

}
