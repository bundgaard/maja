package main

import "fmt"

type Lexer struct {
	data      string
	ch        byte
	index     int
	nextIndex int
}

func (l *Lexer) readChar() {
	if l.nextIndex >= len(l.data) {
		l.ch = 0
	} else {
		l.ch = l.data[l.nextIndex]
	}
	l.index = l.nextIndex
	l.nextIndex++
}

func NewLexer(data string) *Lexer {
	l := &Lexer{data: data}
	l.readChar()
	return l
}

func isSpace(ch byte) bool {
	return ch == ' ' || ch == '\r' || ch == '\n' || ch == '\t'
}
func (l *Lexer) skipSpace() {
	if isSpace(l.ch) {
		l.readChar()
	}
}
func (l *Lexer) NextToken() string {
	l.skipSpace()
	out := ""

	switch l.ch {
	case '(':
		out = string(l.ch)
	case ')':
		out = string(l.ch)
	case '+', '-', '*', '/':
		out = string(l.ch)
	case 0:
		return ""
	default:
		if isLetter(l.ch) {
			out = l.readIdentifier()
			return out
		} else if isDigit(l.ch) {
			out = l.readNumber()
			return out
		} else {
			out = fmt.Sprintf("illegal char %c", l.ch)
		}
	}
	l.readChar()
	return out
}
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
func (l *Lexer) readNumber() string {
	index := l.index
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.data[index:l.index]
}
func (l *Lexer) readIdentifier() string {
	index := l.index
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.data[index:l.index]
}
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z'
}
func tokenize(data string) []string {
	l := NewLexer(data)
	result := make([]string, 0)
	for tok := l.NextToken(); tok != ""; tok = l.NextToken() {
		result = append(result, tok)
	}
	return result
}
