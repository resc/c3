package c3

import (
	"testing"
)

func TestAdd(t *testing.T) {
	l := NewList()

	if !l.Add(1) || l.Len() != 1 {
		t.Errorf("Len(): expected %v, got %v", 1, l.Len())
	}
}
func TestAddNil(t *testing.T) {
	l := ListOf(nil, nil, nil)

	if l.Len() != 3 {
		t.Errorf("Len(): expected %v, got %v", 3, l.Len())
	}
	assertContains(t, l, nil, true)
}

func TestIndexOfNil(t *testing.T) {
	l := ListOf(nil, nil, nil)

	assertIndexOf(t, l, nil, 0)

	index, ok := l.LastIndexOf(nil)
	if !ok || index != 2 {
		t.Error("Expected %v,got%v", 2, index)
	}

	index, ok = l.NextIndexOf(0, nil)
	if !ok || index != 1 {
		t.Error("Expected %v,got%v", 1, index)
	}
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
	l := CollectionOf(1, 2, 3)

	if !l.Contains(2) {
		t.Errorf("Contains(...): expected %v, got %v", false, ToSlice(l))
	}
}

func TestIndexOf(t *testing.T) {
	l := NewList()
	l.Add(1)
	l.Add(2)
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

func assertContains(t *testing.T, l ReadOnlyCollection, item interface{}, expected bool) {
	result := l.Contains(item)
	if expected != result {
		t.Errorf("Contains(%v): expected %v, got %v", item, expected, result)
	}
}

func assertIndexOf(t *testing.T, l ReadOnlyList, item interface{}, expected int) {
	result, ok := l.IndexOf(item)
	if !ok || expected != result {
		t.Errorf("IndexOf(%v): expected %v,%v got %v,%v", item, expected, true, result, ok)
	}
}
