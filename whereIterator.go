package c3

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
