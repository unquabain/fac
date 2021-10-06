package spec

import (
	"strings"
	"sync"
)

// Results represents the changing state of the Spec as the program is run.
type Results interface {
	// GetStdOut returns the accumulated text printed to stdout.
	GetStdOut() string

	// SetStdOut replaces the stored text of stdout.
	SetStdOut(string)

	// GetStdErr returns the accumulated text printed to stderr.
	GetStdErr() string

	// SetStdErr replaces the store text of stderr.
	SetStdErr(string)

	// GetReturnCode returns the return code of the executable after
	// it has been run.
	GetReturnCode() int

	// SetReturnCode replaces the stored return code.
	SetReturnCode(int)

	// GetStatus returns the current state of the Spec.
	GetStatus() Status

	// SetStatus replaces the current state of the Spec.
	SetStatus(Status)
}

type results struct {
	stdOut     string
	stdErr     string
	returnCode int
	status     Status
}

// GetStdOut returns the accumulated text printed to stdout.
// Implements Results interface.
func (r *results) GetStdOut() string {
	return r.stdOut
}

// SetStdOut replaces the stored text of stdout.
// Implements Results interface.
func (r *results) SetStdOut(stdOut string) {
	r.stdOut = stdOut
}

// GetStdErr returns the accumulated text printed to stderr.
// Implements Results interface.
func (r *results) GetStdErr() string {
	return r.stdErr
}

// SetStdErr replaces the store text of stderr.
// Implements Results interface.
func (r *results) SetStdErr(stdErr string) {
	r.stdErr = stdErr
}

// GetReturnCode returns the return code of the executable after
// it has been run.
// Implements Results interface.
func (r *results) GetReturnCode() int {
	return r.returnCode
}

// SetReturnCode replaces the stored return code.
// Implements Results interface.
func (r *results) SetReturnCode(returnCode int) {
	r.returnCode = returnCode
}

// GetStatus returns the current state of the Spec.
// Implements Results interface.
func (r *results) GetStatus() Status {
	return r.status
}

// SetStatus replaces the current state of the Spec.
// Implements Results interface.
func (r *results) SetStatus(status Status) {
	r.status = status
}

// ResultsProxy implements the Results interface, but allows only
// mutex-moderated access to the underlying data. Direct access
// can be achieved via the ResultsProxy.Atomic() method, which
// wraps the whole operation in a write lock.
type ResultsProxy struct {
	*results
	mtx sync.RWMutex
}

// NewResultsProxy generates a new ResultsProxy, intializing the
// internal data object.
func NewResultsProxy() *ResultsProxy {
	return &ResultsProxy{results: new(results)}
}

// Atomic allows multiple operations to be performed on
// the underlying data without worrying about other
// goroutines interleaving their operations.
func (r *ResultsProxy) Atomic(cb func(Results)) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	cb(r.results)
}

// GetStdOut returns the accumulated text printed to stdout.
// Implements Results interface.
func (r *ResultsProxy) GetStdOut() string {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	return r.results.stdOut
}

// SetStdOut replaces the stored text of stdout.
// Implements Results interface.
func (r *ResultsProxy) SetStdOut(stdOut string) {
	r.Atomic(func(results Results) { results.SetStdOut(stdOut) })
}

// AppendStdOut appends the string to the existing value
// in stdout atomically.
func (r *ResultsProxy) AppendStdOut(addendum string) {
	r.Atomic(func(results Results) {
		builder := new(strings.Builder)
		builder.WriteString(results.GetStdOut())
		builder.WriteString(addendum)
		results.SetStdOut(builder.String())
	})
}

// GetStdErr returns the accumulated text printed to stderr.
// Implements Results interface.
func (r *ResultsProxy) GetStdErr() string {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	return r.results.stdErr
}

// SetStdErr replaces the store text of stderr.
// Implements Results interface.
func (r *ResultsProxy) SetStdErr(stdErr string) {
	r.Atomic(func(results Results) { results.SetStdErr(stdErr) })
}

// AppendStdErr appends the string to the existing value
// in stderr atomically.
func (r *ResultsProxy) AppendStdErr(addendum string) {
	r.Atomic(func(results Results) {
		builder := new(strings.Builder)
		builder.WriteString(results.GetStdErr())
		builder.WriteString(addendum)
		results.SetStdErr(builder.String())
	})
}

// GetReturnCode returns the return code of the executable after
// it has been run.
// Implements Results interface.
func (r *ResultsProxy) GetReturnCode() int {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	return r.results.returnCode
}

// SetReturnCode replaces the stored return code.
// Implements Results interface.
func (r *ResultsProxy) SetReturnCode(returnCode int) {
	r.Atomic(func(results Results) { results.SetReturnCode(returnCode) })
}

// GetStatus returns the current state of the Spec.
// Implements Results interface.
func (r *ResultsProxy) GetStatus() Status {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	return r.results.status
}

// SetStatus replaces the current state of the Spec.
// Implements Results interface.
func (r *ResultsProxy) SetStatus(status Status) {
	r.Atomic(func(results Results) { results.SetStatus(status) })
}

// Sets the status to StatusSuccess, but only if it hasn't
// already been set to StatusFailed. Works atomically.
func (r *ResultsProxy) SetSuccess() {
	r.Atomic(func(results Results) {
		if results.GetStatus() != StatusFailed {
			results.SetStatus(StatusSucceeded)
		}
	})
}
