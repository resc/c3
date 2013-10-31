package c3

import (
	"testing"
)

func TestAdd(t *testing.T) {
	l := NewList()
	assert(t, 0, l.Len(), "for Len()")
}
func TestAddNil(t *testing.T) {
	l := ListOf(nil, nil, nil)

	assert(t, 3, l.Len(), "l.Len()")
	assertContains(t, l, nil, true)
}

func TestIndexOfNil(t *testing.T) {
	l := ListOf(nil, nil, nil)

	assertIndexOf(t, l, nil, 0)

	index, ok := l.LastIndexOf(nil)
	assert(t, true, ok, "ok")
	assert(t, 2, index, "index")

	index, ok = l.NextIndexOf(0, nil)
	assert(t, true, ok, "ok")
	assert(t, 1, index, "index")

}

func TestListOf(t *testing.T) {
	l := ReadOnlyListOf(1, 2, 3)
	assertContains(t, l, 0, false)
	assertContains(t, l, 1, true)
	assertContains(t, l, 2, true)
	assertContains(t, l, 3, true)
	assertContains(t, l, 4, false)

	assertIndexOf(t, l, 1, 0)
	assertIndexOf(t, l, 2, 1)
	assertIndexOf(t, l, 3, 2)
}

func TestContains(t *testing.T) {
	l := NewList()
	l.Add(1)
	assertContains(t, l, 1, true)
}

func TestContainsFail(t *testing.T) {
	l := BagOf(1, 2, 3)
	assert(t, false, l.Contains(0), "l.Contains(0)")
	assert(t, true, l.Contains(1), "l.Contains(1)")
	assert(t, true, l.Contains(2), "l.Contains(2)")
	assert(t, true, l.Contains(3), "l.Contains(3)")
	assert(t, false, l.Contains(4), "l.Contains(4)")
}

func TestIndexOf(t *testing.T) {
	l := NewList()
	l.Add(1)
	l.Add(2)
	assertIndexOf(t, l, 1, 0)
	assertIndexOf(t, l, 2, 1)
}

func TestListSwap(t *testing.T) {
	l := NewList()
	l.Add(1)
	l.Add(2)
	l.Swap(0, 1)

	assertIndexOf(t, l, 1, 1)
	assertIndexOf(t, l, 2, 0)
	l.Swap(0, 1)

	assertIndexOf(t, l, 1, 0)
	assertIndexOf(t, l, 2, 1)
}

func TestListIteratorUninitializedValueAccess(t *testing.T) {
	l := ListOf(1, 2, 3)
	i := l.Iterator()
	value := i.Value()
	if value != nil {
		t.Error("Expected nil")
	}
}
func TestListIterator(t *testing.T) {
	l := ListOf(1, 2, 3)
	i := l.Iterator()
	if !i.MoveNext() {
		t.Error("Expected true")
	}
	value, ok := i.Value().(int)
	if !ok || value != 1 {
		t.Errorf("Expected 1, got %v", value)
	}
	if !i.MoveNext() {
		t.Error("Expected true")
	}
	value, ok = i.Value().(int)
	if !ok || value != 2 {
		t.Errorf("Expected 2, got %v", value)
	}
	if !i.MoveNext() {
		t.Error("Expected true")
	}
	value, ok = i.Value().(int)
	if !ok || value != 3 {
		t.Errorf("Expected 3, got %v", value)
	}
	if i.MoveNext() {
		t.Error("Expected false")
	}
}

func TestDelete(t *testing.T) {
	l := NewList()
	l.Add(1)

	if l.Len() != 1 {
		t.Error("Expected Len() == 1")
	}

	if l.Delete(2) {
		t.Error("Should not be able to delete non-existing item")
	}

	if !l.Delete(1) {
		t.Error("Should be able to delete existing item")
	}

	if l.Len() != 0 {
		t.Error("Expected Len() == 0")
	}

	if l.Contains(1) {
		t.Error("Expected Contains(1) == false")
	}
}

func TestInsertAtEnd(t *testing.T) {
	l := ListOf(9, 9, 9)
	ok := l.InsertAt(3, 3)
	if !ok {
		t.Error("Should be able to insert at the end of the list.")
	}
	assertIndexOf(t, l, 3, 3)
}

func TestInsertAtBegin(t *testing.T) {
	l := ListOf(9, 9, 9)
	ok := l.InsertAt(0, 0)
	if !ok {
		t.Error("Should be able to insert at the end of the list.")
	}
	assertIndexOf(t, l, 0, 0)
	assertIndexOf(t, l, 9, 1)
}

func TestIndexOfFail(t *testing.T) {
	l := NewList()
	added := l.Add(1) && l.Add(1)
	if index, ok := l.IndexOf(2); !added || ok || index >= 0 {
		t.Errorf("IndexOf(...): expected %v,%v got %v,%v", 1, true, index, ok)
	}
}
