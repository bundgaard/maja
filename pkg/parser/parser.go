package parser

import (
	"maja/pkg/ast"
	"maja/pkg/scanner"
	"math/big"
	"strconv"
)

type Parser struct {
	l       *scanner.Lexer
	current string
	peek    string
}

func NewParser(l *scanner.Lexer) *Parser {
	p := &Parser{l: l}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.current = p.peek
	p.peek = p.l.NextToken()
}

func (p *Parser) parseList() ast.Cons {
	l := make([]ast.Cons, 0)
	p.nextToken()      // eat (
	token := p.current // define
	if p.peek == "quote" || p.peek == "'" {
		return p.parseQuote()
	}
	for token != ")" && token != "EOF" {
		l = append(l, p.Parse())
		p.nextToken() // eat token
		token = p.current
	}
	return ast.NewList(l)
}

func (p *Parser) Parse() ast.Cons {
	var cons ast.Cons
	token := p.current
	for token == "COMMENT" {
		p.nextToken()
		token = p.current
	}
	for token != "EOF" {
		switch token {
		case "(":
			cons = p.parseList()
			return cons
		case "'":
			cons = p.parseQuote()
			return cons

		default:
			n, err := strconv.ParseInt(token, 0, 64)
			if err != nil {
				return ast.NewSymbol(token)
			}
			return ast.NewNumber(big.NewInt(n))
		}

	}

	return cons
}
func (p *Parser) parseQuote() ast.Cons {
	// '() -> Cons(List: Parse)
	p.nextToken() // eat ' or quote
	token := p.Parse()
	l := make(ast.ConsList, 0)
	l = append(l, ast.NewSymbol("quote"), token)
	out := ast.NewList(l)
	return out
}

/*
 current = (, next = *, tokens = ( define r 10)
parse ( -> parseList (, nextToken, append to list

*/
