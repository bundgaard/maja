package main

import "fmt"

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

func (e *Env) Find(symbol string) map[string]Cons {
	for ks := range e.Environment {
		if ks == symbol {
			return e.Environment
		}
	}
	if e.Outer == nil {
		fmt.Printf("unbound symbol '%s'", symbol)
	} else {
		return e.Outer.Find(symbol)
	}
	return nil
}
func add(list ConsList) Cons {
	acc := list[0].Number
	for i := 1; i < len(list); i++ {
		acc += list[i].Number
	}
	return NewNumber(acc)
}

func subtract(list ConsList) Cons {
	acc := list[0].Number
	for i := 1; i < len(list); i++ {
		acc -= list[i].Number
	}
	return NewNumber(acc)
}

func multiply(list ConsList) Cons {
	acc := list[0].Number
	for i := 1; i < len(list); i++ {
		acc *= list[i].Number
	}
	return NewNumber(acc)
}
func divide(list ConsList) Cons {
	acc := list[0].Number
	for i := 1; i < len(list); i++ {
		acc /= list[i].Number
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
	return argv[0].List[0]
}

func cdr(argv ConsList) Cons {
	newList := Cons{Type: Pair}
	for i := 1; i < len(argv[0].List); i++ {
		newList.List = append(newList.List, argv[0].List[i])
	}
	return newList
}

func sqrt(argv ConsList) Cons {
	return NewNumber(0)
}

func standardEnvironment() *Env {
	env := NewEnvironment(nil)
	env.Add("+", NewProc(add))

	env.Add("+", NewProc(add))
	env.Add("-", NewProc(subtract))
	env.Add("*", NewProc(multiply))
	env.Add("/", NewProc(divide))

	env.Add("<", NewProc(lessThan))
	env.Add(">", NewProc(greaterThan))
	env.Add("=", NewProc(equal))
	env.Add("car", NewProc(car))
	env.Add("cdr", NewProc(cdr))
	env.Add("#f", NewSymbol("#f"))
	env.Add("#t", NewSymbol("#t"))
	env.Add("nil", NewSymbol("nil"))
	env.Add("sqrt", NewProc(sqrt))

	return env
}
