package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func processSyntax(tokens []string) (Cons, []string) {

	var token string

	token, tokens = tokens[0], tokens[1:]

	/*	eat(openparen)
		processList(tokens)
		eat(closeparen)*/

	// expect pairs of parens () (()),
	if token == "(" {
		list := make(ConsList, 0)

		if len(tokens) > 1 && tokens[0] != ")" {
			for tokens[0] != ")" {
				var cons Cons
				cons, tokens = processSyntax(tokens)
				list = append(list, cons)
			}

			tokens = tokens[1:]
			return NewList(list), tokens
		}
		fmt.Println(tokens)
		return NewSymbol("missing close parent"), tokens

	} else if token == ")" {
		log.Println("unexpected ')'")
		return NewSymbol("unexpected ')'"), tokens
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

func evaluate(cons Cons, env Env) (Cons, error) {

	switch cons.Type {
	case Symbol:
		fn, ok := env.Find(cons.Value)[cons.Value]
		if !ok {
			return Cons{}, fmt.Errorf("'%s' not defined", cons.Value)
		}
		return fn, nil
	case Number:
		return cons, nil
	case Pair:
		if len(cons.List) < 1 {
			return cons, nil
		}
		switch cons.List[0].Value {
		case "define":
			value, err := evaluate(cons.List[2], env)

			if err != nil {
				return NewSymbol("error"), err
			}
			env[cons.List[1].Value] = value
		case "begin":
			for i := 1; i < len(cons.List)-1; i++ {
				evaluate(cons.List[1], env)
			}
			return evaluate(cons.List[len(cons.List)-1], env)
		case "lambda", "λ":
			cons.Type = Closure
			return cons, nil
		default:
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
			if proc.Type == Closure {
				return evaluate(proc.List[2], standardEnvironment()) // add arguments from Lambda
			} else if proc.Type == Proc && len(xs) > 0 {
				return proc.Proc(xs), nil
			} else {
				return evaluate(proc, env)
			}

		}
	}
	return Cons{}, nil
}
func main() {
	// input := `(begin (define r 10) (+ 10 r))`
	// input := `(+ 1 2 3 4 5)`

	env := standardEnvironment()
	_ = env
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

		_, err := verifyParenthesis(line)
		if err != nil {
			fmt.Println(err)
			continue
		}
		l := NewLexer(line)
		p := NewParser(l)
		program := p.Parse()
		fmt.Println(">>", program.String())
		output, err := evaluate(program, env)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(output)
		}

	}
	/*tokens := tokenize(input)
	cons, _ := processSyntax(tokens)
	evaluate(cons, env)
	*/
}
