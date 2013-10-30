package c3

type whereIterable struct {
	items Iterable
	where Predicate
}

func (i *whereIterable) Iterator() Iterator {
	return &whereIterator{i.items.Iterator(), i.where}
}
