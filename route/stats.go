package route

import "sync/atomic"

type Counter struct {
	value int64
}

func NewCounter(initial int64) *Counter {
	return &Counter{initial}
}

func (c *Counter) Increment() {
	atomic.AddInt64(&c.value, 1)
}
func (c *Counter) Decrement() {
	atomic.AddInt64(&c.value, -1)
}
func (c *Counter) Count() int64 {
	return atomic.LoadInt64(&c.value)
}

type Stats struct {
	NumberConnections *Counter
}

func NewStats() *Stats {
	return &Stats{
		NumberConnections: &Counter{},
	}
}
