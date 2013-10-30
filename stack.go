package c3

type stack struct {
	l *list
}

func (s *stack) Len() int {
	return s.l.Len()
}

func (s *stack) Clear() {
	s.l.Clear()
}

func (s *stack) Peek() (interface{}, bool) {
	return s.l.Last()
}

func (s *stack) Pop() (interface{}, bool) {
	item, ok := s.l.Last()
	if ok {
		s.l.DeleteAt(s.Len() - 1)
	}
	return item, ok
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

func (s *stack) Consumer() Consumer {
	return &stackConsumer{s, nil}
}
