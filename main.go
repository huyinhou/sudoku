package main

import (
	"flag"
	"fmt"

	"github.com/golang/glog"
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
	glog.V(6).Info(s.MaskString())
}
