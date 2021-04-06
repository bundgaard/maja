package eval

import (
	"fmt"
	"maja/pkg/ast"
)

func evalLambda(cons ast.Cons, env Env, xs ast.ConsList) (ast.Cons, error) {
	fmt.Printf("Begin Closure %#v\n", cons)
	newEnv := NewEnvironment(&env)
	for idx, symbol := range cons.List[1].List {
		fmt.Println("add to new env", symbol, xs[idx])
		newEnv.Add(symbol.Value, xs[idx])
	}
	out, err := evaluate(cons.List[2], newEnv)
	fmt.Printf("EOF Closure %#v\n", out)
	return out, err
}

func evaluate(cons ast.Cons, env Env) (ast.Cons, error) {

	switch cons.Type {
	case ast.Pair:
		if len(cons.List) < 1 {
			return cons, nil
		}
		switch cons.List[0].Value {
		case "eval":
			fmt.Printf("evaluate %#v", cons)

			return evaluate(cons.List[1], env)
		case "map":
			_ = cons.List[0] // ignore map
			fn := cons.List[1]
			args := cons.List[2:]

			result := make(ast.ConsList, 0)
			dummy, _ := evaluate(args[0], env)
			for i := range dummy.List {
				l := make(ast.ConsList, len(dummy.List)+1) // +1 for the proc
				l[0] = fn
				for j := range args {
					argv, _ := evaluate(args[j], env)
					cell := argv.List[i]
					l[j+1] = cell
				}
				cell, err := evaluate(ast.NewList(l), env)
				if err != nil {
					return ast.Cons{}, err
				}
				result = append(result, cell)
			}

			return ast.NewList(result), nil

		case "'", "quote":
			return cons.List[1], nil
		case "define":
			value, err := evaluate(cons.List[2], env)
			if err != nil {
				return ast.NewError(err), err
			}
			env.Environment[cons.List[1].Value] = value
			return env.Environment[cons.List[1].Value], nil

		case "begin":
			for i := 1; i < len(cons.List)-1; i++ {
				evaluate(cons.List[1], env)
			}
			return evaluate(cons.List[len(cons.List)-1], env)
		case "lambda", "λ":
			cons.Type = ast.Closure
			return cons, nil
		case "set!":
			/*_, err := env.Find(cons.List[1].Value)
			if err != nil {
				fmt.Println(err)
				return ast.Cons{}, err
			}
			return ast.Cons{}, nil*/
		case "if":
			// if predicate true false
			arg1, err := evaluate(cons.List[1], env)
			if err != nil {
				fmt.Println("error: if", err)
				return ast.Cons{}, err
			}
			if arg1.Value == "#t" {
				arg2, err := evaluate(cons.List[2], env)
				if err != nil {
					return ast.Cons{}, err
				}
				return arg2, nil
			} else {
				arg2, err := evaluate(cons.List[3], env)
				if err != nil {
					return ast.Cons{}, err
				}
				return arg2, nil
			}

		default:
			// found proc +/-,
			// log.Println("found proc", cons.List, cons.Value, cons.Number, cons.Type.String())
			proc, err := evaluate(cons.List[0], env)
			if err != nil {
				return ast.NewSymbol("error"), err
			}
			xs := ast.Arguments(cons)
			for i := range xs {
				value, err := evaluate(xs[i], env)
				fmt.Println("> xs ", xs, xs[i])
				if err != nil {
					return ast.NewSymbol("error"), err
				}
				xs[i] = value
				fmt.Println(">> xs ", xs, xs[i])
			}
			if proc.Type == ast.Closure {
				return evalLambda(proc, env, xs)
			} else if proc.Type == ast.Proc && len(xs) > 0 {
				return proc.Proc(xs), nil
			} else {
				fmt.Println("nothing to execute", cons)
				return ast.NewError(fmt.Errorf("error: nothing ")), nil
			}

		}
	case ast.String:
		fmt.Printf("string %+v\n", cons)

		return ast.NewString(cons.Value), nil

	case ast.Symbol:
	/*	if cons.Value[0] == '"' {
			return cons, nil
		}
		env, err := env.Find(cons.Value)
		if err != nil {
			return ast.Cons{}, fmt.Errorf("'%s' not defined", cons.Value)
		}
		fn, ok := env[cons.Value]
		if !ok {
			return ast.Cons{}, nil
		}
		return fn, nil*/
	case ast.Number:
		return cons, nil
	}
	return ast.Cons{}, nil
}

func myEval(root ast.Cons, env Env) (ast.Cons, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("error: %v %#v", err, root)
		}
	}()
	// always test for quote first

	switch root.Type {
	case ast.Number:
		return root, nil
	case ast.Symbol:
		// Look up if symbol is a registered proc/lambda
		if root.Value == "quote" {
			return root.List[0], nil
		}
		e := env.Find(root.Value)
		if e == nil {
			return ast.Cons{}, fmt.Errorf("unbound symbol '%s' %#v", root.Value, root)
		}
		return *e, nil
	case ast.String:
		return root, nil
	case ast.Closure:
		root.Type = ast.Closure
		return root, nil
	case ast.Pair:
		// (op List)
		if len(root.List) == 0 {
			return ast.Cons{}, nil // return nil
		}

		if root.List[0].Value == "quote" {
			if len(root.List[1:]) > 1 {
				return ast.NewList(root.List[1:]), nil
			}
			return root.List[1], nil

		}

		return evalPair(root, env)
	}
	return ast.Cons{}, fmt.Errorf("unhandled type %s with %#v", root.Type, root)
}

func evalPair(root ast.Cons, env Env) (ast.Cons, error) {

	// (operator args)
	operator, _ := myEval(root.List[0], env)
	args := root.List[1:]
	fmt.Printf("operator %#v\nargs %#v\n", operator, args)
	fmt.Printf("root %#v\n", root)

	values := make(ast.ConsList, len(args))
	for idx, x := range args {
		value, _ := myEval(x, env)
		values[idx] = value
	}

	switch operator.Type {
	case ast.Proc:
		result := operator.Proc(values)
		return result, nil

	default:
		fmt.Printf("unhandled operator %s %#v\n", operator.Type, operator)
	}

	return ast.Cons{}, nil
}

func Evalautor(root ast.Cons, env Env) (ast.Cons, error) {
	return myEval(root, env) // changed from evaluator
}
