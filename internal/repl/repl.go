package repl

import (
	"bufio"
	"fmt"
	"io"
	"maja/internal/parser"
)

const Prompt = "maja-> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	for {
		fmt.Print(Prompt)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()

		p := parser.NewParser(line)
		ast := p.SExpression()

		fmt.Printf("* %+v\n", ast)
	}
}
