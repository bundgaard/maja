package main

import (
	"fmt"
	"math/big"
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
	Continuation
	Foreign
	Character
	Port
	Vector
	Macro
	Promise
	Environment
	LastSystemType
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
	case Pair:
		return fmt.Sprint(cs.List)
	default:
		return cs.Value
	}

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
