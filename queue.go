package c3

func NewQueue() Queue {
	return &queue{nil, nil, 0, 0}
}

func QueueOf(items ...interface{}) Queue {
	q := NewQueue()
	for item := range items {
		q.Enqueue(item)
	}
	return q
}

type queue struct {
	head    *entry
	tail    *entry
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
	return i.e.item
}

func (q *queue) Consumer() Consumer {
	return &qconsumer{q, nil}
}

type qconsumer struct {
	q    *queue
	item interface{}
}

func (c *qconsumer) MoveNext() bool {
	item, ok := c.q.Dequeue()
	c.item = item
	return ok
}

func (c *qconsumer) Value() interface{} {
	return c.item
}
