package lexer

type Scanner struct {
	buf       []byte
	current   byte
	index     int
	nextIndex int
}

func NewScanner(data string) *Scanner {
	s := &Scanner{buf: []byte(data)}
	s.readChar()
	return s
}

func (s *Scanner) readChar() {
	if s.nextIndex >= len(s.buf) {
		s.current = 0
	} else {
		s.current = s.buf[s.nextIndex]
	}
	s.index = s.nextIndex
	s.nextIndex++

}

func (s *Scanner) NextToken() *Token {
	s.skipSpace()
	var token Token

	switch s.current {
	case '(':
		token = Token{Type: TokenOpenParen, Literal: string(s.current)}
	case ')':
		token = Token{Type: TokenCloseParen, Literal: string(s.current)}
	case '+', '-', '/', '*':
		token = Token{Type: TokenIdentifier, Literal: string(s.current)}
	case '"':
		s.readChar() // EAT "
		token.Literal = s.readString()
		token.Type = TokenAtom
	case 0:
		token = Token{Type: TokenEOF, Literal: ""}
	default:
		if isLetter(s.current) {
			token.Literal = s.readIdentifier()
			token.Type = TokenIdentifier
			return &token
		} else if isNumber(s.current) {
			token.Literal = s.readNumber()
			token.Type = TokenNumber
			return &token
		} else {
			token = Token{Type: TokenILLEGAL, Literal: string(s.current)}
		}

	}

	s.readChar()
	return &token
}
func (s *Scanner) readString() string {
	index := s.index
	for s.current != '"' {
		s.readChar()
	}
	return string(s.buf[index:s.index])
}
func (s *Scanner) readIdentifier() string {
	index := s.index
	for isLetter(s.current) {
		s.readChar()
	}
	return string(s.buf[index:s.index])
}

func (s *Scanner) readNumber() string {
	index := s.index
	for isNumber(s.current) {
		s.readChar()
	}
	return string(s.buf[index:s.index])
}

func (s *Scanner) skipSpace() {
	for isSpace(s.current) {
		s.readChar()
	}
}
func isSpace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\r' || b == '\n'
}
func isNumber(b byte) bool {
	return '0' <= b && b <= '9'
}
func isLetter(b byte) bool {
	return 'a' <= b && b <= 'z' || 'A' <= b && b <= 'Z'
}
func isAlphanum(b byte) bool {
	return b == '_' || isNumber(b) || isLetter(b) // unicode.IsLetter(rune(b))
}
