package main

import (
	"fmt"
	"math/big"
)

type Env struct {
	Environment map[string]Expr
	Outer       *Env
}

func NewEnvironment(outerEnv *Env) Env {
	m := make(map[string]Expr)
	e := Env{Environment: m, Outer: outerEnv}
	return e
}

func (e *Env) Add(symbol string, item Expr) {
	e.Environment[symbol] = item
}

type errUnbound struct {
	Symbol string
}

func (e errUnbound) Error() string {
	return fmt.Sprintf("unbound %s", e.Symbol)
}

func (e *Env) Find(symbol string) (map[string]Expr, error) {
	fmt.Println("env#Find ", symbol, e)
	for ks := range e.Environment {
		if ks == symbol {
			return e.Environment, nil
		}
	}
	if e.Outer == nil {
		return nil, errUnbound{Symbol: symbol}
	} else {
		return e.Outer.Find(symbol)
	}
}
func add(list ConsList) Expr {
	acc := list[0].Number
	for i := 1; i < len(list); i++ {
		acc.Add(acc, list[i].Number)
	}
	return NewNumber(acc)
}

func subtract(list ConsList) Expr {
	acc := list[0].Number
	for i := 1; i < len(list); i++ {
		acc.Sub(acc, list[i].Number)
	}
	return NewNumber(acc)
}

func multiply(list ConsList) Expr {
	acc := big.NewInt(1)
	for i := 0; i < len(list); i++ {
		acc.Mul(acc, list[i].Number)
	}
	return NewNumber(acc)
}
func divide(list ConsList) Expr {
	acc := list[0].Number
	for i := 1; i < len(list); i++ {
		acc.Div(acc, list[i].Number)
	}
	return NewNumber(acc)
}

func equalNumeric(argv ConsList) Expr {
	// = [argv . argv]
	out := NewSymbol("#f")
	for i := 0; i < len(argv)-1; i++ {
		current := argv[i]
		next := argv[i+1]
		if current.Type != Number || next.Type != Number {
			return out
		}
		cmp := current.Number.Cmp(next.Number)
		if cmp == 0 {
			out = NewSymbol("#t")
		} else {
			out = NewSymbol("#f")
		}
	}
	return out
}

func lessThan(argv ConsList) Expr {
	out := NewSymbol("#f")
	for i := 0; i < len(argv)-1; i++ {
		current := argv[i]
		next := argv[i+1]
		if current.Type != Number || next.Type != Number {
			return out
		}
		cmp := current.Number.Cmp(next.Number)
		if cmp == -1 {
			out = NewSymbol("#t")
		} else {
			out = NewSymbol("#f")
		}
	}

	return out
}
func greaterThan(argv ConsList) Expr {
	out := NewSymbol("#f")
	for i := 0; i < len(argv)-1; i++ {
		current := argv[i]
		next := argv[i+1]
		if current.Type != Number || next.Type != Number {
			return out
		}
		cmp := current.Number.Cmp(next.Number)
		if cmp == 1 {
			out = NewSymbol("#t")
		} else {
			out = NewSymbol("#f")
		}
	}
	return out
}

func car(argv ConsList) Expr {
	out := argv[0]
	return out
}

func cdr(argv ConsList) Expr {
	newList := Expr{Type: Pair}
	for i := 1; i < len(argv[0].List); i++ {
		newList.List = append(newList.List, argv[0].List[i])
	}
	return newList
}

func sqrt(argv ConsList) Expr {
	return NewNumber(argv[0].Number.Sqrt(argv[0].Number))
}

// append: take quoted arguments and combine them into a list
func appendFn(argv ConsList) Expr {
	l := make(ConsList, 0)
	for _, entries := range argv {
		l = append(l, entries.List...)
	}
	out := NewList(l)

	return out
}

func apply(args ConsList) Expr {
	// (apply + '(1 2 3 4))

	return args[0].Proc(args[1].List)
}

func modulo(argv ConsList) Expr {
	// modulo _ _ -> NewNumber
	if len(argv) == 2 {
		return NewNumber(argv[0].Number.Mod(argv[0].Number, argv[1].Number))
	}
	return NewSymbol("unexpected number of arguments")
}

func isNumber(argv ConsList) Expr {
	if car(argv).Type == Number {
		return NewSymbol("#t")
	}
	return NewSymbol("#f")
}

func isString(args ConsList) Expr {

	if car(args).Type == String {
		return NewSymbol("#t")
	}
	return NewSymbol("#f")
}
func standardEnvironment() Env {
	env := NewEnvironment(nil)
	env.Add("+", NewProc(add))

	env.Add("+", NewProc(add))
	env.Add("-", NewProc(subtract))
	env.Add("*", NewProc(multiply))
	env.Add("/", NewProc(divide))

	env.Add("modulo", NewProc(modulo))

	env.Add("<", NewProc(lessThan))
	env.Add(">", NewProc(greaterThan))
	env.Add("=", NewProc(equalNumeric))

	env.Add("car", NewProc(car))
	env.Add("cdr", NewProc(cdr))
	env.Add("#f", NewSymbol("#f"))
	env.Add("#t", NewSymbol("#t"))
	env.Add("nil", NewSymbol("nil"))
	env.Add("sqrt", NewProc(sqrt))
	env.Add("append", NewProc(appendFn))
	env.Add("apply", NewProc(apply))
	env.Add("number?", NewProc(isNumber))
	env.Add("string?", NewProc(isString))
	//	env.Add("map", NewProc(mapFn))

	return env
}
