package spec

import (
	"strings"
	"sync"
)

type Results interface {
	GetStdOut() string
	SetStdOut(string)
	GetStdErr() string
	SetStdErr(string)
	GetReturnCode() int
	SetReturnCode(int)
	GetStatus() Status
	SetStatus(Status)
}

type results struct {
	stdOut     string
	stdErr     string
	returnCode int
	status     Status
}

func (r *results) GetStdOut() string {
	return r.stdOut
}

func (r *results) SetStdOut(stdOut string) {
	r.stdOut = stdOut
}

func (r *results) GetStdErr() string {
	return r.stdErr
}

func (r *results) SetStdErr(stdErr string) {
	r.stdErr = stdErr
}

func (r *results) GetReturnCode() int {
	return r.returnCode
}

func (r *results) SetReturnCode(returnCode int) {
	r.returnCode = returnCode
}

func (r *results) GetStatus() Status {
	return r.status
}

func (r *results) SetStatus(status Status) {
	r.status = status
}

type ResultsProxy struct {
	*results
	mtx sync.RWMutex
}

func NewResultsProxy() *ResultsProxy {
	return &ResultsProxy{results: new(results)}
}

func (r *ResultsProxy) Atomic(cb func(Results)) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	cb(r.results)
}

func (r *ResultsProxy) GetStdOut() string {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	return r.results.stdOut
}

func (r *ResultsProxy) SetStdOut(stdOut string) {
	r.Atomic(func(results Results) { results.SetStdOut(stdOut) })
}

func (r *ResultsProxy) AppendStdOut(addendum string) {
	r.Atomic(func(results Results) {
		builder := new(strings.Builder)
		builder.WriteString(results.GetStdOut())
		builder.WriteString(addendum)
		results.SetStdOut(builder.String())
	})
}

func (r *ResultsProxy) GetStdErr() string {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	return r.results.stdErr
}

func (r *ResultsProxy) SetStdErr(stdErr string) {
	r.Atomic(func(results Results) { results.SetStdErr(stdErr) })
}

func (r *ResultsProxy) AppendStdErr(addendum string) {
	r.Atomic(func(results Results) {
		builder := new(strings.Builder)
		builder.WriteString(results.GetStdErr())
		builder.WriteString(addendum)
		results.SetStdErr(builder.String())
	})
}

func (r *ResultsProxy) GetReturnCode() int {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	return r.results.returnCode
}

func (r *ResultsProxy) SetReturnCode(returnCode int) {
	r.Atomic(func(results Results) { results.SetReturnCode(returnCode) })
}

func (r *ResultsProxy) GetStatus() Status {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	return r.results.status
}

func (r *ResultsProxy) SetStatus(status Status) {
	r.Atomic(func(results Results) { results.SetStatus(status) })
}

func (r *ResultsProxy) SetSuccess() {
	r.Atomic(func(results Results) {
		if results.GetStatus() != StatusFailed {
			results.SetStatus(StatusSucceeded)
		}
	})
}
