package c3

import (
	"testing"
)

func TestPopFromEmptyStack(t *testing.T) {
	s := NewStack()

	item, ok := s.Pop()
	if ok || item != nil {
		t.Errorf("Expected nil,false got %v,%v", item, ok)
	}
}

func TestConsumeEmptyStack(t *testing.T) {
	s := NewStack()

	for c := s.Consumer(); c.MoveNext(); {
		value := c.Value()
		t.Error("Didn't expect value %v", value)
	}
}
