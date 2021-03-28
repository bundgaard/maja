package lexer

import "testing"

func TestNextToken(t *testing.T) {
	s := NewScanner(`(+ 1 2 3 4 "hello world")`)
	tests := []struct {
		Expected TokenType
		Got      *Token
	}{
		{Expected: TokenOpenParen, Got: s.NextToken()},
		{Expected: TokenIdentifier, Got: s.NextToken()},
		{Expected: TokenNumber, Got: s.NextToken()},
		{Expected: TokenNumber, Got: s.NextToken()},
		{Expected: TokenNumber, Got: s.NextToken()},
		{Expected: TokenNumber, Got: s.NextToken()},
		{Expected: TokenAtom, Got: s.NextToken()},
		{Expected: TokenCloseParen, Got: s.NextToken()},
		{Expected: TokenEOF, Got: s.NextToken()},
	}
	for _, test := range tests {
		if test.Got.Type != test.Expected {
			t.Errorf("expected=%v, got=%+v", test.Expected, test.Got.Type)
		}
	}

}
