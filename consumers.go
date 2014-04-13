package c3

type consumable struct {
	c <-chan interface{}
}

func (c *consumable) Consumer() Consumer {
	return WrapConsumer(c.c)
}

type consumer struct {
	c        <-chan interface{}
	hasvalue bool
	value    interface{}
}

func (i *consumer) MoveNext() bool {
	value, hasvalue := <-i.c
	i.value = value
	i.hasvalue = hasvalue
	return hasvalue
}

func (i *consumer) Value() interface{} {
	if i.hasvalue {
		return i.value
	}
	return defaultElementValue
}
