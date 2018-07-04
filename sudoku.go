package main

import (
	"bytes"
	"fmt"
	"strings"
)

const (
	sudoCols = 9
	sudoRows = 9
)

type BitMask uint16

func (m *BitMask) IsSet(mask BitMask) bool {
	return (*m & mask) != 0
}

func (m *BitMask) Set(mask BitMask) {
	*m |= mask
}

func (m *BitMask) Clear(mask BitMask) {
	*m &= ^mask
}

func (m *BitMask) OneBit() int {
	i, j := 1, 9
	m1, m2 := BitMask(1)<<1, BitMask(1)<<9
	for i < j {
		if !m.IsSet(m1) {
			m1 = m1 << 1
			i++
			continue
		}
		if !m.IsSet(m2) {
			m2 = m2 >> 1
			j--
			continue
		}
		break
	}
	if i == j {
		return i
	}
	return 0
}

func (m *BitMask) String() string {
	var i uint
	var buf [9]byte
	for i = 1; i < 10; i++ {
		if m.IsSet(1 << i) {
			buf[i-1] = '1'
		} else {
			buf[i-1] = '0'
		}
	}
	return string(buf[:])
}

type Sudoku struct {
	data  [sudoRows][sudoCols]int
	masks [sudoRows][sudoCols]BitMask
}

func NewSudoku(data [][]int) (*Sudoku, error) {
	s := &Sudoku{}
	for i := 0; i < sudoRows; i++ {
		for j := 0; j < sudoCols; j++ {
			num := data[i][j]
			if num < 0 || num > 9 {
				return nil, fmt.Errorf("Invalid number %d at (%d, %d)", num, i, j)
			}
			s.data[i][j] = data[i][j]
		}
	}
	err := s.init()
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Sudoku) init() error {
	for i := 0; i < sudoRows; i++ {
		for j := 0; j < sudoCols; j++ {
			s.masks[i][j] = BitMask(1022)
		}
	}
	for i := 0; i < sudoRows; i++ {
		for j := 0; j < sudoCols; j++ {
			if s.data[i][j] > 0 {
				s.masks[i][j] = 0
				s.setMask(i, j, uint(s.data[i][j]), false)
			}
		}
	}
	return nil
}

func (s *Sudoku) setMask(i, j int, bit uint, set bool) {
	fmt.Printf("setMask(%d, %d, %d, %v)\n", i, j, bit, set)
	mask := BitMask(1) << bit
	old := s.masks[i][j]
	// 行
	s.setMaskBlock(i, i+1, 0, 9, mask, set)
	// 列
	s.setMaskBlock(0, 9, j, j+1, mask, set)
	// 块
	s.setMaskBlock(i/3*3, (i+3)/3*3, j/3*3, (j+3)/3*3, mask, set)
	s.masks[i][j] = old
}

func (s *Sudoku) setMaskBlock(rb, re, cb, ce int, mask BitMask, set bool) {
	for r := rb; r < re; r++ {
		for c := cb; c < ce; c++ {
			m := s.masks[r][c]
			if set {
				m.Set(mask)
			} else {
				m.Clear(mask)
			}
			s.masks[r][c] = m
		}
	}
}

func (s *Sudoku) String() string {
	buf := bytes.NewBuffer(nil)
	for idx, line := range s.data {
		if idx%3 == 0 {
			buf.WriteString(strings.Repeat("-", 25))
			buf.WriteString("\n")
		}
		buf.WriteString(fmt.Sprintf(
			"| %d %d %d | %d %d %d | %d %d %d |\n",
			line[0], line[1], line[2],
			line[3], line[4], line[5],
			line[6], line[7], line[8]))
	}
	buf.WriteString(strings.Repeat("-", 25))
	buf.WriteString("\n")
	return buf.String()
}

func (s *Sudoku) MaskString() string {
	buf := bytes.NewBuffer(nil)
	for idx, line := range s.masks {
		if idx%3 == 0 {
			buf.WriteString(strings.Repeat("-", 97))
			buf.WriteString("\n")
		}
		buf.WriteString(fmt.Sprintf("| %s %s %s | %s %s %s | %s %s %s |\n",
			line[0].String(), line[1].String(), line[2].String(),
			line[3].String(), line[4].String(), line[5].String(),
			line[6].String(), line[7].String(), line[8].String()))
	}
	buf.WriteString(strings.Repeat("-", 97))
	buf.WriteString("\n")
	return buf.String()
}

// 找到最近一个可以确定的数字
func (s *Sudoku) fillOne(row, col *int) bool {
	r, c := *row, *col
	for i := 0; i < sudoRows; i, r = i+1, (r+1)%9 {
		for j := 0; j < sudoCols; j, c = j+1, (c+1)%9 {
			if s.data[r][c] > 0 {
				continue
			}
			pos := s.masks[r][c].OneBit()
			if pos > 0 {
				s.data[r][c] = pos
				s.setMask(r, c, uint(pos), false)
				*row = r
				*col = c
				return true
			}
		}
	}
	return false
}

func (s *Sudoku) Resolve() error {
	row, col := 0, 0
	for {
		if !s.fillOne(&row, &col) {
			break
		}
	}
	return nil
}
