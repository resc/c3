package c3

// WrapConsumable wraps a channel in a Consumable
func WrapConsumable(c <-chan interface{}) Consumable {
	return &consumable{c}
}

type consumable struct {
	c <-chan interface{}
}

func (c *consumable) Consumer() Consumer {
	return WrapConsumer(c.c)
}

// WrapConsumer wraps a channel in a consuming Iterator
func WrapConsumer(c <-chan interface{}) Consumer {
	return &consumer{c, false, nil}
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
	return nil
}
