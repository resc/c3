package c3

func NewQueue() Queue {
	return &queue{0, 0, nil, nil}
}

type queue struct {
	version int
	length  int
	head    *entry
	tail    *entry
}

type entry struct {
	value interface{}
	next  *entry
}

func (q *queue) Len() int {
	return q.length
}

func (q *queue) Contains(item interface{}) bool {
	for e := q.head; e != nil; e = e.next {
		if e.value == item {
			return true
		}
	}
	return false
}

func (q *queue) Peek() (interface{}, bool) {
	if q.head != nil {
		return q.head.value, true
	}
	return nil, false
}

func (q *queue) Dequeue() (interface{}, bool) {
	if q.head != nil {
		e := q.head
		q.head = e.next
		if e.next == nil {
			q.tail = nil
		}
		q.length--
		q.version++
		return e.value, true
	}
	return nil, false
}

func (q *queue) Enqueue(item interface{}) bool {
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
	return &qiterator{q, q.head, q.version}
}

type qiterator struct {
	q       *queue
	e       *entry
	version int
}

func (i *qiterator) MoveNext() bool {
	if i.version != i.q.version {
		panic("Concurrent modification detected.")
	}
	if i.e == nil {
		return false
	}
	i.e = i.e.next
	return i.e != nil
}

func (i *qiterator) Value() interface{} {
	if i.e == nil {
		return nil
	}
	return i.e.value
}

func (q *queue) Consumer() Consumer {
	return &qconsumer{q, nil}
}

type qconsumer struct {
	q     *queue
	value interface{}
}

func (qc *qconsumer) MoveNext() bool {
	value, ok := qc.q.Dequeue()
	qc.value = value
	return ok
}

func (qc *qconsumer) Value() interface{} {
	return qc.value
}
