package c3

import (
	"fmt"
	"path"
	"runtime"
	"testing"
)

func TestToSliceOfIterable(t *testing.T) {
	s := ToSlice(Range(0, 9))
	assert(t, 10, len(s), "len(s)")
}

func TestToSliceOfIterableReverse(t *testing.T) {
	s := ToSlice(Range(9, 0))

	assert(t, 10, len(s), "len(s)")
}

func TestToSliceOfCollection(t *testing.T) {
	r := ListOf(0, 1, 2, 3)
	s := ToSlice(r)
	assert(t, 4, len(s), "len(s)")
}

func TestSortList(t *testing.T) {
	l := ListOf(3, 1, 0, 2)
	if l.Len() < 4 {
		t.Error("empty set!")
	}

	Sort(l, func(a, b interface{}) bool {
		aval := a.(int)
		bval := b.(int)
		return aval < bval
	})

	if l.Len() < 4 {
		t.Error("empty set!")
	}

	var index = 0
	for i := l.Iterator(); i.MoveNext(); {
		val, _ := i.Value().(int)
		if val != index {
			t.Errorf("Expected %v, got %v", index, ToSlice(l))
			return
		}
		index++
	}
}

func TestRepeat(t *testing.T) {
	n := 3
	count := 0
	for i := Repeat(n, n).Iterator(); i.MoveNext(); {
		if val, ok := i.Value().(int); !ok || val != n {
			t.Errorf("Expected value %v, got %v", n, val)
		}
		count++
	}
	if count != n {
		t.Error("Expected ", n, ", got ", count)
	}
}

func TestRange(t *testing.T) {
	s := ToQueue(Range(0, 3))
	for n := 0; n < 4; n++ {
		val, ok := s.Dequeue()
		if !ok || val.(int) != n {
			t.Error("Expected ", n, " got ", val)
		}
	}
}

///////////////////////////////
// testing utility functions //
///////////////////////////////

func fail(t *testing.T, msg string) {
	_, file, line, _ := runtime.Caller(1)
	t.Errorf("\n%v:%v: Failed for %v", path.Base(file), line, msg)
}

func failf(t *testing.T, msg string, v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	t.Errorf("\n%v:%v: Failed for %v", path.Base(file), line, fmt.Sprintf(msg, v...))
}

func assert(t *testing.T, expected, actual interface{}, msg string) {
	if expected != actual {
		_, file, line, _ := runtime.Caller(1)
		t.Errorf("\n%v:%v: Expected %v, got %v for %v.", path.Base(file), line, expected, actual, msg)
	}
}

func assertContains(t *testing.T, l ReadOnlyBag, item interface{}, expected bool) {
	result := l.Contains(item)
	if expected != result {
		_, file, line, _ := runtime.Caller(1)
		t.Errorf("\n%v:%v: Expected %v, got %v for Contains(%v).", path.Base(file), line, expected, result, item)
	}
}

func assertIndexOf(t *testing.T, l ReadOnlyList, item interface{}, expected int) {
	result, ok := l.IndexOf(item)
	if !ok || expected != result {
		_, file, line, _ := runtime.Caller(1)
		t.Errorf("\n\t%v:%v: Expected %v, got %v for IndexOf(%v)", path.Base(file), line, expected, result, item)
	}
}
