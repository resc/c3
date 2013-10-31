package c3

import "testing"

//import "runtime"

func TestEmptyQueue(t *testing.T) {
	q := NewQueue()
	if q.Len() != 0 {
		t.Errorf("Expected 0, got %v", q.Len())
	}
}
func TestQueueIterator(t *testing.T) {
	q := QueueOf(0, 1, 2, 3)
	index := 0
	for i := q.Iterator(); i.MoveNext(); {
		val, ok := i.Value().(int)
		if !ok || val != index {
			t.Error("Expected ", index, " got ", val)
		}
		index++
	}
}

func TestQueue(t *testing.T) {
	q := NewQueue()
	assert(t, 0, q.Len(), "q.Len()")

	q.Enqueue(9999)
	assert(t, 1, q.Len(), "q.Len()")
	assert(t, true, q.Contains(9999), "q.Contains(9999)")

	result, ok := q.Peek()
	assert(t, 9999, result, "result")
	assert(t, true, ok, "ok")
	assert(t, 1, q.Len(), "q.Len()")

	result, ok = q.Dequeue()
	assert(t, 9999, result, "result")
	assert(t, true, ok, "ok")
	assert(t, 0, q.Len(), "q.Len()")

	result, ok = q.Dequeue()
	assert(t, false, ok, "ok")
	assert(t, 0, q.Len(), "q.Len()")
}

func BenchmarkEnqueue1000(b *testing.B) {
	value := wrap(1)
	q := NewQueue()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for n := 0; n < 1000; n++ {
			q.Enqueue(value)
		}
	}
}

func BenchmarkEnqDeq1000(b *testing.B) {
	value := wrap(1)
	q := NewQueue()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for n := 0; n < 1000; n++ {
			q.Enqueue(value)
		}
		for n := 0; n < 1000; n++ {
			q.Dequeue()
		}
	}
}

func BenchmarkEnqConsume1000(b *testing.B) {
	value := wrap(1)
	q := NewQueue()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for n := 0; n < 1000; n++ {
			q.Enqueue(value)
		}
		for i := q.Consumer(); i.MoveNext(); {
			i.Value()
		}
	}
}

func wrap(item interface{}) interface{} {
	return item
}
