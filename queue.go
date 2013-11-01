package c3

type queue struct {
	head    *entry
	tail    *entry
	free *entry
	version int
	length  int
}

type entry struct {
	item interface{}
	next *entry
}

func (q *queue) Len() int {
	return q.length
}

func (q *queue) Clear() {
	if q.length == 0 {
		return
	}
	q.head = nil
	q.tail = nil
	q.free = nil
	q.length = 0
	q.version++
}

func (q *queue) Contains(item interface{}) bool {
	for e := q.head; e != nil; e = e.next {
		if e.item == item {
			return true
		}
	}
	return false
}

func (q *queue) Peek() (interface{}, bool) {
	if q.head != nil {
		return q.head.item, true
	}
	return nil, false
}

func (q *queue) Dequeue() (interface{}, bool) {
// TODO add removed entry to the free list.
	if q.head != nil {
		e := q.head
		q.head = e.next
		if e.next == nil {
			q.tail = nil
		}
		q.length--
		q.version++
		return e.item, true
	}
	return nil, false
}

func (q *queue) Enqueue(item interface{}) bool {
// TODO check if there's any free entries.
	e := &entry{item, nil}
	if q.tail == nil {
		q.head = e
		q.tail = e
	} else {
		q.tail.next = e
		q.tail = e
	}
	q.version++
	q.length++
	return true
}

func (q *queue) Iterator() Iterator {
	return &queueIterator{q, nil, false, q.version}
}

func (q *queue) Consumer() Consumer {
	return &queueConsumer{q, nil}
}
