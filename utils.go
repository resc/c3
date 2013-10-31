package c3

var (
	// An empty Iterable
	emptyIterable Iterable = &nilIterable{}
	// An empty Iterator.
	emptyIterator Iterator = &nilIterator{}
)

func EmptyIterable() Iterable {
	return emptyIterable
}

func EmptyIterator() Iterator {
	return emptyIterator
}

type nilIterable struct{}

func (x *nilIterable) Iterator() Iterator {
	return emptyIterator
}

type nilIterator struct{}

func (x *nilIterator) MoveNext() bool { return false }

func (x *nilIterator) Value() interface{} { return nil }

// Sort sorts the list with the given Lesser function
func Sort(l List, lesser Lesser) {
	s := &Sorter{l, lesser}
	s.Sort()
}

// ToSlice makes a new slice of the items in an Iterable
func ToSlice(c Iterable) []interface{} {
	var slice []interface{}
	if col, ok := c.(ReadOnlyBag); ok {
		slice = make([]interface{}, 0, col.Len())
	} else {
		slice = make([]interface{}, 0, 4)
	}
	for i := c.Iterator(); i.MoveNext(); {
		slice = append(slice, i.Value())
	}
	return slice
}

// ToList creates a new List of the items in an Iterable
func ToList(c Iterable) List {
	l := NewList()
	for i := c.Iterator(); i.MoveNext(); {
		l.Add(i.Value())
	}
	return l
}

// ToReadOnlyList creates a new ReadOnlyList of the items in an Iterable
func ToReadOnlyList(c Iterable) ReadOnlyList {
	l := NewList()
	for i := c.Iterator(); i.MoveNext(); {
		l.Add(i.Value())
	}
	return l
}

// ToBag creates a new Bag of the items in an Iterable
func ToBag(c Iterable) Bag {
	l := NewList()
	for i := c.Iterator(); i.MoveNext(); {
		l.Add(i.Value())
	}
	return l
}

// ToReadOnlyBag creates a new ReadOnlyBag of the items in an Iterable
func ToReadOnlyBag(c Iterable) ReadOnlyBag {
	l := NewList()
	for i := c.Iterator(); i.MoveNext(); {
		l.Add(i.Value())
	}
	return l
}

// ToSet makes a new Set of the unique items in an Iterable
func ToSet(c Iterable) Set {
	set := NewSet()
	for i := c.Iterator(); i.MoveNext(); {
		set.Add(i.Value())
	}
	return set
}

// ToStack makes a new Stack of the items in an Iterable
func ToStack(c Iterable) Stack {
	s := NewStack()
	for i := c.Iterator(); i.MoveNext(); {
		s.Push(i.Value())
	}
	return s
}

// ToQueue makes a new Queue of the items in an Iterable
func ToQueue(c Iterable) Queue {
	s := NewQueue()
	for i := c.Iterator(); i.MoveNext(); {
		s.Enqueue(i.Value())
	}
	return s
}

// For applies the action to each item in the Iterable
func For(c Iterable, action Action) {
	for i := c.Iterator(); i.MoveNext(); {
		action(i.Value())
	}
}

//  Go applies the action to each item in the Iterable
//  on a separate goroutine using an unbuffered channel
func Go(c Iterable, action Action) {
	ch := make(chan interface{})
	defer close(ch)
	go func() {
		for value := range ch {
			action(value)
		}
	}()
	For(c, func(value interface{}) { ch <- value })
}

// GoBuffered applies the action to each item in the Iterable
// on a separate goroutine using a buffered channel
func GoBuffered(c Iterable, bufferSize int, action Action) {
	ch := make(chan interface{}, bufferSize)
	defer close(ch)
	go func() {
		for value := range ch {
			action(value)
		}
	}()
	For(c, func(value interface{}) { ch <- value })
}

// Repeat repeats the item count times.
//
// e.g.:
//		Repeat(3,42) // returns [42,42,42]
func Repeat(count int, item interface{}) Iterable {
	if count < 0 {
		panic("Count parameter invalid")
	}
	if count == 0 {
		return emptyIterable
	}
	return MakeIterable(func() Generate {
		x := 0
		return func() (interface{}, bool) {
			if x == count {
				return nil, false
			}
			x++
			return item, true
		}
	})
}

// Range creates a range iterable that iterates over the ints from start to end inclusive.
//
// e.g.:
//		 Range(0,9) // returns [0,1,2,3,4,5,6,7,8,9]
//		 Range(9,0) // returns [9,8,7,6,5,4,3,2,1,0]
func Range(start, end int) Iterable {
	inc := 1
	if end < start {
		inc = -1
	}
	return MakeIterable(func() Generate {
		x := start - inc
		y := end
		return func() (interface{}, bool) {
			if x == y {
				return nil, false
			}
			x = x + inc
			return x, true
		}
	})
}

func max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}
