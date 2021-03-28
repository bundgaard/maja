package main

import (
	"maja/internal/repl"
	"os"
)

func main() {
	repl.Start(os.Stdin, os.Stdout)
}
