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
