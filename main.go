package main

import (
	"bufio"
	"fmt"
	"maja/pkg/ast"
	"maja/pkg/eval"
	"maja/pkg/parser"
	lexer "maja/pkg/scanner"
	"math/big"
	"os"
	"runtime"
	"strings"
)

func insertInto(env eval.Env, input ...string) {
	for _, line := range input {
		l := lexer.NewLexer(line)
		p := parser.NewParser(l)
		eval.Evalautor(p.Parse(), env)
	}
}
func main() {
	env := eval.StandardEnvironment()
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
				(= 1 n) 1 
				(* n (factorial (- n 1))))))`

	ohInput := `(define oh 
					(lambda (n) 
							(if 
								(= 0 n) 
								"foo" 
								(oh (- n 1)))))`

	_ = ohInput
	_, _, _ = cubeInput, fibInput, factorial
	//insertInto(env, cubeInput, fibInput, factorial)

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
					fn := env.Find(key)

					fn.Proc(ast.ConsList{
						{
							Type: ast.Pair,
							List: ast.ConsList{
								{Type: ast.Number, Number: big.NewInt(10)},
							},
						},
						{
							Type: ast.Pair,
							List: ast.ConsList{
								{Type: ast.Number, Number: big.NewInt(20)},
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
		l := lexer.NewLexer(line)
		p := parser.NewParser(l)
		program := p.Parse()
		output, err := eval.Evalautor(program, env)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(output)
		}

	}
}
