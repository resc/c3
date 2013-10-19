package c3

type stack struct {
	*list
}

func NewStack() Stack {
	return &stack{newList()}
}

func StackOf(items ...interface{}) Stack {
	s := NewStack()
	for item := range items {
		s.Push(item)
	}
	return s
}

func (s *stack) Peek() (interface{}, bool) {
	return s.Last()
}

func (s *stack) Pop() (interface{}, bool) {
	if s.Len() > 0 {
		defer s.DeleteAt(s.Len() - 1)
	}
	return s.Last()
}

func (s *stack) Push(item interface{}) bool {
	return s.Add(item)
}

func (s *stack) Consumer() Consumer {
	return &stackConsumer{s, nil}
}

type stackConsumer struct {
	s     Stack
	value interface{}
}

func (sc *stackConsumer) MoveNext() bool {
	value, ok := sc.s.Pop()
	if ok {
		sc.value = value
		return true
	}
	sc.value = nil
	return false
}

func (sc *stackConsumer) Value() interface{} {
	return sc.value
}
