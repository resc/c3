package c3

import "testing"

func TestConsumerOfChannel(t *testing.T) {
	ch := make(chan interface{}, 10)
	NewQuery(Range(0, 9)).For(func(e interface{}) { ch <- e })
	close(ch)
	count := 0
	for c := WrapConsumer(ch); c.MoveNext(); {
		count++
	}
	assert(t, 10, count, "iterations")
}
