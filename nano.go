package main

import "fmt"

type Nano struct {
}

func verifyParenthesis(data string) (string, error) {
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
