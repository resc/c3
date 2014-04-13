package c3

type queueIterator struct {
	q       *queue
	e       *entry
	done    bool
	version int
}

func (i *queueIterator) MoveNext() bool {
	if i.version != i.q.version {
		panic("Concurrent modification detected.")
	}

	if i.done {
		i.e = nil
		return false
	}

	if i.e == nil {
		i.e = i.q.head
	} else {
		i.e = i.e.next
		i.done = i.e.next == nil
	}
	return i.e != nil
}

func (i *queueIterator) Value() interface{} {
	if i.e == nil {
		return defaultElementValue
	}
	return i.e.item
}
