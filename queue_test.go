package c3

import "testing"

func TestNewQueue(t *testing.T) {
	q := NewQueue()

	if q.Len() != 0 {
		t.Errorf("%v ", q.Len())
	}

	q.Enqueue(9999)

	if q.Len() != 1 {
		t.Errorf("%v ", q.Len())
	}

	if !q.Contains(9999) {
		t.Errorf("%v ", q.Contains(9999))
	}

	result, ok := q.Peek()
	if !ok || result != 9999 {
		t.Errorf("%v %v", result, ok)
	}

	if q.Len() != 1 {
		t.Errorf("%v ", q.Len())
	}

	result, ok = q.Dequeue()
	if !ok || result != 9999 {
		t.Errorf("%v %v", result, ok)
	}

	result, ok = q.Dequeue()
	if ok || result != nil {
		t.Errorf("%v %v", result, ok)
	}

	if q.Len() != 0 {
		t.Errorf("%v ", q.Len())
	}

}
