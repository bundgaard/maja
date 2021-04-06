package parser

import (
	"maja/pkg/ast"
	"maja/pkg/scanner"
	"unicode"
	"unicode/utf8"
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
	token := p.current
	switch token {
	case "EOF":
		panic("unexpected EOF")
	case "(":
		cl := make(ast.ConsList, 0)
		for {
			p.nextToken()
			token = p.current

			if token == ")" {
				return ast.NewList(cl)
			}
			cl = append(cl, p.Parse())
		}
	case ")":
		panic("unexpected )")
	case "'":
		p.nextToken() // eat quote
		out := ast.NewSymbol("quote")
		tokens := p.Parse()
		cl := make(ast.ConsList, len(tokens.List)+1)
		cl[0] = out
		copy(cl[1:], tokens.List)
		return ast.NewList(cl)
	default:
		chr, _ := utf8.DecodeRuneInString(token)
		switch {
		case unicode.IsDigit(chr):
			return ast.NewNumber(token)
		case chr == '"':
			return ast.NewString(token)
		default:
			return ast.NewSymbol(token)
		}

	}
}

func (p *Parser) parseQuote() ast.Cons {
	// '(args) --> (quote (args))
	// (quote (args))
	token := p.Parse()
	l := make(ast.ConsList, 0)
	l = append(l, token)
	out := ast.NewList(l)
	return out
}
