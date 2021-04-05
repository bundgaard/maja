package main

import "testing"

func TestLexer(t *testing.T) {

	tests := []struct {
		input    string
		expected string
	}{
		{"...", "..."},
		{"+", "+"},
		{"+soup+", "+soup+"},
		{"<=?", "<=?"},
		{"->string", "->string"},
		{"a34kTMNs", "a34kTMNs"},
		{"lambda", "lambda"},
		{"list->vector", "list->vector"},
		{"q", "q"},
		{"V17a", "V17a"},
		{"|two words|", "|two words|"},
		{"|two\x20;words|", "|two\x20;words|"},
		{"the-word-recursion-has-many-meanings", "the-word-recursion-has-many-meanings"},
		{"let*", "let*"},
	}

	for idx, test := range tests {
		l := NewLexer(test.input)
		got := l.NextToken()
		if got != test.expected {
			t.Errorf("test[%02d] -- expected=%q. got=%q", idx, test.expected, got)
		}
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
