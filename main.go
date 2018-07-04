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
	for _, line := range data {
		fmt.Println(line)
	}
}
