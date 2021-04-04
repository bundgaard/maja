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
		return env[cons.Value], nil
	case Number:
		return cons, nil
	case Pair:
		if len(cons.List) < 1 {
			return cons, nil
		}
		switch cons.List[0].Value {

		/*
						List map (const SLists& argv) {
			    SList newList(SList::LIST);
			    for (int i = 0; i < argv[1].getList().size(); i++) {
			        SLists n;
			        SList args(SList::LIST);
			        n.push_back(argv[0].getProc());
			        for (int j = 1; j < argv.size(); j++) {
			            args.push(argv[j].getList()[i]);
			        }
			        n.push_back(args);
			        newList.push(apply(n));
			    }
			    return newList;
			}
		*/
		case "map":
			// (map fn lst)
			fn := cons.List[1]
			out := make(ConsList, 0)
			for i := 2; i < len(cons.List); i++ {
				xs, err := evaluate(cons.List[i], env)
				if err != nil {
					fmt.Fprintf(os.Stderr, "error: map %+v %v\n", cons.List[i], err)
					return Cons{}, err
				}
				l := make(ConsList, 0)
				l = append(l, fn)

				for j := range xs.List {

					fmt.Printf("argv[%02d][%02d] %+v\n", i, j, xs.List[i])
					l = append(l, xs.List[j])

				}
				fmt.Println("list", l)
				out = append(out, l[0])
			}

			return NewList(out), nil
		case "'", "quote":
			return cons.List[1], nil
		case "define":
			fmt.Fprintf(os.Stderr, "define %+v\n", cons)
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
				fmt.Fprintf(os.Stderr, "%+v\n", cons)
				fmt.Fprintf(os.Stderr, "xs len %d\nproc %+v\n", len(xs), proc)
				for idx, symbol := range proc.List[1].List {
					fmt.Fprintf(os.Stderr, "adding %s:%s to new environmentn\n", symbol, xs[idx])
					newEnv.Add(symbol.Value, xs[idx])
				}
				out, _ := evaluate(proc.List[2], newEnv)
				fmt.Fprintf(os.Stderr, "evaluated lambda %v\n", out)
				return out, nil
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
