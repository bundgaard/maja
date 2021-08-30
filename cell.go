package main

import (
	"fmt"
)

type Cell struct {
	Tag   Type
	First interface{}
	Rest  interface{}
}

func NewCell(first interface{}, rest interface{}) *Cell {
	return &Cell{First: first, Rest: rest}
}

func First(cell interface{}) *Cell {
	switch x := cell.(type) {
	case *Cell:
		out, ok := x.First.(Cell)
		if !ok {
			fmt.Printf("cannot change type")
		}
		return &out
	default:
		fmt.Printf("%T", x)
		return nil
	}
}

func (c Cell) String() string {
	return ""
}
