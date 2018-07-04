package main

import (
	"bufio"
	"bytes"
	"fmt"
	"testing"
)

func Test_Case1(t *testing.T) {
	mask := BitMask(22)
	fmt.Println(mask.String())

	mask = BitMask(16)
	if mask.OneBit() != 4 {
		t.Error(mask.OneBit())
		return
	}
	mask = BitMask(64)
	if mask.OneBit() != 6 {
		t.Error(mask.OneBit())
		return
	}
	mask = BitMask(123)
	if mask.OneBit() != 0 {
		t.Error("!= 0")
		return
	}
}

func Test_Case2(t *testing.T) {
	buf := bytes.NewBufferString(`092130050
800600309
100097080
750000100
203060408
009000072
040250001
506003007
080074620`)
	data, err := LoadFromReader(bufio.NewReader(buf))
	if err != nil {
		t.Fail()
	}
	s, err := NewSudoku(data)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	fmt.Println(s.String())
	fmt.Println(s.MaskString())
	s.Resolve()
	fmt.Println(s.String())
	fmt.Println(s.MaskString())
}
