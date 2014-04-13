package c3

type listIterator struct {
	l       *list
	version int
	index   int
	value   interface{}
}

func (i *listIterator) MoveNext() bool {
	if i.version != i.l.version {
		i.value = defaultElementValue
		panic("Concurrent modification detected")
	}

	if i.index < len(i.l.items)-1 {
		i.index++
		i.value = i.l.items[i.index]
		return true
	}

	i.value = defaultElementValue
	return false
}

func (i *listIterator) Value() interface{} {
	return i.value
}
