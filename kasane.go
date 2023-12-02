package kasane

import (
	"slices"
	"strings"

	"github.com/mattn/go-runewidth"
)

type Option func(*k)

type k struct {
	rc *runewidth.Condition
}

func new() *k {
	return &k{
		rc: runewidth.NewCondition(),
	}
}

func OverlayString(base, s string, top, left int, opts ...Option) string {
	k := new()
	for _, opt := range opts {
		opt(k)
	}

	bs := strings.Split(base, "\n")
	ss := strings.Split(s, "\n")
	ret := make([]string, 0)
	for i, b := range bs {
		if top <= i && i < top+len(ss) {
			ret = append(ret, k.overlaySingleLineString(b, ss[i-top], left, opts...))
		} else {
			ret = append(ret, b)
		}
	}

	return strings.Join(ret, "\n")
}

func (k *k) overlaySingleLineString(base, s string, left int, opts ...Option) string {

	baseCells := k.toCells(base)
	sCells := k.toCells(s)
	bw := len(baseCells)
	sw := len(sCells)

	for si := 0; si < sw; si++ {
		bi := si + left
		if 0 <= bi && bi < bw {
			sCell := sCells[si]
			baseCell := baseCells[bi]
			if baseCell.double {
				if baseCell.head {
					if bi+1 < bw {
						baseCells[bi+1] = emptyCell()
					}
				} else {
					if bi-1 >= 0 {
						baseCells[bi-1] = emptyCell()
					}
				}
			}
			if sCell.double && !sCell.head && bi == 0 {
				baseCells[bi] = emptyCell()
				continue
			}
			if sCell.double && bi == bw-1 {
				baseCells[bi] = emptyCell()
				break
			}
			baseCells[bi] = sCell
		}
	}

	rs := make([]rune, 0)
	for i := range baseCells {
		if i > 0 {
			if len(baseCells[i].csis) == 0 {
				if len(baseCells[i-1].csis) > 0 {
					rs = append(rs, []rune(csiReset)...)
				}
			} else {
				if len(baseCells[i-1].csis) > 0 && !overlapped(baseCells[i].csis, baseCells[i-1].csis) {
					rs = append(rs, []rune(csiReset)...)
				}
			}
		}
		for _, c := range baseCells[i].csis {
			if i == 0 || (i > 0 && !slices.Contains(baseCells[i-1].csis, c)) {
				rs = append(rs, []rune(c)...)
			}
		}
		if baseCells[i].head {
			rs = append(rs, baseCells[i].r)
		}
		if i == bw-1 && len(baseCells[i].csis) > 0 {
			rs = append(rs, []rune(csiReset)...)
		}
	}

	return string(rs)
}

func ansiStart(r rune) bool {
	return r == '\x1b'
}

func ansiEnd(r rune) bool {
	return (0x40 <= r && r <= 0x5a) || (0x61 <= r && r <= 0x7a)
}

func (k *k) width(r rune) int {
	return k.rc.RuneWidth(r)
}

func (k *k) stringWidth(s string) int {
	ansi := false
	n := 0
	for _, r := range s {
		if ansiStart(r) {
			ansi = true
		} else if ansi {
			if ansiEnd(r) {
				ansi = false
			}
		} else {
			n += k.width(r)
		}
	}
	return n
}

type cell struct {
	r      rune
	double bool
	head   bool
	csis   []string
}

func emptyCell() *cell {
	return &cell{r: ' ', double: false, head: true}
}

const (
	csiReset = "\x1b[0m"
)

func (k *k) toCells(s string) []*cell {
	cells := make([]*cell, 0)
	ansi := false
	ansiBuf := make([]rune, 0)
	csis := make([]string, 0)
	for _, r := range s {
		if ansiStart(r) {
			ansiBuf = append(ansiBuf, r)
			ansi = true
		} else if ansi {
			ansiBuf = append(ansiBuf, r)
			csi := string(ansiBuf)
			if ansiEnd(r) {
				ansi = false
				ansiBuf = make([]rune, 0) // reset
				if csi == csiReset {
					csis = make([]string, 0) // reset
				} else {
					csis = append(csis, csi)
				}
			}
		} else {
			w := k.width(r)
			for i := 0; i < w; i++ {
				cell := &cell{
					r:      r,
					double: w == 2,
					head:   i == 0,
					csis:   copySlice(csis),
				}
				cells = append(cells, cell)
			}
		}
	}
	return cells
}

func copySlice[T any](src []T) []T {
	dst := make([]T, len(src))
	_ = copy(dst, src)
	return dst
}

func overlapped[T comparable](xs, ys []T) bool {
	for _, x := range xs {
		for _, y := range ys {
			if x == y {
				return true
			}
		}
	}
	return false
}
