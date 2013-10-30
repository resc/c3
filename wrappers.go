package c3

// WrapConsumer wraps a channel in a consuming Iterator
func WrapConsumer(c <-chan interface{}) Consumer {
	return &consumer{c, false, nil}
}

// Wraps a slice in a List interface
func WrapList(items []interface{}) List {
	return &list{0, items[:]}
}
