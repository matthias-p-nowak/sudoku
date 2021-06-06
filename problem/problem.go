package problem

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

var (
	lp = []int{0, 1, 2, 4, 5, 6, 8, 9, 10}
	cp = []int{1, 4, 7, 11, 14, 17, 21, 24, 27}
)

type Sudoku struct {
	Given  [9][9]int
	Placed [9][9]int
}

func (s *Sudoku) Reset() {
	var p [9][9]int
	s.Placed = p
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if s.Given[i][j] > 0 {
				s.Placed[i][j] = s.Given[i][j]
			} else {
			}
		}
	}
}

func (s *Sudoku) PrintOut() {
	po := make([][]rune, 0)
	po = append(po, []rune("         │         │         "))
	po = append(po, []rune("         │         │         "))
	po = append(po, []rune("         │         │         "))
	po = append(po, []rune("─────────┼─────────┼─────────"))
	po = append(po, []rune("         │         │         "))
	po = append(po, []rune("         │         │         "))
	po = append(po, []rune("         │         │         "))
	po = append(po, []rune("─────────┼─────────┼─────────"))
	po = append(po, []rune("         │         │         "))
	po = append(po, []rune("         │         │         "))
	po = append(po, []rune("         │         │         "))
	for i := 0; i < 9; i++ {
		l := lp[i]
		for j := 0; j < 9; j++ {
			if s.Placed[i][j] > 0 {
				c := cp[j]
				po[l][c] = rune('0' + s.Placed[i][j])
			}
		}
	}
	for _, line := range po {
		fmt.Printf("  %s\n", string(line))
	}
}

func ReadFile(fn string) (sudoku *Sudoku) {
	f, err := os.Open(fn)
	if err != nil {
		log.Fatal("couldn't open file "+fn, err)
	}
	defer f.Close()
	sudoku = new(Sudoku)
	sc := bufio.NewScanner(f)
	i := 0
	for sc.Scan() {
		givens := sc.Text() + "         "
		for j := 0; j < 9; j++ {
			c := int(givens[j] - '0')
			if c > 0 && c <= 9 {
				sudoku.Given[i][j] = c
			}
		}
		i++
	}
	return
}
