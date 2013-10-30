package c3

type selectIterable struct {
	items    Iterable
	selector Selector
}

func (i *selectIterable) Iterator() Iterator {
	return &selectIterator{i.items.Iterator(), i.selector, nil}
}
