package c3

import (
	"testing"
)

func TestPopFromEmptyStack(t *testing.T) {
	s := NewStack()

	item, ok := s.Pop()
	assert(t, false, ok, "ok")
	assert(t, nil, item, "item")
}

func TestPushStack(t *testing.T) {
	s := NewStack()
	pushOK := s.Push(333)
	assert(t, true, pushOK, "pushOK")

	item, ok := s.Pop()
	assert(t, true, ok, "ok")
	assert(t, 333, item, "item")
}

func TestConsumeEmptyStack(t *testing.T) {
	s := NewStack()

	for c := s.Consumer(); c.MoveNext(); {
		value := c.Value()
		failf(t, "Didn't expect value %v", value)
	}
}

func TestIterateEmptyStack(t *testing.T) {
	s := NewStack()

	for i := s.Iterator(); i.MoveNext(); {
		value := i.Value()
		failf(t, "Didn't expect value %v", value)
	}
}

func BenchmarkStackPush1000(b *testing.B) {
	value := wrap(1)
	s := NewStack()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for n := 0; n < 1000; n++ {
			s.Push(value)
		}
	}
}

func BenchmarkStackPushPop1000(b *testing.B) {
	value := wrap(1)
	s := NewStack()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for n := 0; n < 500; n++ {
			s.Push(value)
		}
		for n := 0; n < 500; n++ {
			s.Pop()
		}
	}
}

func BenchmarkConsumeStack1000(b *testing.B) {
	value := wrap(1)
	s := NewStack()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for n := 0; n < 500; n++ {
			s.Push(value)
		}
		for c := s.Consumer(); c.MoveNext(); {
			c.Value()
		}
	}
}
