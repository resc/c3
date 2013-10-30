package c3

// The c3 query representation
type Q struct {
	result Iterable
}

// Action is invoked for every item in the query result.
type Action func(interface{})

// Predicate if a function that returns true if the predicate holds for the item.
type Predicate func(item interface{}) bool

// Aggregator converts an item and an aggregate into an aggregate result
type Aggregator func(item interface{}, aggregate interface{}) (aggregateResult interface{})

// Selector converts an item into another item
type Selector func(interface{}) interface{}

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
	return nil, false
}

// Last returns the last query result and true,
// or nil and false if there are no results.
func (q *Q) Last() (interface{}, bool) {
	value, ok := interface{}(nil), false
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

// Len counts the number of results.
func (q *Q) Len() int {
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
	return NewQuery(ListOf(q, items)).SelectMany(func(v interface{}) Iterable {
		return v.(Iterable)
	})
}

// Distinct filters non-unique items from the query result.
func (q *Q) Distinct() *Q {
	set := make(map[interface{}]bool)
	return q.Where(func(v interface{}) bool {
		if set[v] {
			return false
		}
		set[v] = true
		return true
	})
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
