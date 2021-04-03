package main

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

func main() {

	input := "λ hello World 💛"

	fmt.Println(len(input))

	for len(input) > 0 {
		r, width := utf8.DecodeRuneInString(input)
		fmt.Printf("0x%x %c %d\n", r, r, width)
		fmt.Println("IsGraphic ", unicode.IsGraphic(r))
		fmt.Println("IsSymbol", unicode.IsSymbol(r))
		input = input[width:]
	}
}
