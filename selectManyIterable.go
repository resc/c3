package c3

type selectManyIterable struct {
	items    Iterable
	selector ManySelector
}

func (i *selectManyIterable) Iterator() Iterator {
	return &selectManyIterator{
		i.items.Iterator(),
		i.selector,
		emptyIterator,
		defaultElementValue,
	}
}
