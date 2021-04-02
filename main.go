package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func processSyntax(tokens []string) (Cons, []string) {

	var token string

	token, tokens = tokens[0], tokens[1:]

	if token == "(" {

		list := make(ConsList, 0)
		for tokens[0] != ")" {
			var cons Cons
			cons, tokens = processSyntax(tokens)
			list = append(list, cons)
		}
		_, tokens = tokens[0], tokens[1:]
		return NewList(list), tokens
	} else {
		return atomic(token), tokens
	}

}

func atomic(token string) Cons {
	n, err := strconv.ParseInt(token, 0, 64)
	if err != nil {
		return NewSymbol(token)
	} else {
		return NewNumber(n)
	}
}

func evaluate(cons Cons, env Environment) (Cons, error) {

	if cons.Type == Symbol {
		fn, ok := env[cons.Value]
		if !ok {
			return Cons{}, fmt.Errorf("'%s' not defined", cons.Value)
		}
		return fn, nil
	} else if cons.Type == Number {
		return cons, nil
	} else if len(cons.List) == 0 {
		return cons, nil
	} else if cons.List[0].Value == "define" {
		// log.Println("found define", cons.List[1].Value, cons.List[2].Number)
		value, err := evaluate(cons.List[2], env)

		if err != nil {
			return NewSymbol("error"), err
		}
		env[cons.List[1].Value] = value
	} else if cons.List[0].Value == "begin" {
		for i := 1; i < len(cons.List)-1; i++ {
			evaluate(cons.List[1], env)
		}
		return evaluate(cons.List[len(cons.List)-1], env)
	} else {
		// found proc +/-,
		// log.Println("found proc", cons.List, cons.Value, cons.Number, cons.Type.String())
		proc, err := evaluate(cons.List[0], env)
		if err != nil {
			return NewSymbol("error"), err
		}
		// log.Println("proc", proc.Type, proc)
		xs := args(cons)
		for i := range xs {
			value, err := evaluate(xs[i], env)
			if err != nil {
				return NewSymbol("error"), err
			}
			xs[i] = value
		}

		if proc.Type == Proc && len(xs) > 0 {
			return proc.Proc(xs), nil
		}

	}
	return Cons{}, nil
}
func main() {
	// input := `(begin (define r 10) (+ 10 r))`
	// input := `(+ 1 2 3 4 5)`
	env := standardEnvironment()
	input := `(begin (define r (+ 1 2 3 4 5)) (+ 10 r))`
	tokens := tokenize(input)
	cons, _ := processSyntax(tokens)
	evaluate(cons, env)
	scanner := bufio.NewScanner(os.Stdin)
	Prompt := "λ -> "

	for {

		fmt.Print(Prompt)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		if line == "" {
			continue
		}

		tokens := tokenize(line)
		cons, _ := processSyntax(tokens)
		output, err := evaluate(cons, env)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(output)
		}

	}
}
