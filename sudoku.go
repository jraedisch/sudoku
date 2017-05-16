// Package sudoku contains sudoku solvers following different strategies.
package sudoku

import (
	"errors"
	"log"
	"math"
	"strconv"
)

// Board contains all fields of a simple, unannotated sudoku.
type Board [][]int

// Copy returns a copy of the board. Helpful to stay as immutible as possible for now.
func (bo Board) Copy() (bo2 Board) {
	size := len(bo)
	bo2 = make(Board, size)
	for i := range bo {
		bo2[i] = make([]int, size)
		copy(bo2[i], bo[i])
	}
	return
}

// FirstEmpty returns coordinates of first empty field and whether one was found.
func (bo Board) FirstEmpty() (y, x int, found bool) {
	for y = range bo {
		for x = range bo[y] {
			if bo[y][x] == 0 {
				found = true
				return
			}
		}
	}
	return
}

const abc = "abcdefghijklmnopqrstuvwxyz"

// Short returns simple sudoku string representation, e.g. "a18b52".
func (bo Board) Short() (s string) {
	for y, row := range bo {
		for x, v := range row {
			if v > 0 {
				s += string(abc[x]) + strconv.Itoa(y+1) + strconv.Itoa(v)
			}
		}
	}
	return
}

// Candidates contains all "penciled" numbers that may occupy a field.
type Candidates int

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

// Single returns whether candidates contain a single candidate.
func (c Candidates) Single() bool {
	// Check if c is a power of two.
	return (c > 0 && ((c & (c - 1)) == 0))
}

// Decimals returns candidates as a decimal array.
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

// AnnotatedBoard enriches a board with various annotation ("penceling") helpers.
type AnnotatedBoard struct {
	Board
	Candidates, Blocks [][]Candidates
	Rows, Cols         []Candidates
}

// NewAnnotatedBoard returns an annotated version of provided board.
func NewAnnotatedBoard(bo Board) (ab AnnotatedBoard, err error) {
	ab = AnnotatedBoard{Board: bo.Copy()}
	return ab.Annotate()
}

// Solved returns true if board is solved correctly.
func (ab AnnotatedBoard) Solved() bool {
	for _, row := range ab.Candidates {
		for _, v := range row {
			if !v.Single() {
				return false
			}
		}
	}
	return true
}

// Annotate (naively) annotates a board with possible candidates for each field.
// All data except board will be overwritten.
func (ab AnnotatedBoard) Annotate() (AnnotatedBoard, error) {
	ab.Board = ab.Board.Copy()
	le := len(ab.Board[0])
	rt := sqrt(le)
	ab.Rows, ab.Cols = make([]Candidates, le), make([]Candidates, le)
	ab.Candidates, ab.Blocks = newBlockCandidates(le), newBlockCandidates(rt)
	for y, row := range ab.Board {
		for x, v := range row {
			if v > 0 {
				if ab.Rows[y].Contains(v) || ab.Cols[x].Contains(v) || ab.Blocks[y/rt][x/rt].Contains(v) {
					return ab, errors.New("Not Solvable.")
				}
				ab.Candidates[y][x] = ab.Candidates[y][x].Add(v)
				ab.Rows[y] = ab.Rows[y].Add(v)
				ab.Cols[x] = ab.Cols[x].Add(v)
				ab.Blocks[y/rt][x/rt] = ab.Blocks[y/rt][x/rt].Add(v)
			}
		}
	}
	for y, row := range ab.Candidates {
		for x, v := range row {
			if v > 1 {
				continue
			}
			ab.Candidates[y][x] = allCandidates(le) &^ ab.Rows[y] &^ ab.Cols[x] &^ ab.Blocks[y/rt][x/rt]
		}
	}
	return ab, nil
}

// A Parser is a function that parses a certain string representation of a sudoku into a board.
type Parser func(unparsed string) (bo Board, err error)

// A Solver is a function that tries to solve an annotated board following a certain strategy.
// An integer should be provided to limit the amount of calculated solutions.
// -1 means, as many solutions as possible will be calculated.
type Solver func(ab AnnotatedBoard, maxSolutions int) (solved bool, solutions []Board)

// SingleCandidate will look through annotated board for fields with just one candidate, fill it, recalculate candidates for all fields, and repeat, until board is solved or no longer solvable.
// It will not return more than one solution.
func SingleCandidate(ab AnnotatedBoard, maxSolutions int) (bool, []Board) {
	solvable := true
	ab.Board = ab.Copy()

	for solvable {
		solvable = false
		for y, row := range ab.Candidates {
			for x, c := range row {
				if !c.Single() || ab.Board[y][x] > 0 {
					continue
				}

				ab.Board[y][x] = c.Decimals()[0]
				solvable = true
			}
		}
		ab, _ = ab.Annotate()
	}
	return ab.Solved(), []Board{ab.Board}
}

// Backtrack implements a simple backtracking solver. It is not performant but guaranteed to finish.
func Backtrack(ab AnnotatedBoard, maxSolutions int) (solved bool, solutions []Board) {
	solutions = []Board{}
	return backtrack(ab, maxSolutions, &solutions), solutions
}

func backtrack(ab AnnotatedBoard, maxSolutions int, solutions *[]Board) bool {
	var err error
	ab, err = ab.Annotate()
	if err != nil {
		return false
	}
	y, x, found := ab.Board.FirstEmpty()
	if !found {
		*solutions = append(*solutions, ab.Board)
		return len(*solutions) >= maxSolutions
	}
	for _, v := range ab.Candidates[y][x].Decimals() {
		ab.Board[y][x] = v
		if backtrack(ab, maxSolutions, solutions) {
			return true
		}
	}
	return false
}

// Helpers

// allCandidates returns all candidates for a field of a sudoku with provided size.
func allCandidates(size int) (c Candidates) {
	for size > 0 {
		c = c.Add(size)
		size--
	}
	return
}

func newBlockCandidates(size int) (bc [][]Candidates) {
	bc = make([][]Candidates, size)
	for y := range bc {
		bc[y] = make([]Candidates, size)
	}
	return bc
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
	log.Print(vs...)
}
