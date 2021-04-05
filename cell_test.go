package main

import (
	"fmt"
	"testing"
)

func TestCellTree(t *testing.T) {

	begin := NewCell(Cell{First: "david"}, nil)
	fmt.Printf("%s\n", First(begin))

}
