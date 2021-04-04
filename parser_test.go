package main

import (
	"testing"
)

func TestParse(t *testing.T) {

	input := `(define r 10)`
	l := NewLexer(input)
	p := NewParser(l)

	cons := p.Parse()
	if cons.Type != Pair {
		t.Errorf("expected %q. got=%q", Pair, cons.Type)
	}

	if len(cons.List) != 3 {
		t.Errorf("expexted length %d. got %d", 3, len(cons.List))
	}

	if cons.List[0].Type != Symbol {
		t.Errorf("expected %q. got=%q", Symbol, cons.List[0].Type)
	}
	if cons.List[1].Type != Symbol {
		t.Errorf("expected %q. got=%q", Symbol, cons.List[1].Type)
	}

	if cons.List[2].Type != Number {
		t.Errorf("expected %q. got=%q", Number, cons.List[2].Type)
	}
}

func TestParseWithBegin(t *testing.T) {
	input := `(begin (define r 10) (+ r r))`
	expected := "[begin [define r 10] [+ r r]]"
	l := NewLexer(input)
	p := NewParser(l)

	cons := p.Parse()
	if expected != cons.String() {
		t.Fatal()
	}
}

func TestParseWithSetBang(t *testing.T) {
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
		l := NewLexer(test.input)
		p := NewParser(l)

		got := p.current
		if got != test.expected {
			t.Errorf("test [%03d] - expected %q. got=%q", tidx, test.expected, got)
		}
	}

}
