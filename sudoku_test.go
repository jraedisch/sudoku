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

var empty4 = Board{
	{1, 0, 0, 0},
	{0, 2, 0, 0},
	{0, 0, 3, 0},
	{0, 0, 0, 4},
}

var annotated4 = AnnotatedBoard{
	{2, 24, 20, 12},
	{24, 4, 18, 10},
	{20, 18, 8, 6},
	{12, 10, 6, 16},
}

func TestAnnotate(t *testing.T) {
	actual := empty4.Annotate()
	if !reflect.DeepEqual(annotated4, actual) {
		t.Errorf("Expected %v, got %v", annotated4, actual)
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
	c = c.Add(9)
	if "1000000000" != b(c) {
		t.Errorf(`Expected "1000000000", got "%s".`, b(c))
	}
	c = c.Add(1)
	if "1000000010" != b(c) {
		t.Errorf(`Expected "1000000010", got "%s".`, b(c))
	}
	c = c.Subtract(7)
	if "1000000010" != b(c) {
		t.Errorf(`Expected "1000000010", got "%s".`, b(c))
	}
	c = c.Subtract(9)
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

func TestNewAnnotatedBoard(t *testing.T) {
	expected := AnnotatedBoard{
		make([]Candidates, 2),
		make([]Candidates, 2),
	}
	actual := newAnnotatedBoard(2)
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestBacktrack(t *testing.T) {

}
