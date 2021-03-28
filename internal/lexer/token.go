package lexer

import (
	"math/big"
)

//go:generate stringer -type TokenType
type TokenType uint8

const (
	TokenEOF TokenType = iota
	TokenILLEGAL
	TokenAtom
	TokenBoolean
	TokenNumber
	TokenIdentifier
	TokenOpenParen
	TokenCloseParen
	TokenDot
	TokenChar
	TokenQuote
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
