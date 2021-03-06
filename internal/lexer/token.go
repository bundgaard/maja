package lexer

import (
	"math/big"
)

//go:generate stringer -type TokenType
type TokenType uint8

const (
	TokenEOF TokenType = iota
	TokenOpenParen
	TokenCloseParen
	TokenIllegal
	TokenAtom
	TokenBoolean
	TokenNumber
	TokenIdentifier
	TokenDot
	TokenChar
	TokenQuote
	TokenDQuote
	TokenBQuote
	TokenComma
	TokenNewline
)

type Token struct {
	Type    TokenType
	Literal string
	Num     *big.Int
}

func (t *Token) String() string {
	return t.Literal
}
