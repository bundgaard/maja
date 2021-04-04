package main

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

type Lexer struct {
	data      string
	ch        rune
	index     int
	nextIndex int
}

func (l *Lexer) readChar() {
	var r rune
	var size int
	if l.nextIndex >= len(l.data) {
		l.ch = 0
	} else {
		r, size = utf8.DecodeRuneInString(l.data[l.nextIndex:])
		l.ch = r
	}
	l.index = l.nextIndex
	l.nextIndex += size
}

func NewLexer(data string) *Lexer {
	l := &Lexer{data: data}
	l.readChar()
	return l
}

func (l *Lexer) skipSpace() {
	for isSpace(l.ch) {
		l.readChar()
	}
}

func (l *Lexer) NextToken() string {

	out := ""
	l.skipSpace()
	switch l.ch {

	case '(':
		out = string(l.ch)
	case ')':
		out = string(l.ch)
	case '+', '-', '*', '/':
		out = string(l.ch)
	case '#':
		out = l.readHash()
	case '\'':
		out = string(l.ch)
	case 0:
		return "EOF"
	default:
		if l.ch == '"' {
			out = l.readString()
			return out

		} else if isLetter(l.ch) {
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

func (l *Lexer) readHash() string {
	index := l.index
	l.readChar() // eat #
	if l.ch == 't' || l.ch == 'f' {
		l.readChar()
	}
	return l.data[index:l.index]
}

func (l *Lexer) readString() string {

	l.readChar() // eat "
	index := l.index
	for l.ch != '"' {
		l.readChar()
	}
	token := l.data[index:l.index]
	l.readChar() // eat last "
	return fmt.Sprintf("\"%s\"", token)

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

	for isLetter(l.ch) ||
		unicode.IsSymbol(l.ch) ||
		l.ch == rune('!') ||
		l.ch == rune('-') ||
		l.ch == rune('?') {
		l.readChar()
	}
	token := l.data[index:l.index]
	return token
}

func isDigit(ch rune) bool {
	return unicode.IsDigit(ch) // return '0' <= ch && ch <= '9'
}

func isSpace(ch rune) bool {
	return unicode.IsSpace(ch) //  == ' ' || ch == '\r' || ch == '\n' || ch == '\t'
}
func isLetter(ch rune) bool {
	return unicode.IsLetter(ch) || unicode.IsSymbol(ch) //return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z'
}
