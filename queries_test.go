package c3

import "testing"

func TestWhere(t *testing.T) {
	l := ListOf(1, 2, 3, 4)
	result := Query(l).Where(isMod2).ToList()

	if result.Len() != 2 {
		t.Error("Expected 2 elements in the list")
	}
}

func TestSkip(t *testing.T) {
	l := ListOf(1, 2, 3, 4)

	result := Query(l).Skip(0).ToList()
	if result.Len() != 4 {
		t.Errorf("Expected 4 elements in the list got %v", result.Len())
	}

	result = Query(l).Skip(1).ToList()
	if result.Len() != 3 {
		t.Errorf("Expected 3 elements in the list got %v", result.Len())
	}

	result = Query(l).Skip(4).ToList()
	if result.Len() != 0 {
		t.Errorf("Expected 0 elements in the list, got %v", result.Len())
	}

	result = Query(l).Skip(5).ToList()
	if result.Len() != 0 {
		t.Errorf("Expected 0 elements in the list, got %v", result.Len())
	}
}

func TestSelectMany(t *testing.T) {
	l := ListOf(1, 2, 3, 4)
	result := Query(l).
		SelectMany(func(v interface{}) Iterable {
		x := v.(int)
		return Range(x*10, x*10+9)
	}).
		ToList()

	if result.Len() != 40 {
		t.Errorf("Expected 40 elements, got %v", ToSlice(result))
	}
}

func TestDistinctNil(t *testing.T) {
	nils := Query(Repeat(3, nil)).Distinct().ToSlice()
	if len(nils) != 1 {
		t.Error("Expected only 1 item in the result")
	}
}

func TestDistinctIntAppend(t *testing.T) {
	set := Query(Repeat(3, 42)).Append(2, 2, 2).Distinct().ToSlice()
	if len(set) != 2 {
		t.Errorf("Expected only 2 items in the result %v", set)
	}
}

func TestDistinctIntPrepend(t *testing.T) {
	set := Query(Repeat(3, 42)).Prepend(2, 2, 2).Distinct().ToSlice()
	if len(set) != 2 {
		t.Errorf("Expected only 2 items in the result %v", set)
	}
}

func TestTake(t *testing.T) {
	l := ListOf(1, 2, 3, 4)

	result := Query(l).Take(0).ToList()
	if result.Len() != 0 {
		t.Errorf("Expected 0 elements in the list got %v", result.Len())
	}

	result = Query(l).Take(1).ToList()
	if result.Len() != 1 {
		t.Errorf("Expected 1 element in the list got %v", result.Len())
	}

	result = Query(l).Take(4).ToList()
	if result.Len() != 4 {
		t.Errorf("Expected 4 elements in the list, got %v", result.Len())
	}

	result = Query(l).Take(5).ToList()
	if result.Len() != 4 {
		t.Errorf("Expected 4 elements in the list, got %v", result.Len())
	}
}

func isMod2(v interface{}) bool {
	return v.(int)%2 == 0
}
