// Package sudoku contains sudoku solvers following different strategies.
package sudoku

import (
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

// Board contains all fields of a simple, unannotated sudoku.
type Board [][]int

// NewEmptyBoard builds an empty board with provided size.
func NewEmptyBoard(size int) Board {
	bo := make(Board, size)
	for i := range bo {
		bo[i] = make([]int, size)
	}
	return bo
}

// Copy returns a copy of the board. Helpful to stay as immutible as possible for now.
func (bo Board) Copy() (bo2 Board) {
	size := bo.Size()
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

// Full returns if there are no empty fields left. It does not check correctness.
func (bo Board) Full() bool {
	_, _, found := bo.FirstEmpty()
	return !found
}

// Size returns the size of the sides of the board
func (bo Board) Size() int {
	return len(bo[0])
}

// NewRandomBoard returns a solved, pseudo random board of provided size.
func NewRandomBoard(size int) Board {
	ab, _ := NewAnnotatedBoard(NewEmptyBoard(size))
	for i, v := range rand.Perm(size) {
		ab.Board[0][i] = v + 1
	}
	_, solutions := Backtrack(ab, 10)

	return solutions[rand.Intn(len(solutions))]
}

const abc = "abcdefghijklmnopqrstuvwxyz"

var reShortNotation = regexp.MustCompile("([a-z][1-9][1-9])")

// NewFromShort tries to parse provided short notation into board with provided size.
func NewFromShort(s string, size int) (bo Board, err error) {
	bo = NewEmptyBoard(size)
	if s == "" {
		return
	}
	if !reShortNotation.MatchString(s) {
		return nil, fmt.Errorf("Malformed Short Notation: %s", s)
	}

	for _, triple := range reShortNotation.FindAllString(s, len(s)/3) {
		y := strings.Index(abc, string(triple[0]))
		x, _ := strconv.Atoi(string(triple[1]))
		v, _ := strconv.Atoi(string(triple[2]))
		bo[y][x-1] = v
	}
	return
}

// Short returns simple sudoku string representation, e.g. "a18b52".
func (bo Board) Short() (s string, err error) {
	if bo.Size() > 9 {
		return "", errors.New("Board Too Large For Short Notation")
	}
	for y, row := range bo {
		for x, v := range row {
			if v > 0 {
				s = s + string(abc[y]) + strconv.Itoa(x+1) + strconv.Itoa(v)
			}
		}
	}
	return
}

// UltraShort returns an even shorter identifier only if board is full.
// It does so by returning only values and removing the last column and row.
// A 9x9 sudoku will thus result in 64 chars.
func (bo Board) UltraShort() (s string, err error) {
	if !bo.Full() {
		return "", errors.New("Board Is Not Full")
	}
	sizeWithoutLast := bo.Size() - 1
	for y, row := range bo {
		for x, v := range row {
			if y < sizeWithoutLast && x < sizeWithoutLast {
				s = s + strconv.Itoa(v)
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

// Remove provided number from candidates (if exists).
func (c Candidates) Remove(v int) Candidates {
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

// Complement returns the candidates not represented assuming a sudoku of provided size.
func (c Candidates) Complement(size int) Candidates {
	return allCandidates(size) - c
}

// AnnotatedBoard contains "penciled" candidates.
type AnnotatedBoard struct {
	Board
	Candidates [][]Candidates
}

// Copy returns a copy of the annotated board. Helpful to stay as immutible as possible for now.
func (ab AnnotatedBoard) Copy() (ab2 AnnotatedBoard) {
	size := ab.Size()
	ab2.Candidates = make([][]Candidates, size)
	for i := range ab.Candidates {
		ab2.Candidates[i] = make([]Candidates, len(ab.Candidates[i]))
		copy(ab2.Candidates[i], ab.Candidates[i])
	}
	ab2.Board = ab.Board.Copy()
	return
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
	size := ab.Size()
	rt := sqrt(size)
	rows, cols := make([]Candidates, size), make([]Candidates, size)
	ab.Candidates = newBlockCandidates(size)
	blocks := newBlockCandidates(rt)
	for y, row := range ab.Board {
		for x, v := range row {
			if v > 0 {
				if rows[y].Contains(v) || cols[x].Contains(v) || blocks[y/rt][x/rt].Contains(v) {
					return ab, errors.New("Not Solvable.")
				}
				ab.Candidates[y][x] = ab.Candidates[y][x].Add(v)
				rows[y] = rows[y].Add(v)
				cols[x] = cols[x].Add(v)
				blocks[y/rt][x/rt] = blocks[y/rt][x/rt].Add(v)
			}
		}
	}
	for y, row := range ab.Candidates {
		for x, v := range row {
			if v > 1 {
				continue
			}
			ab.Candidates[y][x] = allCandidates(size) &^ rows[y] &^ cols[x] &^ blocks[y/rt][x/rt]
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
	ab = ab.Copy()

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

// A Simplifier does not solve a sudoku but tries to remove candidates.
type Simplifier func(ab AnnotatedBoard) (ab2 AnnotatedBoard, succeeded bool)

// CandidateLines tries to simplify by finding candidates of the same kind on the same row or column within a block, so that they can be safely removed from other fields of that line not within that block.
func CandidateLines(ab AnnotatedBoard) (ab2 AnnotatedBoard, succeeded bool) {
	ab2 = ab.Copy()
	size := ab2.Size()
	blockSize := sqrt(size)

	// Iterate over blocks.
	for blkY := 0; blkY < blockSize; blkY++ {
		for blkX := 0; blkX < blockSize; blkX++ {
			// Build other rows and cols for easier removal of found candidates from them.
			rowsNotInBlock, colsNotInBlock := allCandidates(size), allCandidates(size)
			// Build maps for indices per candidates (if a candidate is in a single row/col - that will be a win!).
			inRows, inCols := map[int]Candidates{}, map[int]Candidates{}
			// Iterate over rows in block.
			for yInBlk := 0; yInBlk < blockSize; yInBlk++ {
				y := blkY*blockSize + yInBlk
				rowsNotInBlock = rowsNotInBlock.Remove(y + 1)
				// Iterate over cols in block.
				for xInBlk := 0; xInBlk < blockSize; xInBlk++ {
					x := blkX*blockSize + xInBlk
					colsNotInBlock = colsNotInBlock.Remove(x + 1)
					cs := ab2.Candidates[y][x]
					// Add one-based candidate indices to maps.
					for _, c := range cs.Decimals() {
						inRows[c] = inRows[c].Add(y + 1)
						inCols[c] = inCols[c].Add(x + 1)
					}
				}
			}

			// Remove found line candidates from other lines.
			for c, cols := range inCols {
				if cols.Single() {
					col := cols.Decimals()[0] - 1
					for _, row := range rowsNotInBlock.Decimals() {
						if ab2.Candidates[row-1][col].Contains(c) {
							succeeded = true
							ab2.Candidates[row-1][col] = ab2.Candidates[row-1][col].Remove(c)
						}
					}
				}
			}
			for c, rows := range inRows {
				if rows.Single() {
					row := rows.Decimals()[0] - 1
					for _, col := range colsNotInBlock.Decimals() {
						if ab2.Candidates[row][col-1].Contains(c) {
							succeeded = true
							ab2.Candidates[row][col-1] = ab2.Candidates[row][col-1].Remove(c)
						}
					}
				}
			}
		}
	}

	return
}

// GenerateSimple generates a board that is solvable with only single candidates strategy.
// Minimum param sets the number of minimum numbers that should remain on sudoku.
// minimum <= 0 will be ignored (hardest) and the higher the easier it gets.
// minimum >= size² makes no sense.
func GenerateSimple(random Board, minimum int) (unsolved Board) {
	size := random.Size()
	fields := [][3]int{}
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			fields = append(fields, [3]int{y, x, random[y][x]})
		}
	}

	for i := range fields {
		j := rand.Intn(i + 1)
		fields[i], fields[j] = fields[j], fields[i]
	}

	ab, _ := NewAnnotatedBoard(random)
	fieldCount := len(fields)
	validMinimum := minimum > 0 && minimum < fieldCount

	for i, f := range fields {
		ab.Board[f[0]][f[1]] = 0
		ab, _ = ab.Annotate()
		solvable, _ := SingleCandidate(ab, 1)
		if !solvable {
			ab.Board[f[0]][f[1]] = f[2]
		}
		if validMinimum && (fieldCount-i-1) == minimum {
			return ab.Board
		}
	}
	return ab.Board
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
