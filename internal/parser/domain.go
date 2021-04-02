package parser

import (
	"maja/internal/lexer"
)

type Program struct {
	Expressions []Expression
}

type Expression struct {
	car  *Expression
	atom *lexer.Token
	cdr  *Expression
}

func (e *Expression) String() string {
	if e == nil {
		return "nil"
	}

	if e.atom != nil {
		return e.atom.String()
	}

	str := "("
	str += e.car.String()
	str += " . "
	str += e.cdr.String()
	str += ")"
	return str
}
func (e *Expression) Car() *Expression {
	return e.car
}
func (e *Expression) Cdr() *Expression {
	return e.cdr
}

type Parser struct {
	scanner *lexer.Scanner
}

func NewParser(data string) *Parser {
	return &Parser{scanner: lexer.NewScanner(data)}
}
func (p *Parser) SExpression() *Expression {
	token := p.scanner.NextToken()

	switch token.Type {
	case lexer.TokenEOF:
		return nil
	case lexer.TokenAtom, lexer.TokenNumber:
		return &Expression{atom: token}
	case lexer.TokenIdentifier:
		return &Expression{atom: token}
	case lexer.TokenOpenParen:
		return p.list()
	}
	return nil
}

func (p *Parser) list() *Expression {
	token := p.scanner.NextToken()
	switch token.Type {
	case lexer.TokenAtom, lexer.TokenNumber:
		return &Expression{car: &Expression{atom: token}, cdr: p.list()}
	case lexer.TokenIdentifier:
		return &Expression{car: &Expression{atom: token}, cdr: p.list()}
	case lexer.TokenOpenParen:
		return &Expression{car: p.SExpression(), cdr: p.list()}
	case lexer.TokenCloseParen:
		return nil
	}
	panic("not supported")
}

func (p *Parser) Program() *Expression {
	return nil
}
