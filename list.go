package c3

type list struct {
	version int
	items   []interface{}
}

func (l *list) Iterator() Iterator {
	return &listIterator{l, l.version, -1, defaultElementValue}
}

func (l *list) Add(item interface{}) bool {
	l.items = append(l.items, item)
	l.version++
	return true
}

func (l *list) Swap(i, j int) {
	ival := l.items[i]
	jval := l.items[j]
	l.items[i] = jval
	l.items[j] = ival
	l.version++
}

func (l *list) Clear() {
	if l.Len() == 0 {
		return
	}
	if l.Len() <= 1024 {
		for i := 0; i < len(l.items); i++ {
			l.items[i] = defaultElementValue
		}
		l.items = l.items[:0]
	} else {
		l.items = make([]interface{}, 0, 4)
	}
	l.version++
}

func (l *list) InsertAt(index int, item interface{}) bool {
	if 0 > index || index > len(l.items) {
		return false
	}

	if index == len(l.items) {
		return l.Add(item)
	}

	l.items = append(l.items, defaultElementValue)
	copy(l.items[index+1:], l.items[index:])
	l.items[index] = item
	l.version++
	return true
}

func (l *list) First() (interface{}, bool) {
	return l.Get(0)
}

func (l *list) Last() (interface{}, bool) {
	return l.Get(l.Len() - 1)
}

func (l *list) Get(index int) (interface{}, bool) {
	if 0 > index || index >= len(l.items) {
		return defaultElementValue, false
	}
	return l.items[index], true
}

func (l *list) Contains(item interface{}) bool {
	_, ok := l.IndexOf(item)
	return ok
}

func (l *list) IndexOf(item interface{}) (int, bool) {
	return l.NextIndexOf(-1, item)
}

func (l *list) NextIndexOf(offset int, item interface{}) (int, bool) {
	for index := max(-1, offset) + 1; 0 <= index && index < l.Len(); index++ {
		if l.items[index] == item {
			return index, true
		}
	}
	return -1, false
}

func (l *list) LastIndexOf(item interface{}) (int, bool) {
	return l.PrevIndexOf(l.Len(), item)
}

func (l *list) PrevIndexOf(offset int, item interface{}) (int, bool) {
	for index := min(offset, l.Len()) - 1; 0 <= index && index < l.Len(); index-- {
		if l.items[index] == item {
			return index, true
		}
	}
	return -1, false
}

func (l *list) Delete(item interface{}) bool {
	if index, ok := l.IndexOf(item); ok {
		return l.DeleteAt(index)
	}
	return false
}

func (l *list) DeleteAt(index int) bool {
	if 0 > index || index >= len(l.items) {
		return false
	}
	last := len(l.items) - 1
	if index != last {
		copy(l.items[index:], l.items[index+1:])
	}
	l.items[last] = defaultElementValue
	l.items = l.items[:last]
	l.version++
	return true
}

func (l *list) Len() int {
	return len(l.items)
}
