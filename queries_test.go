package c3

import "testing"

func TestWhere(t *testing.T) {
	l := ListOf(1, 2, 3, 4)
	result := NewQuery(l).Where(isMod2).ToList()

	if result.Len() != 2 {
		t.Error("Expected 2 elements in the list")
	}
}

func TestSkip(t *testing.T) {
	l := ListOf(1, 2, 3, 4)

	result := NewQuery(l).Skip(0).ToList()
	if result.Len() != 4 {
		t.Errorf("Expected 4 elements in the list got %v", result.Len())
	}

	result = NewQuery(l).Skip(1).ToList()
	if result.Len() != 3 {
		t.Errorf("Expected 3 elements in the list got %v", result.Len())
	}

	result = NewQuery(l).Skip(4).ToList()
	if result.Len() != 0 {
		t.Errorf("Expected 0 elements in the list, got %v", result.Len())
	}

	result = NewQuery(l).Skip(5).ToList()
	if result.Len() != 0 {
		t.Errorf("Expected 0 elements in the list, got %v", result.Len())
	}
}

func TestSelectMany(t *testing.T) {
	l := ListOf(1, 2, 3, 4)
	result := NewQuery(l).
		SelectMany(func(v interface{}) Iterable {
		x := v.(int)
		return Range(x*10, x*10+9)
	}).
		ToList()

	if result.Len() != 40 {
		t.Errorf("Expected 40 elements, got %v", ToSlice(result))
	}
}

func TestCount(t *testing.T) {
	n := 3
	q := NewQuery(Repeat(n, nil))

	c1 := q.Count()
	if c1 != n {
		t.Errorf("Expected %v, got %v", n, c1)
	}

	c2 := q.Count()
	if c2 != n {
		t.Errorf("Expected %v, got %v", n, c2)
	}
}

/*
func TestShuffles(t *testing.T) {
	slice := NewQuery(Range(0, 15)).Shuffle().ToSlice()
	t.Error(slice)
}
*/

func BenchmarkShuffle1000(b *testing.B) {
	var q = NewQuery(testSource(nil))
	shuffle := q.Shuffle()
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		shuffle.Run()
	}
}

func TestEmptyShuffle(t *testing.T) {
	n := 0
	q := NewQuery(EmptyIterable())

	result := q.Shuffle().ToList()

	if result.Len() != n {
		t.Errorf("Expected %v results, got %v", n, result.Len())
	}
}

func TestNonEmptyShuffle(t *testing.T) {
	n := 10
	q := NewQuery(Range(0, n-1))

	result := q.Shuffle().ToList()

	if result.Len() != n {
		t.Errorf("Expected %v results, got %v", n, result.Len())
	}
}

func TestSort(t *testing.T) {
	l := NewQuery(testSource(t)).Shuffle().ToList()
	if l.Len() != 1000 {
		t.Error("wrong count")
	}

	l = NewQuery(l).Sort(func(a, b interface{}) bool {
		return a.(int) < b.(int)
	}).ToList()

	if l.Len() != 1000 {
		t.Error("wrong count")
	}

	var index = 0
	for i := l.Iterator(); i.MoveNext(); {
		if i.Value().(int) != index {
			t.Errorf("Expected %v, got %v", index, ToSlice(l))
			return
		}
		index++
	}
}

func TestShuffle(t *testing.T) {
	// smaller,equal and larger number of items
	// than the internal buffer of Shuffle uses.
	for n := 16; n <= 64; n = n * 2 {
		q := NewQuery(Range(0, n-1))

		result := q.Shuffle().ToList()

		if result.Len() != n {
			t.Errorf("Expected %v results, got %v", n, result.Len())
		}

		different := false
		for i := 0; i < result.Len(); i++ {
			val, _ := result.Get(i)
			different = val.(int) != i
			if different {
				// good enough for me...
				break
			}
		}

		if !different {
			t.Error("Expected a randomized result, got: ", ToSlice(result))
		}
	}
}

func TestDistinctNil(t *testing.T) {
	set := NewQuery(Repeat(3, nil)).Distinct().ToSlice()
	if len(set) != 1 {
		t.Errorf("Expected 1 item in the result, got %v", set)
	}
}

func TestDistinctIntAppend(t *testing.T) {
	set := NewQuery(Repeat(3, 42)).Append(2, 2, 2).Distinct().ToSlice()
	if len(set) != 2 {
		t.Errorf("Expected 2 items in the result %v", set)
	}
}

func TestDistinctIntPrepend(t *testing.T) {
	set := NewQuery(Repeat(3, 42)).Prepend(2, 2, 2).Distinct().ToSlice()
	if len(set) != 2 {
		t.Errorf("Expected 2 items in the result %v", set)
	}
}

func TestTake(t *testing.T) {
	l := ListOf(1, 2, 3, 4)

	result := NewQuery(l).Take(0).ToList()
	if result.Len() != 0 {
		t.Errorf("Expected 0 elements in the list got %v", result.Len())
	}

	result = NewQuery(l).Take(1).ToList()
	if result.Len() != 1 {
		t.Errorf("Expected 1 element in the list got %v", result.Len())
	}

	result = NewQuery(l).Take(4).ToList()
	if result.Len() != 4 {
		t.Errorf("Expected 4 elements in the list, got %v", result.Len())
	}

	result = NewQuery(l).Take(5).ToList()
	if result.Len() != 4 {
		t.Errorf("Expected 4 elements in the list, got %v", result.Len())
	}
}

func isMod2(v interface{}) bool {
	return v.(int)%2 == 0
}

func testSource(t *testing.T) Iterable {
	n := 1000
	source := ToQueue(Range(0, n-1))
	if t != nil {
		if source.Len() != n {
			panic("No items!")
		}
		index := 0
		for i := source.Iterator(); i.MoveNext(); {
			val, ok := i.Value().(int)
			if !ok || val != index {
				t.Errorf("unexpected sequence expected %v got %v", index, val)
			}
			index++
		}
	}
	return source
}
