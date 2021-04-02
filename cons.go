package main

import "fmt"

//go:generate stringer -type Type
type Type int

const (
	Symbol Type = iota
	Number
	List
	Proc
	Lambda
)

type ConsList []Cons
type Cons struct {
	Type   Type
	Value  string
	Number int64
	List   ConsList
	Proc   func(ConsList) Cons
}

func (cs Cons) String() string {
	if cs.Number != 0 {
		return fmt.Sprint(cs.Number)
	}
	return cs.Value
}

func NewSymbol(value string) Cons {
	return Cons{Type: Symbol, Value: value}
}
func NewNumber(number int64) Cons {
	return Cons{Type: Number, Number: number}
}
func NewList(list []Cons) Cons {
	return Cons{Type: List, List: list}
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
