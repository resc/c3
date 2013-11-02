package c3

const (
	// queue_max_free_length is an arbitratry large-ish number of free entries
	// to keep around to prevent busywork for the garbage collector.
	// You can see the effect in the performance difference between
	// BenchmarkEnqueue1000 and BenchmarkEnqDeq1000.
	queue_max_free_length int = 1024
)

type queue struct {
	head *entry
	tail *entry
	free *entry

	version    int
	length     int
	freeLength int
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
	q.freeLength = 0
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
	if q.head != nil {
		// remove entry from queue
		e := q.head
		item := e.item
		e.item = nil
		q.head = e.next
		if e.next == nil {
			q.tail = nil
		}

		// add freed entry to free list but don't keep
		// too many of them, it's a waste of space
		if q.freeLength < queue_max_free_length {
			e.next = q.free
			q.free = e
			q.freeLength++
		}

		q.length--
		q.version++
		return item, true
	}
	return nil, false
}

func (q *queue) Enqueue(item interface{}) bool {
	e := q.free
	if e != nil {
		// get entry from the free list
		q.free = e.next
		q.freeLength--
		e.item = item
		e.next = nil
	} else {
		// create a new entry
		e = &entry{item, nil}
	}

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
