package c3

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
	i.iterator = emptyIterator
	i.value = defaultElementValue
	return false
}

func (i *selectManyIterator) Value() interface{} {
	return i.value
}
