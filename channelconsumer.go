package c3

// ChannelConsumer wraps a channel in a consuming Iterator
func ChannelConsumer(c <-chan interface{}) Consumer {
	return &chaniter{c, false, nil}
}

type chaniter struct {
	c        <-chan interface{}
	hasvalue bool
	value    interface{}
}

func (i *chaniter) MoveNext() bool {
	value, hasvalue := <-i.c
	i.value = value
	i.hasvalue = hasvalue
	return hasvalue
}

func (i *chaniter) Value() interface{} {
	if i.hasvalue {
		return i.value
	}
	return nil
}
