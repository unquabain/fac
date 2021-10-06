package util

import (
	"fmt"
	"sync"
)

type ErrorList struct {
	errors []error
	mtx    sync.RWMutex
}

func NewErrorList() *ErrorList {
	el := new(ErrorList)
	el.errors = make([]error, 0)
	return el
}

func (el *ErrorList) Append(err error) []error {
	el.mtx.Lock()
	defer el.mtx.Unlock()
	el.errors = append(el.errors, err)
	return el.errors
}

func (el *ErrorList) Appendf(format string, args ...interface{}) []error {
	return el.Append(fmt.Errorf(format, args...))
}

func (el *ErrorList) Errors() []error {
	el.mtx.RLock()
	defer el.mtx.RUnlock()
	return el.errors
}

func (el *ErrorList) Len() int {
	el.mtx.RLock()
	defer el.mtx.RUnlock()
	return len(el.errors)
}
