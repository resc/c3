package c3

// Query provides the entry point to the c3 query api.
//
// Usage:
//		list := c3.ListOf(1,2,3)
//		q := Query(list).
//			Where(/* filter function here */).
//			Select( /* selector function here */).
//			ToList() /* collect the results */
func Query(items Iterable) *Q {
	return &Q{items}
}

// The c3 query representation
type Q struct {
	result Iterable
}

// Iterator provides an iterator for the query results.
func (q *Q) Iterator() Iterator {
	return q.result.Iterator()
}

// Action is invoked for every item in the query result.
type Action func(interface{})

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

// ToList  puts the query results in a new List
func (q *Q) ToList() List {
	l := NewList()
	for i := q.Iterator(); i.MoveNext(); {
		l.Add(i.Value())
	}
	return l
}

// ToSlice puts the query results in a new slice
func (q *Q) ToSlice() []interface{} {
	return ToSlice(q)
}

// ToSet puts the unique query results in a new set
func (q *Q) ToSet() Set {
	return ToSet(q)
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

// Len counts the number of results.
func (q *Q) Len() int {
	count := 0
	for i := q.Iterator(); i.MoveNext(); {
		count++
	}
	return count
}

// Any returns true if there are results, false if there are not any results.
func (q *Q) Any() bool {
	for i := q.Iterator(); i.MoveNext(); {
		return true
	}
	return false
}

// Predicate if a function that returns true if the predicate holds for the item.
type Predicate func(item interface{}) bool

// All returns true if the predicate holds for all results, false otherwise.
func (q *Q) All(predicate Predicate) bool {
	for i := q.Iterator(); i.MoveNext(); {
		if !predicate(i.Value()) {
			return false
		}
	}
	return true
}

// Prepend prepends the items to the query result.
func (q *Q) Prepend(items ...interface{}) *Q {
	return Query(ListOf(items...)).Concat(q)
}

// Appens appends the items to the query result.
func (q *Q) Append(items ...interface{}) *Q {
	return q.Concat(ListOf(items...))
}

// Concat appends the items to the query result.
func (q *Q) Concat(items Iterable) *Q {
	return Query(ListOf(q, items)).SelectMany(func(v interface{}) Iterable {
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

// Aggregator converts an item and an aggregate into an aggregate result
type Aggregator func(item interface{}, aggregate interface{}) (aggregateResult interface{})

// Aggregate applies the action to every item in the query result
// and combines them in a single result.
func (q *Q) Aggregate(aggregate interface{}, action Aggregator) interface{} {
	for i := q.Iterator(); i.MoveNext(); {
		aggregate = action(i.Value(), aggregate)
	}
	return aggregate
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

// Filters the items using the filter function.
// If filter returns true, the item is included
// in the result, otherwise it is skipped.
func (q *Q) Where(filter Predicate) *Q {
	return &Q{&whereIterable{q.result, filter}}
}

type whereIterable struct {
	items Iterable
	where Predicate
}

func (i *whereIterable) Iterator() Iterator {
	return &whereIterator{i.items.Iterator(), i.where}
}

type whereIterator struct {
	items Iterator
	where Predicate
}

func (i *whereIterator) MoveNext() bool {
	for i.items.MoveNext() {
		if i.where(i.items.Value()) {
			return true
		}
	}
	return false
}

func (i *whereIterator) Value() interface{} {
	return i.items.Value()
}

// ManySelector converts 1 item into many items
type ManySelector func(interface{}) Iterable

// SelectMany uses the selector to create an Iterator for each
// item, and concatenates all the results in a single flat result set.
func (q *Q) SelectMany(selector ManySelector) *Q {
	return &Q{&selectManyIterable{q.result, selector}}
}

type selectManyIterable struct {
	items    Iterable
	selector ManySelector
}

func (i *selectManyIterable) Iterator() Iterator {
	return &selectManyIterator{
		i.items.Iterator(),
		i.selector,
		EmptyIterator,
		nil,
	}
}

type selectManyIterator struct {
	items    Iterator
	selector ManySelector
	iterator Iterator
	value    interface{}
}

func (i *selectManyIterator) MoveNext() bool {
	if i.iterator.MoveNext() {
		i.value = i.iterator.Value()
		return true
	}
	for i.items.MoveNext() {
		value := i.items.Value()
		i.iterator = i.selector(value).Iterator()
		if !i.iterator.MoveNext() {
			continue
		}
		i.value = i.iterator.Value()
		return true
	}
	i.iterator = EmptyIterator
	i.value = nil
	return false
}

func (i *selectManyIterator) Value() interface{} {
	return i.value
}

// Selector converts an item into another item
type Selector func(interface{}) interface{}

// Select uses the selector to create a new result for each item.
func (q *Q) Select(selector Selector) *Q {
	return &Q{&selectIterable{q.result, selector}}
}

type selectIterable struct {
	items    Iterable
	selector Selector
}

func (i *selectIterable) Iterator() Iterator {
	return &selectIterator{i.items.Iterator(), i.selector, nil}
}

type selectIterator struct {
	items    Iterator
	selector Selector
	value    interface{}
}

func (i *selectIterator) MoveNext() bool {
	if i.items.MoveNext() {
		i.value = i.selector(i.items.Value())
		return true
	}
	i.value = nil
	return false
}

func (i *selectIterator) Value() interface{} {
	return i.value
}
