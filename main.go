package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
	"strings"
)

func evaluate(cons Cons, env *Env) (Cons, error) {

	switch cons.Type {
	case Symbol:
		if cons.Value[0] == '"' {
			return NewSymbol(cons.Value[1 : len(cons.Value)-1]), nil
		}
		env, err := env.Find(cons.Value)
		if err != nil {
			return Cons{}, fmt.Errorf("'%s' not defined", cons.Value)
		}
		fn, ok := env[cons.Value]
		if !ok {
			return Cons{}, fmt.Errorf("'%s' not found in environment", cons)
		}
		return fn, nil
	case Number:
		return cons, nil
	case Pair:
		if len(cons.List) < 1 {
			return cons, nil
		}
		switch cons.List[0].Value {
		case "map":
			_ = cons.List[0] // ignore map
			fn := cons.List[1]
			args := cons.List[2:]

			result := make(ConsList, 0)
			dummy, _ := evaluate(args[0], env)
			for i := range dummy.List {
				l := make(ConsList, len(dummy.List)+1) // +1 for the proc
				l[0] = fn
				for j := range args {
					argv, _ := evaluate(args[j], env)
					cell := argv.List[i]
					l[j+1] = cell
				}
				cell, err := evaluate(NewList(l), env)
				if err != nil {
					return Cons{}, err
				}
				result = append(result, cell)
			}

			return NewList(result), nil

		case "'", "quote":
			return cons.List[1], nil
		case "define":
			value, err := evaluate(cons.List[2], env)
			if err != nil {
				return NewSymbol("error"), err
			}
			env.Environment[cons.List[1].Value] = value
			return env.Environment[cons.List[1].Value], nil

		case "begin":
			for i := 1; i < len(cons.List)-1; i++ {
				evaluate(cons.List[1], env)
			}
			return evaluate(cons.List[len(cons.List)-1], env)
		case "lambda", "λ":
			cons.Type = Closure
			return cons, nil
		case "set!":
			_, err := env.Find(cons.List[1].Value)
			if err != nil {
				fmt.Println(err)
				return Cons{}, err
			}
			return Cons{}, nil
		case "if":
			/* return evaluate(s.getList()[1],env).val()=="#t" ?
			evaluate(s.getList()[2],env) :
				(s.getList()[3].val() == "else" ?
					evaluate(s.getList()[4],env) :
						SList());
			*/
			arg1, err := evaluate(cons.List[1], env)
			if err != nil {
				fmt.Println("error: if", err)
				return Cons{}, err
			}
			if arg1.Value == "#t" {
				arg2, err := evaluate(cons.List[2], env)
				if err != nil {
					fmt.Println("error: if true", err)

				}
				fmt.Println("arg2", arg2)

			} else {
				return Cons{}, nil
			}

		default:
			// found proc +/-,
			// log.Println("found proc", cons.List, cons.Value, cons.Number, cons.Type.String())
			proc, err := evaluate(cons.List[0], env)
			if err != nil {
				return NewSymbol("error"), err
			}
			xs := args(cons)
			for i := range xs {
				value, err := evaluate(xs[i], env)
				if err != nil {
					return NewSymbol("error"), err
				}
				xs[i] = value
			}
			if proc.Type == Closure {
				newEnv := NewEnvironment(env)
				for idx, symbol := range proc.List[1].List {
					newEnv.Add(symbol.Value, xs[idx])
				}
				out, err := evaluate(proc.List[2], newEnv)
				return out, err
			} else if proc.Type == Proc && len(xs) > 0 {
				return proc.Proc(xs), nil
			} else {
				fmt.Println("nothing to execute", cons)
			}

		}
	}
	return Cons{}, nil
}
func main() {
	env := standardEnvironment()
	scanner := bufio.NewScanner(os.Stdin)
	Prompt := "λ -> "

	cubeInput := `(define cube (lambda (x) (* x x x)))`
	l := NewLexer(cubeInput)
	p := NewParser(l)
	evaluate(p.Parse(), env)
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
		fields := strings.Fields(line)
		if fields[0][0] == '?' {
			switch fields[0] {
			case "?env":

				if len(fields) == 1 {
					fmt.Println("Environment")
					fmt.Println(env.Environment)

				} else {
					key := fields[1]
					env, err := env.Find(key)
					if err != nil {
						fmt.Fprintf(os.Stderr, "%s not found in environment.", key)
					}
					fn, ok := env[key]
					if !ok {
						fmt.Fprintf(os.Stderr, "%s not found in environment", key)
					}
					fn.Proc(ConsList{
						{
							Type: Pair,
							List: ConsList{
								{Type: Number, Number: big.NewInt(10)},
							},
						},
						{
							Type: Pair,
							List: ConsList{
								{Type: Number, Number: big.NewInt(20)},
							},
						},
					})
				}
				continue
			case "?exit":
				fmt.Fprintf(os.Stderr, "Goodbye\n")
				os.Exit(0)
			}
		}

		_, err := verifyParenthesis(line, nil)
		if err != nil {
			fmt.Println(err)
			continue
		}
		l := NewLexer(line)
		p := NewParser(l)
		program := p.Parse()
		output, err := evaluate(program, env)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(output)
		}

	}
}
