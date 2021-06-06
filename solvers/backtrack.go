package solvers

import (
	"fmt"

	"github.com/matthias-p-nowak/sudoku/problem"
)

type placed struct {
	// placed digit k into row i column j
	i, j, k int
	n       *placed
}
type backTrack struct {
	digits [9][9][9]bool
	value  [9][9]int
	// rows   [9][9]bool
	// cols   [9][9]bool
	// boxes  [9][9]bool
	placed *placed
	steps  int
}

func (bt *backTrack) mark(i, j, k int) {
	bi := (i / 3) * 3
	bj := (j / 3) * 3
	bt.value[i][j] = k
	for l := 0; l < 9; l++ {
		bt.digits[i][l][k-1] = false
		bt.digits[l][j][k-1] = false
		bt.digits[i][j][l] = false
	}
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			bt.digits[bi+i][bj+j][k-1] = false
		}
	}
}

func (bt *backTrack) reset() {
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			bt.value[i][j] = 0
			// bt.rows[i][j] = true
			// bt.cols[i][j] = true
			// bt.boxes[i][j] = true
			for k := 0; k < 9; k++ {
				bt.digits[i][j][k] = true
			}
		}
	}
	for p := bt.placed; p != nil; p = p.n {
		bt.mark(p.i, p.j, p.k)
	}
}

func (bt *backTrack) initialize(s *problem.Sudoku) {
	op := new(placed) // temporary placeholder
	bt.placed = op
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			k := s.Given[i][j]
			if k > 0 {
				op.n = new(placed)
				op = op.n
				op.i = i
				op.j = j
				op.k = k
			}
		}
	}
	bt.placed = bt.placed.n
	bt.reset()
	bt.steps = 0
}

func (bt *backTrack) branch(i, j int) bool {
	// initial, go to the end
	t := bt.placed
	for t.n != nil {
		t = t.n
	}
	var values []int
	for k := 0; k < 9; k++ {
		if bt.digits[i][j][k] {
			values = append(values, k+1)
		}
	}
	fmt.Printf("  branching on %#v\n", values)
	for idx, v := range values {
		fmt.Printf("this branch %d value %d\n", idx, v)
		t.n = new(placed)
		t.n.i = i
		t.n.j = j
		t.n.k = v
		bt.reset()
		bt.steps++
		r := bt.solve()
		if r {
			return true
		} else {
			fmt.Printf(" branch %d (%d,%d)=%d failed\n", idx, i, j, v)
		}
	}
	return false
}

func (bt *backTrack) solve() bool {
	// initial, go to the end
	t := bt.placed
	for t.n != nil {
		t = t.n
	}
	for {
		var df [10]int
		mdf := 10
		var mi, mj int
		co := 0
		for i := 0; i < 9; i++ {
			for j := 0; j < 9; j++ {
				if bt.value[i][j] > 0 {
					// filled, continue
					continue
				}
				// open cell
				co++
				c := 0
				for k := 0; k < 9; k++ {
					if bt.digits[i][j][k] {
						c++
					}
				}
				if c < mdf {
					mdf = c
					mi = i
					mj = j
				}
				df[c]++
			}
		}
		if mdf < 10 {
			fmt.Printf("found cell (%d,%d) dof=%d\n", mi, mj, mdf)
		}
		s := 0
		for k := 0; k <= 9; k++ {
			if df[k] > 0 {
				fmt.Printf(" dof=%d -> %d\n", k, df[k])
				s += df[k]
			}
		}
		fmt.Printf(" ---------\n %d to go\n", s)
		if co == 0 {
			// finished
			fmt.Println("jippi")
			return true
		}
		if mdf == 0 {
			// infeasible
			return false
		}
		if mdf == 1 {
			// continue straight on
			k := 0
			for k = 0; k < 9; k++ {
				if bt.digits[mi][mj][k] {
					break
				}
			}
			k++
			fmt.Printf("continuing with (%d,%d):=%d \n", mi, mj, k)
			t.n = new(placed)
			t = t.n
			t.i = mi
			t.j = mj
			t.k = k
			bt.mark(mi, mj, k)
			bt.steps++
			continue
		} else {
			return bt.branch(mi, mj)
		}
		fmt.Println("### error ###")
		break
	}
	return false
}

func SolveBacktrack(s *problem.Sudoku) {
	bt := new(backTrack)
	bt.initialize(s)
	if bt.placed != nil {
		bt.solve()
		fmt.Printf("steps taken %d\n", bt.steps)
		s.Reset()
		for placed := bt.placed; placed != nil; placed = placed.n {
			s.Placed[placed.i][placed.j] = placed.k
		}
		s.PrintOut()
	} else {
		fmt.Println("nothing placed")
	}
}
