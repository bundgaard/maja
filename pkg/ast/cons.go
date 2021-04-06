package ast

import (
	"fmt"
	"log"
	"math/big"
	"strings"
)

//go:generate stringer -type Type
type Type int

const (
	String Type = iota
	Number
	Symbol
	Proc
	Pair
	Closure
	Error
)

type ConsList []Cons
type Cons struct {
	Type   Type
	Value  string
	Number *big.Int
	List   ConsList
	Proc   func(ConsList) Cons
}

func createBigInt(value interface{}) *big.Int {
	zero := big.NewInt(0)
	switch x := value.(type) {
	case string:
		zero.SetString(x, 10)
	case *big.Int:
		return x
	case int:
		zero.SetInt64(int64(x))
	default:
		log.Fatalf("error: createBigInt unhandled type %T", x)
	}
	return zero
}

func (cs Cons) String() string {
	switch cs.Type {
	case Number:
		return fmt.Sprint(cs.Number)
	case Pair:
		return cs.ListToString()
	case Proc, Closure:
		return cs.Type.String()
	case String:
		return fmt.Sprintf("%s", cs.Value)
	default:
		return cs.Value
	}

}
func (cs Cons) ListToString() string {
	var out strings.Builder

	out.WriteString("(")
	for i, cons := range cs.List {
		if i > 0 && i < len(cs.List) {
			out.WriteString(" ")
		}
		out.WriteString(cons.String())

	}
	out.WriteString(")")
	return out.String()
}
func NewString(value string) Cons {
	return Cons{Type: String, Value: value}
}
func NewSymbol(value string) Cons {
	return Cons{Type: Symbol, Value: value}
}
func NewNumber(value interface{}) Cons {
	number := createBigInt(value)
	return Cons{Type: Number, Number: number}
}
func NewList(list []Cons) Cons {
	return Cons{Type: Pair, List: list}
}
func NewProc(fn func(ConsList) Cons) Cons {
	return Cons{Type: Proc, Proc: fn}
}
func NewError(err error) Cons {
	return Cons{Type: Error, Value: fmt.Sprintf("%v", err)}
}
func Arguments(cons Cons) ConsList {
	unproc := cons.List
	args := make(ConsList, 0)
	for i := 1; i < len(unproc); i++ {
		args = append(args, unproc[i])
	}
	return args

}
