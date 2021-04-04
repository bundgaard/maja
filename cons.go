package main

import (
	"fmt"
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
)

type ConsList []Cons
type Cons struct {
	Type   Type
	Value  string
	Number *big.Int
	List   ConsList
	Proc   func(ConsList) Cons
}

func (cs Cons) String() string {
	switch cs.Type {
	case Number:
		return fmt.Sprint(cs.Number)
	case Pair, Closure:
		return cs.ListToString()
	case Proc:
		return cs.Type.String()
	default:
		return cs.Value
	}

}
func (cs Cons) ListToString() string {
	var out strings.Builder

	out.WriteString("( ")
	for _, cons := range cs.List {
		out.WriteString(cons.String() + " ")

	}
	out.WriteString(")")
	return out.String()
}
func NewSymbol(value string) Cons {
	return Cons{Type: Symbol, Value: value}
}
func NewNumber(number *big.Int) Cons {
	return Cons{Type: Number, Number: number}
}
func NewList(list []Cons) Cons {
	return Cons{Type: Pair, List: list}
}
func NewProc(fn func(ConsList) Cons) Cons {
	return Cons{Type: Proc, Proc: fn}
}

func args(cons Cons) ConsList {
	unproc := cons.List
	args := make(ConsList, 0)
	for i := 1; i < len(unproc); i++ {
		args = append(args, unproc[i])
	}
	return args

}
