package c3

import "testing"

func TestToSliceOfIterable(t *testing.T) {
	s := ToSlice(Range(0, 9))
	if len(s) != 10 {
		t.Errorf("Expected 10, got %v", s)
	}
}
func TestToSliceOfIterableReverse(t *testing.T) {
	s := ToSlice(Range(9, 0))
	if len(s) != 10 {
		t.Errorf("Expected 10, got %v", s)
	}
}
func TestToSliceOfCollection(t *testing.T) {
	r := ListOf(0, 1, 2, 3)
	s := ToSlice(r)
	if len(s) != 4 {
		t.Errorf("Expected 4, got %v", s)
	}
}
