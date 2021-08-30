package main

import (
	"math/big"
	"strconv"
)

type Parser struct {
	l       *Lexer
	current string
	peek    string
}

func NewParser(l *Lexer) *Parser {
	p := &Parser{l: l}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.current = p.peek
	p.peek = p.l.NextToken()
}

func (p *Parser) parseList() Expr {
	l := make([]Expr, 0)
	p.nextToken()      // eat (
	token := p.current // define
	for token != ")" && token != "EOF" {
		l = append(l, p.Parse())
		p.nextToken() // eat token
		token = p.current
	}
	return NewList(l)
}

func (p *Parser) Parse() Expr {
	var cons Expr
	token := p.current
	for token != "EOF" {
		switch token {
		case "(":
			cons = p.parseList()
			return cons
		case "'", "quote": // TODO potential bug [[[]]]
			cons = p.parseQuote()
			return cons
		default:
			n, err := strconv.ParseInt(token, 0, 64)
			if err != nil {
				return NewSymbol(token)
			}
			return NewNumber(big.NewInt(n))

		}

	}

	return cons
}
func (p *Parser) parseQuote() Expr {
	// '() -> Cons(List: Parse)
	p.nextToken() // eat '
	token := p.Parse()
	l := make(ConsList, 0)
	l = append(l, NewSymbol("'"), token)
	out := NewList(l)
	return out
}

/*
 current = (, next = *, tokens = ( define r 10)
parse ( -> parseList (, nextToken, append to list

*/
