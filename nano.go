package main

import (
	"fmt"
	"strings"
)

type Nano struct {
}
type NanoFn func(string, error) (string, error)

func NanoCompose(fns ...NanoFn) func(string) (string, error) {
	return func(data string) (string, error) {
		var result string
		errs := make([]error, 0)
		_ = errs
		var err error
		for _, fn := range fns {
			result, err = fn(data, err)
			if err != nil {
				errs = append(errs, err)
			}
		}
		if len(errs) > 0 {
			return "", fmt.Errorf("a lot of errors")
		}
		return result, nil
	}

}

func verifyParenthesis(data string, wrappedErr error) (string, error) {
	opened := 0
	closed := 0

	for i := 0; i < len(data); i++ {
		if data[i] == '(' {
			opened++
		} else if data[i] == ')' {
			closed++
		}
	}

	if opened != closed {
		return "", fmt.Errorf("failed to find equal open and closed parenthesis")
	}
	return data, nil

}

/*
translateQuote

(+ '(1 2 3 4)) -> (+ (quote(1 2 3 4)))
*/
func translateQuote(data string, wrappedErr error) (string, error) {
	indexOfQuote := strings.Index(data, "'")
	fmt.Println("translateQuote", indexOfQuote)
	return data, nil
}
