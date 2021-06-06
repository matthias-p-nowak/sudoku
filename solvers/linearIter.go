package solvers

import (
	"fmt"

	"github.com/matthias-p-nowak/sudoku/problem"
)

type linVar struct {
	cnt   int
	val   float32
	digit int
}

type position struct {
	digits []*linVar
}

func (p *position) addVar(v *linVar) {
	p.digits = append(p.digits, v)
}

type restriction struct {
	vars []*linVar
	name string
}

func (r *restriction) addVar(l *linVar) {
	r.vars = append(r.vars, l)
}

func (r *restriction) empty() bool {
	return len(r.vars) == 0
}

type linProb struct {
	digits [9][9][9]bool
	allPos [9][9]position
	res    []*restriction
	pos    []*position
}

func (li *linProb) addRes(r *restriction) {
	li.res = append(li.res, r)
}

func (li *linProb) mark(i, j, k int) {
	for l := 0; l < 9; l++ {
		// deleting row
		li.digits[i][l][k-1] = false
		// deleting column
		li.digits[l][j][k-1] = false
		// deleting digits
		li.digits[i][j][l] = false
	}
	bi := (i / 3) * 3
	bj := (j / 3) * 3
	for m := 0; m < 3; m++ {
		for n := 0; n < 3; n++ {
			li.digits[bi+m][bj+n][k-1] = false
		}
	}
}

func (li *linProb) initialize(s *problem.Sudoku) {
	var (
		allVars [9][9][9]*linVar
		resRows [9][9]*restriction
		resCols [9][9]*restriction
		resBoxs [9][9]*restriction
	)
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			r := new(restriction)
			r.name = fmt.Sprintf("r%dd%d", i, j+1)
			resRows[i][j] = r
			r = new(restriction)
			r.name = fmt.Sprintf("c%dd%d", i, j+1)
			resCols[i][j] = r
			r = new(restriction)
			r.name = fmt.Sprintf("b%dd%d", i, j+1)
			resBoxs[i][j] = r
			for k := 0; k < 9; k++ {
				li.digits[i][j][k] = true

			}
		}
	}
	// from givens
	cntGivens := 0
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			k := s.Given[i][j]
			if k > 0 {
				li.mark(i, j, k)
				cntGivens++
			}
		}
	}
	//
	cntVars := 0
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			for k := 0; k < 9; k++ {
				if li.digits[i][j][k] {
					v := new(linVar)
					v.digit = k + 1
					v.cnt = 1 // crude start
					allVars[i][j][k] = v
					cntVars++
					resRows[i][k].addVar(v)
					resCols[j][k].addVar(v)
					b := (i/3)*3 + (j / 3)
					resBoxs[b][k].addVar(v)
					li.allPos[i][j].addVar(v)
				}
			}
		}
	}
	li.res = nil
	for m := 0; m < 9; m++ {
		for n := 0; n < 9; n++ {
			r := resRows[m][n]
			if !r.empty() {
				li.addRes(r)
			}
			r = resCols[m][n]
			if !r.empty() {
				li.addRes(r)
			}
			r = resBoxs[m][n]
			if !r.empty() {
				li.addRes(r)
			}
			if len(li.allPos[m][n].digits) > 0 {
				li.pos = append(li.pos, &li.allPos[m][n])
			}
		}
	}
	fmt.Printf("got %d variables %d givens, %d restrictions (total %d)\n", cntVars, cntGivens, len(li.res), 81-cntGivens+len(li.res))
}

func (li *linProb) solve() {
	var (
		maxRes *restriction
	)
	for iter := 0; iter < 1000; iter++ {
		// adjust the values based on the counts
		for _, pos := range li.pos {
			cnt := float32(0)
			for _, dg := range pos.digits {
				cnt += float32(dg.cnt)
			}
			for _, dg := range pos.digits {
				dg.val = float32(dg.cnt) / cnt
			}
		}
		// go through the restrictions
		maxVal := float32(1)
		maxRes = nil
		for _, res := range li.res {
			val := float32(0)
			for _, v := range res.vars {
				val += v.val
			}
			if val < maxVal {
				maxRes = res
				maxVal = val
			}
		}
		if maxRes == nil {
			break
		}
		fmt.Printf("increasing vars for %s (%f)\n", maxRes.name, maxVal)
		for _, v := range maxRes.vars {
			v.cnt++
		}
	}
}

func SolveLinIt(s *problem.Sudoku) {
	defer fmt.Println("linear iteration done")
	li := new(linProb)
	li.initialize(s)
	li.solve()
}
