package main

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

type Lexer struct {
	data        string
	ch          rune
	peek        rune
	index       int
	nextIndex   int
	hasVertical bool
}

func (l *Lexer) readChar() {
	var r rune
	var size int
	if l.nextIndex >= len(l.data) {
		l.ch = 0
	} else {
		r, size = utf8.DecodeRuneInString(l.data[l.nextIndex:])
		l.ch = r
		r, _ = utf8.DecodeRuneInString(l.data[l.nextIndex+size:])
		l.peek = r
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
	for unicode.IsSpace(l.ch) {
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
		} else if unicode.IsDigit(l.ch) {
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
	for l.ch != '"' && l.ch != 0 {
		l.readChar()
	}
	token := l.data[index:l.index]

	l.readChar() // eat last "
	return fmt.Sprintf("\"%s\"", token)

}

func (l *Lexer) readNumber() string {
	index := l.index
	for unicode.IsDigit(l.ch) {
		l.readChar()
	}
	return l.data[index:l.index]
}
func (l *Lexer) readIdentifier() string {
	index := l.index

	if l.ch == rune('|') {
		l.hasVertical = true
	}
	l.readChar() // '|'
	for isLetter(l.ch) || unicode.IsDigit(l.ch) ||
		unicode.IsSymbol(l.ch) ||
		l.ch == rune('!') ||
		l.ch == rune('.') ||
		l.ch == rune('+') ||
		l.ch == rune('-') ||
		l.ch == rune('*') ||
		l.ch == rune('/') ||
		(l.hasVertical && unicode.IsSpace(l.ch)) ||
		l.ch == rune(';') ||
		l.ch == rune('?') {
		if l.ch == rune('|') {
			l.hasVertical = false
		}
		l.readChar()
	}

	token := l.data[index:l.index]
	return token
}

func isLetter(ch rune) bool {
	return unicode.IsLetter(ch) ||
		unicode.IsSymbol(ch) ||
		ch == rune('-') ||
		ch == rune('*') ||
		ch == rune('/') ||
		ch == rune('.')
}
