package c3

type stack struct {
	l *list
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

func (s *stack) Len() int {
	return s.l.Len()
}

func (s *stack) Peek() (interface{}, bool) {
	return s.l.Last()
}

func (s *stack) Pop() (interface{}, bool) {
	if s.l.Len() > 0 {
		defer s.l.DeleteAt(s.Len() - 1)
	}
	return s.l.Last()
}

func (s *stack) Contains(item interface{}) bool {
	return s.l.Contains(item)
}

func (s *stack) Push(item interface{}) bool {
	return s.l.Add(item)
}

func (s *stack) Iterator() Iterator {
	return &stackIterator{s.l, s.l.Len(), s.l.version, nil}
}

type stackIterator struct {
	l       *list
	index   int
	version int
	value   interface{}
}

func (si *stackIterator) MoveNext() bool {
	if si.version != si.l.version {
		si.value = nil
		panic("Concurrent modification detected")
	}

	si.index--

	if si.index >= 0 {
		si.value = si.l.items[si.index]
		return true
	}
	si.value = nil
	return false
}

func (sc *stackIterator) Value() interface{} {
	return sc.value
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
