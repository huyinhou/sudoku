package main

import (
	"bufio"
	"bytes"
	"fmt"
	"testing"
)

func TestLoadSudoku(t *testing.T) {
	// 不足9行
	buf := bytes.NewBufferString(`012345678`)
	_, err := LoadFromReader(bufio.NewReader(buf))
	if err == nil {
		t.Fail()
	}
	t.Log(err)

	// 非法字符
	buf = bytes.NewBufferString(`0a2345678`)
	_, err = LoadFromReader(bufio.NewReader(buf))
	if err == nil {
		t.Fail()
	}
	t.Log(err)

	// 不足9列
	buf = bytes.NewBufferString(`01234567`)
	_, err = LoadFromReader(bufio.NewReader(buf))
	if err == nil {
		t.Fail()
	}
	t.Log(err)

	// OK
	buf = bytes.NewBufferString(`012345678
012345678
012345678
012345678
012345678d
012345678
012345678
012345678
012345678`)
	data, err := LoadFromReader(bufio.NewReader(buf))
	if err != nil {
		t.Fail()
	}
	for _, line := range data {
		fmt.Println(line)
	}
}
