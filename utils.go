package c3

// Makes a slice of the items in an Iterable
func ToSlice(c Iterable) []interface{} {
	var slice []interface{}
	if col, ok := c.(ReadOnlyCollection); ok {
		slice = make([]interface{}, 0, col.Len())
	} else {
		slice = make([]interface{}, 0, 4)
	}
	for i := c.Iterator(); i.MoveNext(); {
		slice = append(slice, i.Value())
	}
	return slice
}

// Makes a set of the unique items in an Iterable
func ToSet(c Iterable) Set {
	set := NewSet()
	for i := c.Iterator(); i.MoveNext(); {
		set.Add(i.Value())
	}
	return set
}

// Applies the action to each item in the Iterable
func For(c Iterable, action func(interface{})) {
	for i := c.Iterator(); i.MoveNext(); {
		action(i.Value())
	}
}

//  Applies the action to each item in the Iterable
//  on a separate goroutine using an unbuffered channel
func Go(c Iterable, action func(interface{})) {
	ch := make(chan interface{})
	defer close(ch)
	go func() {
		for value := range ch {
			action(value)
		}
	}()
	For(c, func(value interface{}) { ch <- value })
}

// Applies the action to each item in the Iterable
// on a separate goroutine using a buffered channel
func GoBuffered(c Iterable, bufferSize int, f func(interface{})) {
	ch := make(chan interface{}, bufferSize)
	defer close(ch)
	go func() {
		for value := range ch {
			f(value)
		}
	}()
	For(c, func(value interface{}) { ch <- value })
}

// Converts the Generator function into an Iterable
func ToIterable(g Generator) Iterable {
	return &genIterable{g}
}

type genIterable struct {
	g Generator
}

func (i *genIterable) Iterator() Iterator {
	return &genIterator{i.g(), nil}
}

type genIterator struct {
	g     Generate
	value interface{}
}

func (i *genIterator) MoveNext() bool {
	value, ok := i.g()
	if ok {
		i.value = value
		return true
	} else {
		i.value = nil
		return false
	}
}

func (i *genIterator) Value() interface{} {
	return i.value
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
		return EmptyIterable
	}
	return ToIterable(func() Generate {
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

// Creates a range iterable that iterates over the ints from start to end inclusive.
//
// e.g.:
//		 Range(0,9) // returns [0,1,2,3,4,5,6,7,8,9]
//		 Range(9,0) // returns [9,8,7,6,5,4,3,2,1,0]
func Range(start, end int) Iterable {
	inc := 1
	if end < start {
		inc = -1
	}
	return ToIterable(func() Generate {
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
