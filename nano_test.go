package main

import (
	"fmt"
	"testing"
)

func TestVerifyParenthesis(t *testing.T) {
	tests := []struct {
		input string
	}{
		{"()"},
		{"(())"},
		{"(()()())"},
	}

	for _, test := range tests {
		actual, err := verifyParenthesis(test.input, nil)
		if err != nil {
			t.Fail()
		}
		if actual != test.input {
			t.Fail()
		}
	}

}

func TestNanoCompose(t *testing.T) {
	tests := []struct {
		input string
	}{
		{"()"},
		{"('(1 2 3 4 5))"},
	}
	c := NanoCompose(verifyParenthesis, translateQuote)
	for _, test := range tests {
		fmt.Println(c(test.input))
	}

}
