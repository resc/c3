package c3

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
	i.value = defaultElementValue
	return false
}

func (i *selectIterator) Value() interface{} {
	return i.value
}
