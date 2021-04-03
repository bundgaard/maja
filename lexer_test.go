package main

import "testing"

func TestLexer(t *testing.T) {
	input := `(+ 1 2 3 4)`
	l := NewLexer(input)
	if l.NextToken() != "(" {
		t.Fail()
	}
}

func TestLexerWithString(t *testing.T) {
	input := `"hej"`
	l := NewLexer(input)
	actual := l.NextToken()
	if actual != "\"hej\"" {
		t.Error(actual, "\"hej\"")
	}
}
