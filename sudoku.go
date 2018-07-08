package main

import (
	"bytes"
	"errors"
	"fmt"
	"math/bits"
	"strings"

	"github.com/golang/glog"
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

func (m *BitMask) OneBit() bool {
	return bits.OnesCount16(uint16(*m)) == 1
}

func (m *BitMask) LowBit() int {
	for i, mask := 1, BitMask(1)<<1; i <= 9; i, mask = i+1, mask<<1 {
		if m.IsSet(mask) {
			return i
		}
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

func getBlockByPos(i, j int) int {
	return i/3*3 + j/3
}

type Sudoku struct {
	data   [sudoRows][sudoCols]int
	masks  [sudoRows][sudoCols]BitMask // [row][col]BitMask	每个单元格中可以填入的数字
	rmask  [sudoRows]BitMask           // 每一行已经使用的数字
	cmask  [sudoCols]BitMask           // 每一列已经使用的数字
	bmask  [9]BitMask                  // 每个block已经使用的数字
	remain int
}

func NewSudoku(data [][]int) (*Sudoku, error) {
	s := &Sudoku{
		remain: 81,
	}
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
				s.remain--
				s.masks[i][j] = 0
				s.setMask(i, j, uint(s.data[i][j]), false)
			}
		}
	}
	return nil
}

func (s *Sudoku) setMask(i, j int, bit uint, set bool) {
	glog.V(6).Infof("setMask(%d, %d, %d, %v)\n", i, j, bit, set)
	mask := BitMask(1) << bit
	old := s.masks[i][j]
	if s.rmask[i].IsSet(mask) {
		panic(fmt.Sprintf("invalid row %d mask bit %d at (%d, %d) = %s",
			i, bit, i, j, s.rmask[i].String()))
	}
	s.rmask[i].Set(mask)
	if s.cmask[j].IsSet(mask) {
		panic(fmt.Sprintf("invalid col %d mask bit %d at (%d, %d) = %s",
			j, bit, i, j, s.cmask[j].String()))
	}
	s.cmask[j].Set(mask)
	block := getBlockByPos(i, j)
	if s.bmask[block].IsSet(mask) {
		panic(fmt.Sprintf("invalid block %d mask bit %d at (%d, %d) = %s",
			block, bit, i, j, s.bmask[block].String()))
	}
	s.bmask[block].Set(mask)

	// 行
	s.setMaskBlock(i, i+1, 0, 9, mask, set)
	// 列
	s.setMaskBlock(0, 9, j, j+1, mask, set)
	// 块
	s.setMaskBlock(i/3*3, (i+3)/3*3, j/3*3, (j+3)/3*3, mask, set)
	s.masks[i][j] = old
}

func (s *Sudoku) newTraverseResolver() *traverseResolver {
	tr := &traverseResolver{}
	for r := 0; r < 9; r++ {
		tr.bmask[r] = s.bmask[r]
		tr.cmask[r] = s.cmask[r]
		for n, m := 1, BitMask(1)<<1; n <= 9; n, m = n+1, m<<1 {
			if !s.rmask[r].IsSet(m) {
				tr.numbers[r] = append(tr.numbers[r], n)
			}
		}
		for c := 0; c < 9; c++ {
			if s.data[r][c] == 0 {
				tr.columns[r] = append(tr.columns[r], c)
			}
		}
		if len(tr.numbers[r]) != len(tr.columns[r]) {
			panic(fmt.Sprintf("row %d state error, numbers: %v, columns: %v",
				r, tr.numbers[r], tr.columns[r]))
		}
	}
	return tr
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

func (s *Sudoku) fillBlank(r, c int, num int) {
	s.remain--
	s.data[r][c] = num
	s.setMask(r, c, uint(num), false)
}

// 找到最近一个可以确定的数字
func (s *Sudoku) fillOne(row, col *int) bool {
	r, c := *row, *col
	for i := 0; i < sudoRows; i, r = i+1, (r+1)%9 {
		for j := 0; j < sudoCols; j, c = j+1, (c+1)%9 {
			if s.data[r][c] > 0 {
				continue
			}
			mask := s.masks[r][c]
			if mask.OneBit() {
				num := mask.LowBit()
				if num == 0 {
					panic(fmt.Sprintf("invalid mask %s at (%d, %d)", mask.String(), r, c))
				}
				s.fillBlank(r, c, num)
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
	for s.fillOne(&row, &col) {
	}
	if s.remain == 0 {
		return nil
	}
	tr := s.newTraverseResolver()
	if tr.Resolve() == nil {
		for _, r := range tr.Result() {
			s.fillBlank(r.Row, r.Col, r.Num)
		}
		return nil
	}
	return errors.New("no solution")
}

type resolveResult struct {
	Row int
	Col int
	Num int
}

// 遍历求解
type traverseResolver struct {
	numbers [sudoRows][]int
	columns [sudoRows][]int
	cmask   [sudoCols]BitMask // 每一列已经使用的数字
	bmask   [9]BitMask        // 每个block已经使用的数字
}

func (tr *traverseResolver) traverse(r int) bool {
	if r >= 9 {
		return true
	}
	n := len(tr.columns[r]) // 这一行已经填满了
	if n == 0 {
		return tr.traverse(r + 1)
	}
	return tr.traverseRow(r, 0, len(tr.columns[r]))
}

func (tr *traverseResolver) traverseRow(r, b, e int) bool {
	if b == e {
		return tr.traverse(r + 1)
	}
	columns := tr.columns[r]
	numbers := tr.numbers[r]
	glog.V(6).Infof("row %d range [%d, %d] numbers: %v\n", r, b, e, numbers)
	for i := b; i < e; i++ {
		c := columns[b] // 确定第b列应该填哪个数字
		n := numbers[i]
		m := BitMask(1) << uint(n)
		block := getBlockByPos(r, c)
		// 检查第i个number能否放倒(r, c)的位置
		if tr.cmask[c].IsSet(m) || tr.bmask[block].IsSet(m) {
			continue
		}
		glog.V(6).Infof("number at[%d,%d]=%d, m=%s, cmask[%d]=%s, bmask[%d]=%s\n",
			r, c, n, m.String(), c, tr.cmask[c].String(), block, tr.bmask[block].String())
		tr.cmask[c].Set(m)
		tr.bmask[block].Set(m)
		numbers[b], numbers[i] = numbers[i], numbers[b]
		if tr.traverseRow(r, b+1, e) {
			return true
		}
		numbers[b], numbers[i] = numbers[i], numbers[b]
		tr.cmask[c].Clear(m)
		tr.bmask[block].Clear(m)
	}
	return false
}

func (tr *traverseResolver) Result() []*resolveResult {
	result := []*resolveResult{}
	for i := 0; i < 9; i++ {
		for j := 0; j < len(tr.columns[i]); j++ {
			r := &resolveResult{
				Row: i,
				Col: tr.columns[i][j],
				Num: tr.numbers[i][j],
			}
			result = append(result, r)
		}
	}
	return result
}

func (tr *traverseResolver) Resolve() error {
	if tr.traverse(0) {
		return nil
	}
	return errors.New("no solution")
}
