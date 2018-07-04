package main

import (
	"bufio"
	"fmt"
	"os"
	"unicode"
)

func LoadSudoku(file string) ([][]int, error) {
	var reader *bufio.Reader
	if file == "" {
		reader = bufio.NewReader(os.Stdin)
	} else {
		f, err := os.Open(file)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		reader = bufio.NewReader(f)
	}
	return LoadFromReader(reader)
}

func LoadFromReader(reader *bufio.Reader) ([][]int, error) {
	var board [][]int
	for i := 0; i < 9; i++ {
		line, _, err := reader.ReadLine()
		if err != nil {
			return nil, err
		}
		if len(line) < 9 {
			return nil, fmt.Errorf("len(%s) < 9", string(line))
		}

		var data []int
		for j := 0; j < 9; j++ {
			if !unicode.IsDigit(rune(line[j])) {
				return nil, fmt.Errorf("%v is not digit", line[j])
			}
			data = append(data, int(line[j]-'0'))
		}
		board = append(board, data)
	}
	return board, nil
}
