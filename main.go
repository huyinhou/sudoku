package main

import (
	"flag"
	"fmt"
)

func main() {
	file := flag.String("f", "", "data file")
	flag.Parse()

	data, err := LoadSudoku(*file)
	if err != nil {
		fmt.Print(err)
		return
	}
	s, err := NewSudoku(data)
	if err != nil {
		fmt.Print(err)
		return
	}
	s.Resolve()
	fmt.Println(s.String())
	fmt.Println(s.MaskString())
}
