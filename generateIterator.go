package c3

type generateIterator struct {
	g     Generate
	value interface{}
}

func (i *generateIterator) MoveNext() bool {
	value, ok := i.g()
	if ok {
		i.value = value
		return true
	}

	i.value = nil
	return false
}

func (i *generateIterator) Value() interface{} {
	return i.value
}
