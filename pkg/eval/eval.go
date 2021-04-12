package eval

import (
	"fmt"
	"maja/pkg/ast"
)

func evalLambda(cons ast.Cons, env Env, xs ast.ConsList) (ast.Cons, error) {
	fmt.Printf("Begin Closure %#v\n", cons)
	newEnv := NewEnvironment(env)
	for idx, symbol := range cons.List[1].List {
		fmt.Println("add to new env", symbol, xs[idx])
		newEnv.Add(symbol.Value, xs[idx])
	}

	fmt.Println("Middle Closure", newEnv)
	out, err := myEval(cons.List[2], newEnv)
	fmt.Printf("End Closure %#v\n", out)
	return out, err
}

func myEval(root ast.Cons, env Env) (ast.Cons, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("error: %v %#v", err, root)
		}
	}()

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
		switch root.List[0].Value {
		case "if":
			// (if test true false)
			return evalIf(root, env)
		case "lambda":
			root.Type = ast.Closure
			return root, nil
		case "define":
			// (define <name> <args>)
			return evalDefine(root, env)
		case "quote":
			if len(root.List[1:]) > 1 {
				return ast.NewList(root.List[1:]), nil
			}
			return root.List[1], nil
		default:
			return evalPair(root, env)
		}
	}
	return ast.Cons{}, fmt.Errorf("unhandled type %s with %#v", root.Type, root)
}

// (if test true false)
func evalIf(root ast.Cons, env Env) (ast.Cons, error) {
	_ = root.List[0]
	test := root.List[1]
	consequence := root.List[2]
	alternative := root.List[3]

	testEvaluated, err := myEval(test, env)
	if err != nil {
		return ast.Cons{}, err
	}

	if testEvaluated.Value == "#t" {
		outputConsequence, err := myEval(consequence, env)
		if err != nil {
			return ast.Cons{}, err
		}
		return outputConsequence, nil
	}
	outputAlternative, err := myEval(alternative, env)
	if err != nil {
		return ast.Cons{}, err
	}
	return outputAlternative, nil
}

// (define <name> <args>)
func evalDefine(root ast.Cons, env Env) (ast.Cons, error) {
	fmt.Printf("%s %s %#v\n", root.List[0], root.List[1], root.List[2])
	_ = root.List[0]
	name := root.List[1]
	args, _ := myEval(root.List[2], env)
	env.Add(name.Value, args)

	return ast.Cons{}, nil
}
func evalPair(root ast.Cons, env Env) (ast.Cons, error) {

	// (operator args)
	operator, err := myEval(root.List[0], env)
	if err != nil {
		return ast.Cons{}, err
	}
	args := root.List[1:]
	fmt.Printf("operator %#v\nargs %#v\n", operator, args)
	fmt.Printf("root %#v\n", root)

	values := make(ast.ConsList, len(args))
	for i := len(args) - 1; i > -1; i-- {
		fmt.Printf("Begin Value[%02d] %#v\n", i, args[i])
		value := args[i]
		output, _ := myEval(value, env)
		fmt.Printf("End Value[%02d] %#v\n", i, output)
	}
	for idx, x := range args {
		value, _ := myEval(x, env)
		values[idx] = value
	}

	switch operator.Type {
	case ast.Proc:
		result := operator.Proc(values)
		return result, nil
	case ast.Closure:

		return evalLambda(operator, env, values)
		// (sq 10) --> ((lambda (n) (* n n)) 10)
	default:
		fmt.Printf("unhandled operator %s %#v\n", operator.Type, operator)
	}

	return ast.Cons{}, nil
}

func Evalautor(root ast.Cons, env Env) (ast.Cons, error) {
	return myEval(root, env) // changed from evaluator
}
