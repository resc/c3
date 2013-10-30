package c3

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
