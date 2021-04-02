package main

import "strconv"

type Parser struct {
	l       *Lexer
	current string
	peek    string

	errors []string
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

func (p *Parser) parseList() Cons {
	l := make([]Cons, 0)
	p.nextToken()      // eat (
	token := p.current // define
	for token != ")" && token != "EOF" {
		l = append(l, p.Parse())
		p.nextToken() // eat token
		token = p.current
	}
	return NewList(l)
}

func (p *Parser) Parse() Cons {
	var cons Cons
	token := p.current
	for token != "EOF" {
		switch token {
		case "(":
			cons = p.parseList()
			return cons
		default:
			n, err := strconv.ParseInt(token, 0, 64)
			if err != nil {
				return NewSymbol(token)
			}
			return NewNumber(n)

		}

	}

	return cons
}

/*
 current = (, next = *, tokens = ( define r 10)
parse ( -> parseList (, nextToken, append to list

*/
