package c3

type queueIterator struct {
	q       *queue
	e       *entry
	version int
}

func (i *queueIterator) MoveNext() bool {
	if i.version != i.q.version {
		panic("Concurrent modification detected.")
	}
	if i.e == nil {
		return false
	}
	i.e = i.e.next
	return i.e != nil
}

func (i *queueIterator) Value() interface{} {
	if i.e == nil {
		return nil
	}
	return i.e.item
}
