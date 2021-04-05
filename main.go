package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
	"runtime"
	"strings"
)

func evaluate(cons Cons, env Env) (Cons, error) {

	switch cons.Type {
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
			// if predicate true false
			arg1, err := evaluate(cons.List[1], env)
			if err != nil {
				fmt.Println("error: if", err)
				return Cons{}, err
			}
			if arg1.Value == "#t" {
				arg2, err := evaluate(cons.List[2], env)
				if err != nil {
					return Cons{}, err

				}
				return arg2, nil
			} else {
				arg2, err := evaluate(cons.List[3], env)
				if err != nil {
					return Cons{}, err
				}
				return arg2, nil
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
				fmt.Println("> xs ", xs, xs[i])
				if err != nil {
					return NewSymbol("error"), err
				}
				xs[i] = value
				fmt.Println(">> xs ", xs, xs[i])
			}
			if proc.Type == Closure {
				fmt.Printf("Closure %+v\n", proc)
				newEnv := NewEnvironment(&env)
				for idx, symbol := range proc.List[1].List {
					fmt.Println("add to new env", symbol, xs[idx])
					newEnv.Add(symbol.Value, xs[idx])
				}
				out, err := evaluate(proc.List[2], newEnv)
				fmt.Println("Closure", out)
				return out, err
			} else if proc.Type == Proc && len(xs) > 0 {
				return proc.Proc(xs), nil
			} else {
				fmt.Println("nothing to execute", cons)
				return Cons{}, fmt.Errorf("nothing to execute")
			}

		}
	case String:
		fmt.Printf("string %+v\n", cons)

		return NewString(cons.Value), nil

	case Symbol:
		if cons.Value[0] == '"' {
			return cons, nil
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
	}
	return Cons{}, nil
}

func insertInto(env Env, input ...string) {
	for _, line := range input {
		l := NewLexer(line)
		p := NewParser(l)
		evaluate(p.Parse(), env)
	}
}
func main() {
	env := standardEnvironment()
	scanner := bufio.NewScanner(os.Stdin)
	Prompt := "λ -> "

	cubeInput := `(define cube (lambda 
		(x) 
		(* x x x)))`
	fibInput := `(define fib (lambda 
		(n) (if 
			(< n 2) 
			n
			(+ (fib (- n 1))
				(fib (- n 2)) 
			))))`
	factorial := `(define factorial (lambda 
			(n)
			(if 
				(= 0 n) 1 
				(* n (factorial (- n 1))))))`
	insertInto(env, cubeInput, fibInput, factorial)

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
			case "?mem":
				ms := runtime.MemStats{}
				runtime.ReadMemStats(&ms)
				fmt.Println("Alloc", ms.Alloc/1024/1024, "MiB")
				fmt.Println("Total alloc", ms.TotalAlloc/1024/1024, "MiB")

				fmt.Println("Sys", ms.Sys/1024/1024, "MiB")
				fmt.Println("NumGC", ms.NumGC)
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
