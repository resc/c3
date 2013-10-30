package c3

// NewList creates a new, empty List.
func NewList() List {
	return newList()
}

// NewBag creates a new, empty Bag.
func NewBag() Bag {
	return NewList()
}

// NewQueue creates a new, empty Queue.
func NewQueue() Queue {
	return &queue{nil, nil, 0, 0}
}

// NewSet creates a new, empty Set.
func NewSet() Set {
	return &set{0, make(map[interface{}]bool)}
}

// NewStack creates a new, empty Stack.
func NewStack() Stack {
	return &stack{newList()}
}

// NewQuery provides a entry point to the c3 query api.
//
// Usage:
//		list := c3.ListOf(1,2,3)
//		q := Query(list).
//			Where(/* filter function here */).
//			Select( /* selector function here */).
//			ToList() /* collect the results */
func NewQuery(items Iterable) *Q {
	return &Q{items}
}

// QueryOf provides a entry point to the c3 query api.
//
// Usage:
//		q := c3.QueryOf(1,2,3).
//			    Where(/* filter function here */).
//			    Select( /* selector function here */).
//			    ToList() /* collect the results */
func QueryOf(items ...interface{}) *Q {
	return NewQuery(IterableOf(items...))
}

// StackOf creates a new Stack with the given items.
func StackOf(items ...interface{}) Stack {
	s := NewStack()
	for item := range items {
		s.Push(item)
	}
	return s
}

// SetOf creates a new Set containing the unique items.
func SetOf(items ...interface{}) Set {
	set := NewSet()
	for item := range items {
		set.Add(item)
	}
	return set
}

// QueueOf creates a Queue with the given items.
func QueueOf(items ...interface{}) Queue {
	q := NewQueue()
	for item := range items {
		q.Enqueue(item)
	}
	return q
}

// ListOf creates a List with the given items.
func ListOf(items ...interface{}) List {
	l := len(items)
	x := make([]interface{}, l)
	copy(x, items)
	return &list{0, x}
}

// MakeIterable converts the Generator function into an Iterable.
func MakeIterable(g Generator) Iterable {
	return &generatorIterable{g}
}

// MakeIterator converts the Generate function into an Iterator.
func MakeIterator(g Generate) Iterator {
	return &generateIterator{g, nil}
}

// IterableOf creates an Iterable with the given items.
func IterableOf(items ...interface{}) Iterable {
	return ListOf(items...)
}

// IteratorOf creates an Iterator with the given items.
func IteratorOf(items ...interface{}) Iterator {
	return ListOf(items...).Iterator()
}

// ReadOnlyBagOf creates a ReadOnlyBag with the given items.
func ReadOnlyBagOf(items ...interface{}) ReadOnlyBag {
	return ListOf(items...)
}

// ReadOnlyListOf creates a ReadOnlyList with the given items.
func ReadOnlyListOf(items ...interface{}) ReadOnlyList {
	return ListOf(items...)
}

// BagOf creates a Bag with the given items.
func BagOf(items ...interface{}) Bag {
	return ListOf(items...)
}

func newList() *list {
	return &list{0, make([]interface{}, 0, 4)}
}
