// c3 stands for common containers collection and
// provides a few simple to use containers.
//
// The containers provided:
//		- Iterable: a container that allow iterating over its items.
//		- Bag: an unordered container that allows duplicate items.
//		- Set: an unordered container that does not allow duplicate items.
//		- List: an indexable container items.
//		- Queue: an fifo container.
//		- Stack: an lifo container.
//
// It also provides a query api for those containers that looks like C#'s Linq
package c3

// Iterator provides a way to iterate over a container
//
// Usage:
//		for i := iterable.Iterator(); i.MoveNext(); {
//			value := i.Value()
//		}
// Iterator.Value() only return a valid value if the preceding call to
// Iterator.MoveNext() returned true
type Iterator interface {
	// Move the iterator to the next item, returns true on succes
	// or false if the are no more item
	MoveNext() bool
	// returns the value at the current iterator position, or nil
	Value() interface{}
}

// Iterable is a container of items that can be iterated over.
type Iterable interface {
	// returns a new Iterator positioned at the start of the container.
	Iterator() Iterator
}

/*
type Equatable interface {
	Equals(other interface{}) bool
	Hashcode() int
}

type Comparable interface {
	Compare(other interface{}) int
}
*/

// ReadOnlyBag provides a length and a test for
// determining if an item is present in a container
type ReadOnlyBag interface {
	Iterable
	// Returns the item count
	Len() int
	// Returns true if the item is in the container, false otherwise.
	Contains(item interface{}) bool
}

// Indexable provides methods to get items from a container by index
type Indexable interface {
	// Returns the first item and true,
	// or nil and false if there is no first item
	First() (interface{}, bool)
	// Returns the item at the index and true,
	// or nil and false if the index is out of bounds
	Get(index int) (interface{}, bool)
	// Returns the last item and true,
	// or nil and false if there is no last item
	Last() (interface{}, bool)
	// Returns the first index of the item and true,
	// or -1 and false if there is no such item
	IndexOf(item interface{}) (int, bool)
	// Returns the index of the next item before the offset and true,
	// or -1 and false if there is no such item
	PrevIndexOf(offset int, item interface{}) (int, bool)
	// Returns the index of the next item after the offset and true,
	// or -1 and false if there is no such item
	NextIndexOf(offset int, item interface{}) (int, bool)
	// Returns the last index of the item and true,
	// or -1 and false if there is no such item
	LastIndexOf(item interface{}) (int, bool)
}

// ReadOnlyList is an indexable readonly list
type ReadOnlyList interface {
	ReadOnlyBag
	Indexable
}

// Bag is an unordered mutable container
type Bag interface {
	ReadOnlyBag
	// Add adds an item to the container,
	// returns true if the container was modified,
	// false if it was not modified
	Add(item interface{}) bool
	// Delete removes an item from the container,
	// returns true if the container was modified,
	// false if it was not modified
	Delete(item interface{}) bool
}

type List interface {
	Bag
	Indexable
	// Inserts the item at the given index,
	// returns true if the container was modified,
	// false if it was not modified.
	InsertAt(index int, item interface{}) bool
	// Deletes the item at the given index,
	// returns true if the container was modified,
	// false if it was not modified.
	DeleteAt(index int) bool
}

// A set type with basic set operations.
// See also http://en.wikipedia.org/wiki/Set_theory
type Set interface {
	Bag
	// Union computes the union of the set.
	// i.e. all items in this set and the other set, without duplicates
	Union(other Set) Set
	// SymmetricDifference computes the symmetric difference of the sets.
	// i.e. all the items that are either in this set,
	// or in the other set, but not in both.
	SymmetricDifference(other Set) Set
	// Difference computes the items that are in this set but not in the other set.
	Difference(other Set) Set
	// Intersection computes the items that are present in both sets.
	Intersection(other Set) Set
}

// Peeker provides a method to look at the next item without removing it from the container.
type Peeker interface {
	// Peek returns the next item without removing it from the container.
	Peek() (interface{}, bool)
}

// Consume removes and returns the next item in a sequence.
// Returns the next item and true or nil and false if there are no more items.
type Consume func() (interface{}, bool)

// Generator creates a new Generate function
type Generator func() Generate

// Generate computes the next item in a sequence.
// Returns the next item and true or nil and false if there are no more items.
type Generate func() (interface{}, bool)

// Consumer is an Iterator that removes items from a container.
type Consumer interface {
	Iterator
}

// Consumable is a container that provides a consuming iterator.
type Consumable interface {
	Consumer() Consumer
}

// A simple queuing container.
type Queue interface {
	ReadOnlyBag
	Peeker
	Consumable
	// Appends an item at the tail of the queue,
	// returns true if the queue was modified,
	// false if it was not modified.
	Enqueue(item interface{}) bool
	// Removes an item from the head of the queue,
	// returns the item and true if the queue was modified,
	// or nil and false if it was not modified.
	Dequeue() (interface{}, bool)
}

// A simple stack container
type Stack interface {
	ReadOnlyBag
	Peeker
	Consumable
	// Adds an item at the top of the stack,
	// returns true if the stack was modified,
	// false if it was not modified.
	Push(item interface{}) bool
	// Removes an item from the top of the stack,
	// returns the item and true if the stack was modified,
	// nil and false if it was not modified.
	Pop() (interface{}, bool)
}
