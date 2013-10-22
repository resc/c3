package c3

import "testing"

//import "runtime"

func TestEmptyQueue(t *testing.T) {
	q := NewQueue()
	if q.Len() != 0 {
		t.Errorf("Expected 0, got %v", q.Len())
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

func BenchmarkEnqueue(b *testing.B) {
	q := NewQueue()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		q.Enqueue(i)
	}
}

func BenchmarkDequeue(b *testing.B) {
	q := NewQueue()

	// Too fast to benchmark, memory usage wil
	// blow this up beyond 10.000.000 iterations.
	if b.N > 10000000 {
		b.Logf("%v iterations will kill the benchmark due to mem usage", b.N)
		return
	}

	for i := 0; i < b.N; i++ {
		q.Enqueue(i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		q.Dequeue()
	}
}

func BenchmarkConsumeQueue(b *testing.B) {
	q := NewQueue()

	// Too fast to benchmark, memory usage wil
	// blow this up beyond 10.000.000 iterations.
	if b.N > 10000000 {
		b.Logf("%v iterations will kill the benchmark due to mem usage", b.N)
		return
	}

	for i := 0; i < b.N; i++ {
		q.Enqueue(i)
	}

	b.ResetTimer()

	for i := q.Consumer(); i.MoveNext(); {
		i.Value()
	}
}
