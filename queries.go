package c3

import "math/rand"
import "time"

// The c3 query representation
type Q struct {
	result Iterable
}

// Action is invoked for every item in the query result.
type Action func(item interface{})

// Predicate if a function that returns true if the predicate holds for the item.
type Predicate func(item interface{}) bool

// Aggregator converts an item and an aggregate into an aggregate result
type Aggregator func(item interface{}, aggregate interface{}) (aggregateResult interface{})

// Selector converts an item into another item
type Selector func(item interface{}) interface{}

// Lesser compares 2 items, such that a<b:true, false otherwise
type Lesser func(a, b interface{}) bool

// ManySelector converts 1 item into zero or more items
type ManySelector func(interface{}) Iterable

// Iterator provides an iterator for the query results.
func (q *Q) Iterator() Iterator {
	return q.result.Iterator()
}

// Filters the items using the filter function.
// If filter returns true, the item is included
// in the result, otherwise it is skipped.
func (q *Q) Where(filter Predicate) *Q {
	return &Q{&whereIterable{q.result, filter}}
}

// Select uses the selector to create a new result for each item.
func (q *Q) Select(selector Selector) *Q {
	return &Q{&selectIterable{q.result, selector}}
}

// SelectMany uses the selector to create an Iteratable containing zero or more
// items for each item, and concatenates all the results.
func (q *Q) SelectMany(selector ManySelector) *Q {
	return &Q{&selectManyIterable{q.result, selector}}
}

// For applies the action to every item in the query result.
func (q *Q) For(action Action) {
	For(q, action)
}

// Run runs the query and discards the results.
func (q *Q) Run() {
	for i := q.Iterator(); i.MoveNext(); {
		// do nothing.
	}
}

// Go applies the action to every item in the query result on a
// seperate goroutine using an unbuffered channel.
func (q *Q) Go(action Action) {
	Go(q, action)
}

// Go applies the action to every item in the query result on a
// seperate goroutine using a buffered channel of the supplied size.
func (q *Q) GoBuffered(bufferSize int, action Action) {
	GoBuffered(q, bufferSize, action)
}

// ToSlice puts the query results in a new slice
func (q *Q) ToSlice() []interface{} {
	return ToSlice(q)
}

// ToList puts the query results in a new List
func (q *Q) ToList() List {
	result, ok := q.result.(List)
	if ok {
		return result
	}
	return ToList(q)
}

// ToReadOnlyList puts the query results in a new ReadOnlyList
func (q *Q) ToReadOnlyList() ReadOnlyList {
	return ToList(q)
}

// ToReadOnlyList puts the query results in a new ReadOnlyList
func (q *Q) ToReadOnlyBag() ReadOnlyBag {
	return ToList(q)
}

// ToBag puts the query results in a new Bag
func (q *Q) ToBag() Bag {
	return ToList(q)
}

// ToSet puts the unique query results in a new Set
func (q *Q) ToSet() Set {
	return ToSet(q)
}

// ToQueue puts the unique query results in a new Queue
func (q *Q) ToQueue() Queue {
	return ToQueue(q)
}

// ToStack puts the unique query results in a new Stack
func (q *Q) ToStack() Stack {
	return ToStack(q)
}

// Aggregate applies the action to every item in the query result
// and combines them in a single result.
func (q *Q) Aggregate(aggregate interface{}, action Aggregator) interface{} {
	for i := q.Iterator(); i.MoveNext(); {
		aggregate = action(i.Value(), aggregate)
	}
	return aggregate
}

// First returns the first query result and true,
// or nil and false if there are no results.
func (q *Q) First() (interface{}, bool) {
	for i := q.Iterator(); i.MoveNext(); {
		return i.Value(), true
	}
	return defaultElementValue, false
}

// Last returns the last query result and true,
// or nil and false if there are no results.
func (q *Q) Last() (interface{}, bool) {
	value, ok := defaultElementValue, false
	for i := q.Iterator(); i.MoveNext(); {
		value, ok = i.Value(), true
	}
	return value, ok
}

// Any returns true if there are results, false if there are not any results.
func (q *Q) Any() bool {
	for i := q.Iterator(); i.MoveNext(); {
		return true
	}
	return false
}

// All returns true if the predicate holds for all results, false otherwise.
func (q *Q) All(predicate Predicate) bool {
	for i := q.Iterator(); i.MoveNext(); {
		if !predicate(i.Value()) {
			return false
		}
	}
	return true
}

// Contains returns true if the query results contain the item, false otherwise.
func (q *Q) Contains(item interface{}) bool {
	return q.Where(func(x interface{}) bool {
		return x == item
	}).Any()
}

// Count counts the number of results.
func (q *Q) Count() int {
	count := 0
	for i := q.Iterator(); i.MoveNext(); {
		count++
	}
	return count
}

// Take truncates the results after count results have been computed.
// If there are less results Take returns only the available results.
func (q *Q) Take(count int) *Q {
	taken := 0
	return q.Where(func(v interface{}) bool {
		if taken >= count {
			return false
		}
		taken++
		return taken <= count
	})
}

// Prepend prepends the items to the query result.
func (q *Q) Prepend(items ...interface{}) *Q {
	return NewQuery(ListOf(items...)).Concat(q)
}

// Appens appends the items to the query result.
func (q *Q) Append(items ...interface{}) *Q {
	return q.Concat(ListOf(items...))
}

// Concat appends the items to the query result.
func (q *Q) Concat(items Iterable) *Q {
	return NewQuery(&concatIterable{q, items})
}

// Distinct filters non-unique items from the query result.
func (q *Q) Distinct() *Q {
	set := make(map[interface{}]bool)
	return q.Where(func(v interface{}) bool {
		if !set[v] {
			set[v] = true
			return true
		}

		return false
	})
}

type concatIterable struct {
	a, b Iterable
}

func (i *concatIterable) Iterator() Iterator {
	return &concatIterator{i.a.Iterator(), i.b.Iterator(), defaultElementValue}
}

type concatIterator struct {
	a, b  Iterator
	value interface{}
}

func (i *concatIterator) Value() interface{} {
	return i.value
}

func (i *concatIterator) MoveNext() bool {
	if i.a.MoveNext() {
		i.value = i.a.Value()
		return true
	}
	if i.b.MoveNext() {
		i.value = i.b.Value()
		return true
	}
	i.value = defaultElementValue
	return false
}

// Tee applies the action to every item in the query result
// and passes every result on to the next query operator.
func (q *Q) Tee(action Action) *Q {
	return q.Where(func(v interface{}) bool {
		action(v)
		return true
	})
}

// Skip skips the first count results and returns all results after that.
// If there are less results Skip returns an empty result set.
func (q *Q) Skip(count int) *Q {
	skipped := 0
	return q.Where(func(v interface{}) bool {
		if skipped >= count {
			return true
		}
		skipped++
		return skipped > count
	})
}

// Sort sorts the result set using the lesser function.
func (q *Q) Sort(lesser Lesser) *Q {
	l := q.ToList()
	Sort(l, lesser)
	return &Q{l}
}

// Shuffle randomizes the order of the result set.
func (q *Q) Shuffle() *Q {
	return &Q{MakeIterable(func() Generate {
		// make and fill the shuffle buffer
		bufCap := 32
		buf := make([]interface{}, 0, bufCap)
		i := q.Iterator()
		for len(buf) < bufCap && i.MoveNext() {
			buf = append(buf, i.Value())
		}

		// setup the iterator state
		shuffleDone := len(buf) == 0
		sourceDone := len(buf) < bufCap
		if shuffleDone {
			// just return
			return func() (interface{}, bool) {
				return defaultElementValue, false
			}
		}

		// setup the randomizer
		rndSeed := time.Now().UnixNano()
		rndSource := rand.NewSource(rndSeed)
		rnd := rand.New(rndSource)

		// return the Generate function.
		return func() (interface{}, bool) {

			if shuffleDone {
				return defaultElementValue, false
			}

			// get the value from the buffer
			bufLen := len(buf)
			index := rnd.Intn(bufLen)
			value := buf[index]

			// get the next value from the source
			if !sourceDone {
				if i.MoveNext() {
					buf[index] = i.Value()
					return value, true
				} else {
					sourceDone = true
				}
			}

			// move the last buffer value into the current value's location
			lastIndex := bufLen - 1
			// if index==lastIndex this is an expensive NOP
			buf[index] = buf[lastIndex]

			// clear the last value's location to help the
			// garbage collector and then shorten the buffer
			buf[lastIndex] = defaultElementValue
			buf = buf[:lastIndex]

			shuffleDone = len(buf) == 0
			return value, true
		}
	})}
}
