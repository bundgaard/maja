package main

import (
	"fmt"
	"math/big"
)

type Env struct {
	Environment map[string]Cons
	Outer       *Env
}

func NewEnvironment(outerEnv *Env) *Env {
	m := make(map[string]Cons)
	e := Env{Environment: m, Outer: outerEnv}
	return &e
}

func (e *Env) Add(symbol string, item Cons) {
	e.Environment[symbol] = item
}

type errUnbound struct {
	Symbol string
}

func (e errUnbound) Error() string {
	return fmt.Sprintf("unbound %s", e.Symbol)
}

func (e *Env) Find(symbol string) (map[string]Cons, error) {
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
func add(list ConsList) Cons {
	acc := list[0].Number
	for i := 1; i < len(list); i++ {
		acc.Add(acc, list[i].Number)
	}
	return NewNumber(acc)
}

func subtract(list ConsList) Cons {
	acc := list[0].Number
	for i := 1; i < len(list); i++ {
		acc.Sub(acc, list[i].Number)
	}
	return NewNumber(acc)
}

func multiply(list ConsList) Cons {
	fmt.Println("multiply#start", list)
	acc := big.NewInt(1)
	for i := 0; i < len(list); i++ {
		acc.Mul(acc, list[i].Number)
	}
	fmt.Println("multiply#out", acc)
	return NewNumber(acc)
}
func divide(list ConsList) Cons {
	acc := list[0].Number
	for i := 1; i < len(list); i++ {
		acc.Div(acc, list[i].Number)
	}
	return NewNumber(acc)
}

func equal(list ConsList) Cons {
	return NewSymbol("#f")
}

func lessThan(list ConsList) Cons {
	return NewSymbol("#f")
}
func greaterThan(list ConsList) Cons {
	return NewSymbol("#f")
}

func car(argv ConsList) Cons {
	out := argv[0].List[0]
	return out
}

func cdr(argv ConsList) Cons {
	newList := Cons{Type: Pair}
	for i := 1; i < len(argv[0].List); i++ {
		newList.List = append(newList.List, argv[0].List[i])
	}
	return newList
}

func sqrt(argv ConsList) Cons {
	return NewNumber(argv[0].Number.Sqrt(argv[0].Number))
}

// append: take quoted arguments and combine them into a list
func appendFn(argv ConsList) Cons {
	fmt.Println("append ", argv)

	l := make(ConsList, 0)
	for _, entries := range argv {
		l = append(l, entries.List...)
	}
	out := NewList(l)
	fmt.Println("append out", out)
	return out
}

func apply(args ConsList) Cons {
	// (apply + '(1 2 3 4))
	fmt.Printf("apply %#v\n", args)
	return args[0].Proc(args[1].List)
}

func mapFn(argv ConsList) Cons {
	fmt.Printf("map %+v\n", argv)
	lst := make(ConsList, 0)
	lst = append(lst, argv[0])
	for _, x := range argv[1].List {

		fmt.Printf("%+v\n", x)
		lst = append(lst, x)
	}

	fmt.Printf("map %+v\n", lst)
	return apply(lst)

}
func modulo(argv ConsList) Cons {
	// modulo _ _ -> NewNumber
	if len(argv) == 2 {
		return NewNumber(argv[0].Number.Mod(argv[0].Number, argv[1].Number))
	}
	return NewSymbol("unexpected number of arguments")
}
func standardEnvironment() *Env {
	env := NewEnvironment(nil)
	env.Add("+", NewProc(add))

	env.Add("+", NewProc(add))
	env.Add("-", NewProc(subtract))
	env.Add("*", NewProc(multiply))
	env.Add("/", NewProc(divide))

	env.Add("modulo", NewProc(modulo))

	env.Add("<", NewProc(lessThan))
	env.Add(">", NewProc(greaterThan))
	env.Add("=", NewProc(equal))
	env.Add("car", NewProc(car))
	env.Add("cdr", NewProc(cdr))
	env.Add("#f", NewSymbol("#f"))
	env.Add("#t", NewSymbol("#t"))
	env.Add("nil", NewSymbol("nil"))
	env.Add("sqrt", NewProc(sqrt))
	env.Add("append", NewProc(appendFn))
	env.Add("apply", NewProc(apply))
	// env.Add("map", NewProc(mapFn))

	return env
}
