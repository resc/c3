package c3

type stackIterator struct {
	l       *list
	index   int
	version int
	value   interface{}
}

func (si *stackIterator) MoveNext() bool {
	if si.version != si.l.version {
		si.value = defaultElementValue
		panic("Concurrent modification detected")
	}

	si.index--

	if si.index >= 0 {
		si.value = si.l.items[si.index]
		return true
	}
	si.value = defaultElementValue
	return false
}

func (sc *stackIterator) Value() interface{} {
	return sc.value
}
