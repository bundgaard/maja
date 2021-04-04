package main

import (
	"bufio"
	"fmt"
	"os"
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
		return env[cons.Value], nil
	case Number:
		return cons, nil
	case Pair:
		if len(cons.List) < 1 {
			return cons, nil
		}
		switch cons.List[0].Value {
		case "quote", "'":
			return cons.List[1], nil
		case "define":
			value, err := evaluate(cons.List[2], env)

			if err != nil {
				return NewSymbol("error"), err
			}
			env.Environment[cons.List[1].Value] = value
		case "begin":
			for i := 1; i < len(cons.List)-1; i++ {
				evaluate(cons.List[1], env)
			}
			return evaluate(cons.List[len(cons.List)-1], env)
		case "lambda", "λ":
			cons.Type = Closure
			return cons, nil
		// case "if":
		// evaluate(cons.List[1], env) // return #T or #F
		case "set!":
			fmt.Printf("set %+v %+v\n", cons, cons.List[1])
			temp, err := env.Find(cons.List[1].Value)
			if err != nil {
				fmt.Println(err)
				return Cons{}, nil
			}
			fmt.Printf("temp %+v\n ", temp)
			/*temp := env.Find(cons.List[1].Value)
			temp[cons.List[1].Value], _ = evaluate(cons.List[2], env)*/
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
				newEnv := NewEnvironment(env)
				newEnv.Outer = env
				for idx, symbol := range proc.List[1].List {
					newEnv.Add(symbol.Value, xs[idx])
				}
				// newEnv.Add(proc.List[1].List, xs)
				return evaluate(proc.List[2], newEnv) // add arguments from Lambda
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
