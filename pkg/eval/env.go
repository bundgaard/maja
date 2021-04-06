package eval

import (
	"fmt"
	"maja/pkg/ast"
	"math/big"
)

type Env struct {
	Environment map[string]ast.Cons
	Outer       *Env
}

func NewEnvironment(outerEnv *Env) Env {
	m := make(map[string]ast.Cons)
	e := Env{Environment: m, Outer: outerEnv}
	return e
}

func (e *Env) Add(symbol string, item ast.Cons) {
	e.Environment[symbol] = item
}

type errUnbound struct {
	Symbol string
}

func (e errUnbound) Error() string {
	return fmt.Sprintf("unbound %s", e.Symbol)
}

func (e *Env) Find(symbol string) *ast.Cons {
	fmt.Println("env#Find ", symbol, e)
	fn, ok := e.Environment[symbol]
	if ok {
		return &fn
	}
	if e.Outer == nil {
		return nil
	} else {
		return e.Outer.Find(symbol)
	}
}

func add(list ast.ConsList) ast.Cons {
	fmt.Printf("add %+v\n", list)
	acc := big.NewInt(0)
	for i := 0; i < len(list); i++ {
		acc.Add(acc, list[i].Number)
	}
	return ast.NewNumber(acc)
}

func subtract(list ast.ConsList) ast.Cons {
	acc := big.NewInt(0)

	for i := 0; i < len(list); i++ {
		acc.Sub(acc, list[i].Number)
	}
	return ast.NewNumber(acc)
}

func multiply(list ast.ConsList) ast.Cons {
	acc := big.NewInt(1)
	for i := 0; i < len(list); i++ {
		acc.Mul(acc, list[i].Number)
	}
	return ast.NewNumber(acc)
}
func divide(list ast.ConsList) ast.Cons {
	acc := big.NewInt(1)
	for i := 0; i < len(list); i++ {
		acc.Div(acc, list[i].Number)
	}
	return ast.NewNumber(acc)
}

func equalNumeric(argv ast.ConsList) ast.Cons {
	// = [argv . argv]
	out := ast.NewSymbol("#f")
	for i := 0; i < len(argv)-1; i++ {
		current := argv[i]
		next := argv[i+1]
		if current.Type != ast.Number || next.Type != ast.Number {
			return out
		}
		cmp := current.Number.Cmp(next.Number)
		if cmp == 0 {
			out = ast.NewSymbol("#t")
		} else {
			out = ast.NewSymbol("#f")
		}
	}
	return out
}

func lessThan(argv ast.ConsList) ast.Cons {
	out := ast.NewSymbol("#f")
	for i := 0; i < len(argv)-1; i++ {
		current := argv[i]
		next := argv[i+1]
		if current.Type != ast.Number || next.Type != ast.Number {
			return out
		}
		cmp := current.Number.Cmp(next.Number)
		if cmp == -1 {
			out = ast.NewSymbol("#t")
		} else {
			out = ast.NewSymbol("#f")
		}
	}

	return out
}
func greaterThan(argv ast.ConsList) ast.Cons {
	out := ast.NewSymbol("#f")
	for i := 0; i < len(argv)-1; i++ {
		current := argv[i]
		next := argv[i+1]
		if current.Type != ast.Number || next.Type != ast.Number {
			return out
		}
		cmp := current.Number.Cmp(next.Number)
		if cmp == 1 {
			out = ast.NewSymbol("#t")
		} else {
			out = ast.NewSymbol("#f")
		}
	}
	return out
}

func car(argv ast.ConsList) ast.Cons {
	out := argv[0]
	return out
}

func cdr(argv ast.ConsList) ast.Cons {
	newList := ast.Cons{Type: ast.Pair}
	for i := 1; i < len(argv[0].List); i++ {
		newList.List = append(newList.List, argv[0].List[i])
	}
	return newList
}

func sqrt(argv ast.ConsList) ast.Cons {
	return ast.NewNumber(argv[0].Number.Sqrt(argv[0].Number))
}

// append: take quoted arguments and combine them into a list
func appendFn(argv ast.ConsList) ast.Cons {
	l := make(ast.ConsList, 0)
	for _, entries := range argv {
		l = append(l, entries.List...)
	}
	out := ast.NewList(l)

	return out
}

func apply(args ast.ConsList) ast.Cons {
	// (apply + '(1 2 3 4))

	return args[0].Proc(args[1].List)
}

func modulo(argv ast.ConsList) ast.Cons {
	// modulo _ _ -> NewNumber
	if len(argv) == 2 {
		return ast.NewNumber(argv[0].Number.Mod(argv[0].Number, argv[1].Number))
	}
	return ast.NewSymbol("unexpected number of arguments")
}

func isNumber(argv ast.ConsList) ast.Cons {
	if car(argv).Type == ast.Number {
		return ast.NewSymbol("#t")
	}
	return ast.NewSymbol("#f")
}

func isString(args ast.ConsList) ast.Cons {

	if car(args).Type == ast.String {
		return ast.NewSymbol("#t")
	}
	return ast.NewSymbol("#f")
}
func StandardEnvironment() Env {
	env := NewEnvironment(nil)
	env.Add("+", ast.NewProc(add))

	env.Add("+", ast.NewProc(add))
	env.Add("-", ast.NewProc(subtract))
	env.Add("*", ast.NewProc(multiply))
	env.Add("/", ast.NewProc(divide))

	env.Add("modulo", ast.NewProc(modulo))

	env.Add("<", ast.NewProc(lessThan))
	env.Add(">", ast.NewProc(greaterThan))
	env.Add("=", ast.NewProc(equalNumeric))

	env.Add("car", ast.NewProc(car))
	env.Add("cdr", ast.NewProc(cdr))
	env.Add("#f", ast.NewSymbol("#f"))
	env.Add("#t", ast.NewSymbol("#t"))
	env.Add("nil", ast.NewSymbol("nil"))
	env.Add("sqrt", ast.NewProc(sqrt))
	env.Add("append", ast.NewProc(appendFn))
	env.Add("apply", ast.NewProc(apply))
	env.Add("number?", ast.NewProc(isNumber))
	env.Add("string?", ast.NewProc(isString))
	//	env.Add("map", NewProc(mapFn))

	/// disjointess
	/*
		boolean?
		symbol?
		char?
		vector?
		procedure?
		pair?
		number?
		string?
		port?

	*/
	return env
}
