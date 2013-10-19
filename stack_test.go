package c3

import (
	"testing"
)

func TestPopFromEmptyStack(t *testing.T) {
	l := NewStack()
	item, ok := l.Pop()
	if ok || item != nil {
		t.Errorf("Expected nil,false got %v,%v", item, ok)
	}
}
