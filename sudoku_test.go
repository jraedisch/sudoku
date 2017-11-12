package sudoku

import (
	"reflect"
	"testing"
)

var solved4 = Board{
	{1, 2, 3, 4},
	{3, 4, 1, 2},
	{4, 1, 2, 3},
	{2, 3, 4, 1},
}

var unsolved4 = Board{
	{1, 0, 0, 0},
	{0, 2, 0, 0},
	{0, 0, 3, 0},
	{0, 0, 0, 4},
}

var unsolved4b = Board{
	{1, 2, 0, 0},
	{3, 0, 0, 2},
	{4, 1, 2, 0},
	{2, 3, 4, 1},
}

var unsolved4c = Board{
	{0, 0, 0, 0},
	{0, 0, 0, 0},
	{0, 0, 0, 0},
	{0, 0, 0, 0},
}

var unsolved9 = Board{
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
}

var unsolved9b = Board{
	{9, 0, 7, 0, 0, 4, 3, 0, 1},
	{6, 0, 4, 3, 0, 1, 9, 0, 7},
	{3, 0, 1, 9, 0, 7, 6, 0, 4},
	{0, 9, 6, 7, 4, 5, 2, 1, 3},
	{7, 0, 5, 2, 1, 3, 0, 9, 6},
	{2, 1, 3, 8, 9, 0, 7, 4, 5},
	{5, 7, 9, 4, 0, 0, 1, 0, 2},
	{4, 0, 0, 1, 3, 2, 5, 7, 9},
	{1, 0, 2, 5, 7, 9, 4, 0, 0},
}

var solved9 = Board{
	{9, 8, 7, 6, 5, 4, 3, 2, 1},
	{6, 5, 4, 3, 2, 1, 9, 8, 7},
	{3, 2, 1, 9, 8, 7, 6, 5, 4},
	{8, 9, 6, 7, 4, 5, 2, 1, 3},
	{7, 4, 5, 2, 1, 3, 8, 9, 6},
	{2, 1, 3, 8, 9, 6, 7, 4, 5},
	{5, 7, 9, 4, 6, 8, 1, 3, 2},
	{4, 6, 8, 1, 3, 2, 5, 7, 9},
	{1, 3, 2, 5, 7, 9, 4, 6, 8},
}

var annotated4 = [][]Candidates{
	{2, 24, 20, 12},
	{24, 4, 18, 10},
	{20, 18, 8, 6},
	{12, 10, 6, 16},
}

func TestShort(t *testing.T) {
	_, err := NewFromShort("", 9)
	if err != nil {
		t.Error("Unexpected error for empty short notation.")
	}

	expected := "a11a22a33a44b13b24b31b42c14c21c32c43d12d23d34d41"
	actual, _ := solved4.Short()
	if expected != actual {
		t.Errorf("Expected short notations to match:\n%s\n%s", expected, actual)
	}

	short, _ := unsolved9b.Short()
	parsed, _ := NewFromShort(short, 9)
	if !reflect.DeepEqual(unsolved9b, parsed) {
		t.Errorf("Expected original to equal parsed board:\n%+v\n%+v", unsolved9b, parsed)
	}
}

func TestUltraShort(t *testing.T) {
	expected := "9876543265432198321987658967452174521389213896745794681346813257"
	actual, err := solved9.UltraShort()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if expected != actual {
		t.Errorf("Expected short notations to match:\n%s\n%s", expected, actual)
	}
}

func TestFull(t *testing.T) {
	if true != solved9.Full() {
		t.Error("Expected solved board to be full.")
	}
}

func TestGenerateSimple(t *testing.T) {
	size := 9
	solved := NewRandomBoard(size)
	unsolved := GenerateSimple(solved, 70)

	for y, r := range unsolved {
		for x, c := range r {
			if c != 0 && c != solved[y][x] {
				t.Fatalf("Expected solved to contain unsolved:\n%+v\n%+v", solved, unsolved)
			}
		}
	}

	expectedLength := 210
	actualShort, _ := unsolved.Short()
	actualLength := len(actualShort)

	if expectedLength != actualLength {
		t.Errorf("expected length %d, got %d", expectedLength, actualLength)
	}
}

func TestCandidateLines(t *testing.T) {
	ab, _ := NewAnnotatedBoard(unsolved9b)

	// Col 2
	if !ab.Candidates[0][1].Contains(8) {
		t.Error("Expected Candidates[0][1] to contain 8.")
	}
	if !ab.Candidates[4][1].Contains(8) {
		t.Error("Expected Candidates[4][1] to contain 8.")
	}
	if !ab.Candidates[7][1].Contains(8) {
		t.Error("Expected Candidates[7][1] to contain 8.")
	}
	if !ab.Candidates[8][1].Contains(8) {
		t.Error("Expected Candidates[8][1] to contain 8.")
	}

	// Row 7
	if !ab.Candidates[6][4].Contains(6) {
		t.Error("Expected Candidates[6][4] to contain 6.")
	}
	if !ab.Candidates[6][5].Contains(6) {
		t.Error("Expected Candidates[6][5] to contain 6.")
	}
	if !ab.Candidates[6][7].Contains(6) {
		t.Error("Expected Candidates[6][7] to contain 6.")
	}

	ab2, _ := CandidateLines(ab)

	// Col 2
	if !ab2.Candidates[0][1].Contains(8) {
		t.Error("Expected simplified Candidates[0][1] to still contain 8.")
	}
	if ab2.Candidates[4][1].Contains(8) {
		t.Error("Expected simplified Candidates[4][1] not to contain 8.")
	}
	if ab2.Candidates[7][1].Contains(8) {
		t.Error("Expected simplified Candidates[7][1] not to contain 8.")
	}
	if ab2.Candidates[8][1].Contains(8) {
		t.Error("Expected simplified Candidates[8][1] not to contain 8.")
	}

	// Row 7
	if !ab2.Candidates[6][4].Contains(6) {
		t.Error("Expected simplified Candidates[6][4] to contain 6.")
	}
	if ab2.Candidates[6][5].Contains(6) {
		t.Error("Expected simplified Candidates[6][5] to contain 6.")
	}
	if ab2.Candidates[6][7].Contains(6) {
		t.Error("Expected simplified Candidates[6][7] not to contain 6.")
	}
}

func TestFirstEmpty(t *testing.T) {
	y, x, found := unsolved4b.FirstEmpty()
	if !found {
		t.Error("Expected value to be found.")
	}
	if y != 0 {
		t.Errorf("Expected 0, got %d", y)
	}
	if x != 2 {
		t.Errorf("Expected 2, got %d", x)
	}
	_, _, found = solved4.FirstEmpty()
	if found {
		t.Error("Expected value not to be found.")
	}
}

func TestBacktrack(t *testing.T) {
	annotated, _ := NewAnnotatedBoard(unsolved9)
	solved, solutions := Backtrack(annotated, 4)
	if !solved {
		t.Error("Expected board to be solved.")
	}
	if len(solutions) != 4 {
		t.Errorf("Expected 4 solutions, got %d", len(solutions))
	}
	if !reflect.DeepEqual(solved9, solutions[0]) {
		t.Errorf("Expected %+v, got %+v", solved9, solutions[0])
	}
}

func TestSingleCandidate(t *testing.T) {
	expectedSolutions := []Board{solved4}
	annotated, _ := NewAnnotatedBoard(unsolved4b)
	solved, solutions := SingleCandidate(annotated, 1)
	if !solved {
		t.Error("Expected board to be solved.")
	}
	if !reflect.DeepEqual(expectedSolutions, solutions) {
		t.Errorf("Expected %+v, got %+v", expectedSolutions, solutions)
	}
}

func TestSingleCandidateUnsolvable(t *testing.T) {
	expectedSolutions := []Board{unsolved4}
	annotated, _ := NewAnnotatedBoard(unsolved4)
	solved, solutions := SingleCandidate(annotated, 1)
	if solved {
		t.Error("Expected board not to be solved.")
	}
	if !reflect.DeepEqual(expectedSolutions, solutions) {
		t.Errorf("Expected %+v, got %+v", expectedSolutions, solutions)
	}
}

func TestAnnotate(t *testing.T) {
	actual, _ := NewAnnotatedBoard(unsolved4)
	if !reflect.DeepEqual(annotated4, actual.Candidates) {
		t.Errorf("Expected %v, got %v", annotated4, actual.Candidates)
	}
}

func TestDecimals(t *testing.T) {
	expected := []int{4, 3, 2, 1}
	actual := allCandidates(4).Decimals()
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %+v, got %+v", expected, actual)
	}
}

func TestCandidates(t *testing.T) {
	var c Candidates
	if "0" != b(c) {
		t.Errorf(`Expected "0", got "%s"`, b(c))
	}
	if c.Single() {
		t.Errorf("Expected %+v not to be single.", c)
	}
	c = c.Add(9)
	if "1000000000" != b(c) {
		t.Errorf(`Expected "1000000000", got "%s".`, b(c))
	}
	if !c.Single() {
		t.Errorf("Expected %+v to be single.", c)
	}
	c = c.Add(1)
	if "1000000010" != b(c) {
		t.Errorf(`Expected "1000000010", got "%s".`, b(c))
	}
	if c.Single() {
		t.Errorf("Expected %+v not to be single.", c)
	}
	c = c.Remove(7)
	if "1000000010" != b(c) {
		t.Errorf(`Expected "1000000010", got "%s".`, b(c))
	}
	c2 := c.Complement(9)
	if "111111100" != b(c2) {
		t.Errorf(`Expected "111111100", got "%s".`, b(c2))
	}
	c = c.Remove(9)
	if "10" != b(c) {
		t.Errorf(`Expected "10", got "%s".`, b(c))
	}
	if !c.Contains(1) {
		t.Errorf("Expected %+v to contain 1.", c)
	}
	if c.Contains(9) {
		t.Errorf("Expected %+v not to contain 9.", c)
	}
	c = allCandidates(4)
	if "11110" != b(c) {
		t.Errorf(`Expected "11110", got "%s"`, b(c))
	}
}

func TestNewBlockCandidates(t *testing.T) {
	expected := [][]Candidates{
		make([]Candidates, 2),
		make([]Candidates, 2),
	}
	actual := newBlockCandidates(2)
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}
