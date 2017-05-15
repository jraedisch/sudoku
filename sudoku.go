// Package sudoku contains sudoku solvers following different strategies.
package sudoku

import (
	"fmt"
	"math"
	"strconv"
)

// Board contains all fields of a simple, unannotated sudoku.
type Board [][]int

// AnnotatedBoard contains one slice per field to hold possible "penciled" candidates.
type AnnotatedBoard [][]Candidates

// Candidates contains all "penciled" numbers that may occupy a field.
type Candidates int

// allCandidates returns all candidates for a field of a sudoku with provided size.
func allCandidates(size int) (c Candidates) {
	for size > 0 {
		c = c.Add(size)
		size--
	}
	return
}

// Add adds provided number to candidates (if not exists).
func (c Candidates) Add(v int) Candidates {
	return c | 1<<uint(v)
}

// Subtract removes provided number from candidates (if exists).
func (c Candidates) Subtract(v int) Candidates {
	return c &^ (1 << uint(v))
}

// Contains checks whether candidates contain provided number.
func (c Candidates) Contains(v int) bool {
	return c&(1<<uint(v)) != 0
}

func newAnnotatedBoard(size int) (ab AnnotatedBoard) {
	ab = make(AnnotatedBoard, size)
	for y := range ab {
		ab[y] = make([]Candidates, size)
	}
	return ab
}

// Solvable checks whether a board is still solvable (also true if solved correctly).
func (bo Board) Solvable() bool {
	return false
}

// Annotate (naively) a board with possible candidates for each field.
func (bo Board) Annotate() (ab AnnotatedBoard) {
	le := len(bo[0])
	rt := sqrt(le)
	rows, cols := make([]Candidates, le), make([]Candidates, le)
	blocks := newAnnotatedBoard(rt)
	ab = newAnnotatedBoard(le)
	for y, row := range bo {
		for x, v := range row {
			if v > 0 {
				ab[y][x] = ab[y][x].Add(v)
				rows[y] = rows[y].Add(v)
				cols[x] = cols[x].Add(v)
				blocks[y/rt][x/rt] = blocks[y/rt][x/rt].Add(v)
			}
		}
	}
	for y, row := range ab {
		for x, v := range row {
			if v > 1 {
				continue
			}
			ab[y][x] = allCandidates(le) &^ rows[y] &^ cols[x] &^ blocks[y/rt][x/rt]
		}
	}
	return ab
}

// A Parser is a function that parses a certain string representation of a sudoku into a board.
type Parser func(unparsed string) (bo Board, err error)

// A Solver is a function that tries to solve an annotated board following a certain strategy.
// An integer should be provided to limit the amount of calculated solutions.
// -1 means, as many solutions as possible will be calculated.
type Solver func(ab AnnotatedBoard, maxSolutions int) (solved bool, solutions []Board)

// SingleCandidate will look through annotated board for fields with just one candidate, fill it, recalculate candidates for all fields and repeat.
func SingleCandidate(ab AnnotatedBoard, maxSolutions int) (solved bool, solutions []Board) {
	return
}

// Backtrack implements a simple backtracking solver.
func Backtrack(ab AnnotatedBoard, maxSolutions int) (solved bool, solutions []Board) {
	return false, nil
}

// Helpers

// Decimals returns candidates as a more readable decimal array.
func (c Candidates) Decimals() (dcs []int) {
	i := int(c)
	for max := log2(i); max > 0; max-- {
		sq := pow2(max)
		if sq <= i {
			i -= sq
			dcs = append(dcs, max)
		}
	}
	return
}

func sqrt(i int) int {
	return int(math.Sqrt(float64(i)))
}

func pow2(exp int) int {
	return int(math.Pow(float64(2), float64(exp)))
}

func log2(i int) int {
	return int(math.Log2(float64(i)))
}

func b(i Candidates) string {
	return strconv.FormatInt(int64(i), 2)
}

func l(vs ...interface{}) {
	fmt.Print(vs...)
}
