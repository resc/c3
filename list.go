package c3

// Creates a new, empty list
func NewList() List {
	return newList()
}
func newList() *list {
	return &list{0, make([]interface{}, 0, 4)}
}

// Wraps a slice in a List interface
func WrapList(items []interface{}) List {
	return &list{0, items[:]}
}

// Creates a new, empty Collection
func NewCollection() Collection {
	return NewList()
}

// Creates a list with the given items.
func ListOf(items ...interface{}) List {
	l := len(items)
	x := make([]interface{}, l)
	if l != copy(x, items) {
		panic("Didn't copy all items into the list")
	}
	return &list{0, x}
}

// Creates a collection with the given items.
func IterableOf(items ...interface{}) Iterable {
	return ListOf(items...)
}

// Creates a collection with the given items.
func ReadOnlyCollectionOf(items ...interface{}) ReadOnlyCollection {
	return ListOf(items...)
}

// Creates a collection with the given items.
func ReadOnlyListOf(items ...interface{}) ReadOnlyList {
	return ListOf(items...)
}

// Creates a collection with the given items.
func CollectionOf(items ...interface{}) Collection {
	return ListOf(items...)
}

type listiter struct {
	version int
	index   int
	l       *list
}

func (i *listiter) MoveNext() bool {
	if i.version != i.l.version {
		panic("Concurrent modification detected")
	}

	if i.index < len(i.l.items) {
		i.index++
	}

	return i.index < len(i.l.items)
}

func (i *listiter) Value() interface{} {
	l := i.l
	index := i.index
	if i.version != l.version {
		panic("Concurrent modification detected")
	}

	if 0 > index || index >= l.Len() {
		return nil
	}

	return l.items[index]
}

type list struct {
	version int
	items   []interface{}
}

func (l *list) Iterator() Iterator {
	return &listiter{l.version, -1, l}
}

func (l *list) Add(item interface{}) bool {
	l.items = append(l.items, item)
	l.version++
	return true
}

func (l *list) InsertAt(index int, item interface{}) bool {
	if 0 > index || index > len(l.items) {
		return false
	}

	if index == len(l.items) {
		return l.Add(item)
	}

	l.items = append(l.items, nil)
	copy(l.items[index+1:], l.items[index:])
	l.items[index] = item
	l.version++
	return true
}

func (l *list) First() (interface{}, bool) {
	return l.Get(0)
}

func (l *list) Last() (interface{}, bool) {
	return l.Get(l.Len() - 1)
}

func (l *list) Get(index int) (interface{}, bool) {
	if 0 > index || index >= len(l.items) {
		return nil, false
	}
	return l.items[index], true
}

func (l *list) Contains(item interface{}) bool {
	_, ok := l.IndexOf(item)
	return ok
}

func (l *list) IndexOf(item interface{}) (int, bool) {
	return l.NextIndexOf(-1, item)
}

func (l *list) NextIndexOf(offset int, item interface{}) (int, bool) {
	for index := max(-1, offset) + 1; 0 <= index && index < l.Len(); index++ {
		if l.items[index] == item {
			return index, true
		}
	}
	return -1, false
}

func (l *list) LastIndexOf(item interface{}) (int, bool) {
	return l.PrevIndexOf(l.Len(), item)
}

func (l *list) PrevIndexOf(offset int, item interface{}) (int, bool) {
	for index := min(offset, l.Len()) - 1; 0 <= index && index < l.Len(); index-- {
		if l.items[index] == item {
			return index, true
		}
	}
	return -1, false
}

func (l *list) Delete(item interface{}) bool {
	if index, ok := l.IndexOf(item); ok {
		return l.DeleteAt(index)
	}
	return false
}

func (l *list) DeleteAt(index int) bool {
	if 0 > index || index >= len(l.items) {
		return false
	}
	last := len(l.items) - 1
	if index != last {
		copy(l.items[index:], l.items[index+1:])
	}
	l.items[last] = nil
	l.items = l.items[:last]
	l.version++
	return true
}

func (l *list) Len() int {
	return len(l.items)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
