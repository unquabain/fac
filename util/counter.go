package util

import "sync"

type Counter struct {
	count int
	mtx   sync.RWMutex
}

func (c *Counter) Add(val int) int {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.count += val
	return c.count
}

func (c *Counter) Sub(val int) int {
	return c.Add(-val)
}

func (c *Counter) Inc() int {
	return c.Add(1)
}

func (c *Counter) Dec() int {
	return c.Add(-1)
}

func (c *Counter) Set(val int) int {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.count = val
	return c.count
}

func (c *Counter) Val() int {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	return c.count
}
