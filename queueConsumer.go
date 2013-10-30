package c3

type queueConsumer struct {
	q    *queue
	item interface{}
}

func (c *queueConsumer) MoveNext() bool {
	item, ok := c.q.Dequeue()
	c.item = item
	return ok
}

func (c *queueConsumer) Value() interface{} {
	return c.item
}
