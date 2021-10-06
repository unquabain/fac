package util

import (
	"fmt"
	"sync"
)

// ErrorList is a list of errors protected by a mutex for
// use in accumulating errors from multiple, concurrent
// goroutines.
type ErrorList struct {
	errors []error
	mtx    sync.RWMutex
}

// NewErrorList creates a new ErrorList, initializing
// its internal structures.
func NewErrorList() *ErrorList {
	el := new(ErrorList)
	el.errors = make([]error, 0)
	return el
}

// Append appends a new error to the list atomically.
func (el *ErrorList) Append(err error) []error {
	el.mtx.Lock()
	defer el.mtx.Unlock()
	el.errors = append(el.errors, err)
	return el.errors
}

// Appendf is like fmt.Errorf, but appends the generated error
// to the list atomically.
func (el *ErrorList) Appendf(format string, args ...interface{}) []error {
	return el.Append(fmt.Errorf(format, args...))
}

// Errors returns the accumulated errors atomically.
func (el *ErrorList) Errors() []error {
	el.mtx.RLock()
	defer el.mtx.RUnlock()
	return el.errors
}

// Len is how many errors have been accumulated. Atomic.
func (el *ErrorList) Len() int {
	el.mtx.RLock()
	defer el.mtx.RUnlock()
	return len(el.errors)
}
