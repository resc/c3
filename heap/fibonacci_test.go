package heap

import "testing"

func TestFibInsert(t *testing.T) {
	h := NewFibonacci()
	h.Insert(1, 0.5)
	h.Insert(2, 0.2)
	h.Insert(3, 0.1)
	if x, ok := h.DeleteMin(); !ok || x != 3 {
		t.Error("DeleteMin failed")
	}
}
