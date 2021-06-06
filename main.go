package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/matthias-p-nowak/sudoku/cmd"
	"github.com/matthias-p-nowak/sudoku/problem"
	"github.com/matthias-p-nowak/sudoku/solvers"
)

func main() {
	fmt.Println("welcome to github.com/matthias-p-nowak/sudoku")
	defer fmt.Println("all done")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	ia := flag.Bool("i", false, "interactive usage")
	flag.Parse()
	for _, fileName := range flag.Args() {
		fmt.Printf("reading file %s\n", fileName)
		s := problem.ReadFile(fileName)
		s.Reset()
		s.PrintOut()
		if *ia {
			cmd.Interact(s)
		} else {
			solvers.SolveBacktrack(s)
			solvers.SolveLinIt(s)
		}
	}
}
