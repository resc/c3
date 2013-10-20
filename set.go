package c3

// NewSet creates a new, empty Set
func NewSet() Set {
	return &set{make(map[interface{}]bool)}
}

// SetOf creates a new Set containing the unique items
func SetOf(items ...interface{}) Set {
	set := NewSet()
	for item := range items {
		set.Add(item)
	}
	return set
}

type set struct {
	items map[interface{}]bool
}

func (s *set) Add(item interface{}) bool {
	if s.Contains(item) {
		return false
	}
	s.items[item] = true
	return true
}

func (s *set) Contains(item interface{}) bool {
	return s.items[item]
}

func (s *set) Delete(item interface{}) bool {
	if s.Contains(item) {
		delete(s.items, item)
		return true
	}
	return false
}

func (s *set) Len() int {
	return len(s.items)
}

func (s *set) Iterator() Iterator {
	// TODO optimize this, it's horrible...
	// TODO check version.
	var items = make([]interface{}, 0, s.Len())
	for k, _ := range s.items {
		items = append(items, k)
	}
	return ToIterable(func() Generate {
		iter := items[:]
		return func() (interface{}, bool) {
			if len(iter) == 0 {
				return nil, false
			}
			item := iter[0]
			iter = iter[1:]
			return item, true
		}
	}).Iterator()
}

func (s *set) Intersection(other Set) Set {
	result := NewSet()
	os, ok := other.(*set)
	if ok {
		// fast path
		for item, _ := range os.items {
			if s.Contains(item) {
				result.Add(item)
			}
		}
	} else {
		// slow path
		for i := other.Iterator(); i.MoveNext(); {
			if s.Contains(i.Value()) {
				result.Add(i.Value())
			}
		}
	}
	return result
}

func (s *set) Difference(other Set) Set {
	result := NewSet()
	for item, _ := range s.items {
		if !other.Contains(item) {
			result.Add(item)
		}
	}
	return result
}

func (s *set) Union(other Set) Set {
	result := NewSet()
	for item, _ := range s.items {
		result.Add(item)
	}

	os, ok := other.(*set)
	if ok {
		// fast path
		for item, _ := range os.items {
			result.Add(item)
		}
	} else {
		// slow path
		for i := other.Iterator(); i.MoveNext(); {
			result.Add(i.Value())
		}
	}
	return result
}

func (s *set) SymmetricDifference(other Set) Set {
	result := NewSet()
	for item, _ := range s.items {
		if !other.Contains(item) {
			result.Add(item)
		}
	}

	os, ok := other.(*set)
	if ok {
		// fast path
		for item, _ := range os.items {
			if !s.Contains(item) {
				result.Add(item)
			}
		}
	} else {
		// slow path
		for i := other.Iterator(); i.MoveNext(); {
			if !s.Contains(i.Value()) {
				result.Add(i.Value())
			}
		}
	}
	return result
}
