package util

import "sync"

// Counter is a thread-safe wrapper around an integer.
// It's like sync.WaitGroup, but it doesn't block anything,
// just keeps track of the number.
type Counter struct {
	count int
	mtx   sync.RWMutex
}

// Add adds a value to the existing value atomically.
func (c *Counter) Add(val int) int {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.count += val
	return c.count
}

// Sub subtracts a value from the existing value atomically.
func (c *Counter) Sub(val int) int {
	return c.Add(-val)
}

// Inc adds one to the existing value atomically.
func (c *Counter) Inc() int {
	return c.Add(1)
}

// Dec subtracts one from the existing value atomically.
func (c *Counter) Dec() int {
	return c.Add(-1)
}

// Set replaces the current value atomically.
func (c *Counter) Set(val int) int {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.count = val
	return c.count
}

// Val returns the current value atomically.
func (c *Counter) Val() int {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	return c.count
}
