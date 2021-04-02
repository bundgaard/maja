package main

type Environment map[string]Cons

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
	newList := Cons{Type: List}
	for i := 1; i < len(argv[0].List); i++ {
		newList.List = append(newList.List, argv[0].List[i])
	}
	return newList
}

func sqrt(argv ConsList) Cons {
	return NewNumber(0)
}

func standardEnvironment() Environment {
	env := make(Environment)
	env["+"] = NewProc(add)
	env["-"] = NewProc(subtract)
	env["*"] = NewProc(multiply)
	env["/"] = NewProc(divide)

	env["<"] = NewProc(lessThan)
	env[">"] = NewProc(greaterThan)
	env["="] = NewProc(equal)
	env["car"] = NewProc(car)
	env["cdr"] = NewProc(cdr)
	env["#f"] = NewSymbol("#f")
	env["#t"] = NewSymbol("#t")
	env["nil"] = NewSymbol("nil")
	env["sqrt"] = NewProc(sqrt)

	return env
}
