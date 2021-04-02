package eval

import (
	"fmt"
	"maja/internal/parser"
)

type Eval struct {
	p *parser.Parser
}

func NewEval(data string) *Eval {
	return &Eval{p: parser.NewParser(data)}
}

func (e *Eval) Eval() {

	fmt.Println(e.p.SExpression())
	for expr := e.p.SExpression(); expr != nil; expr = expr.Cdr() {
		fmt.Println(expr)
	}

}
