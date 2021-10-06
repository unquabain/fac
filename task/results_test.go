package task

import (
	"fmt"
	"testing"
)

func TestAtomic(t *testing.T) {
	r := NewResultsProxy()
	r.SetStdOut(`test`)
	r.Atomic(func(rr Results) {
		rr.SetStdErr(rr.GetStdOut())
		rr.SetStdOut(fmt.Sprintf(`%s OK`, rr.GetStdOut()))
	})

	if actual := r.GetStdOut(); actual != `test OK` {
		t.Fatalf(`expected stdout to be "test OK", was %q`, actual)
	}

	if actual := r.GetStdErr(); actual != `test` {
		t.Fatalf(`expected stderr to be "test", was %q`, actual)
	}
}

func TestAppendStdOut(t *testing.T) {
	r := NewResultsProxy()
	r.SetStdOut(`test`)
	r.AppendStdOut(` OK`)
	if actual := r.GetStdOut(); actual != `test OK` {
		t.Fatalf(`expected stdout to be "test OK", was %q`, actual)
	}
}

func TestAppendStdErr(t *testing.T) {
	r := NewResultsProxy()
	r.SetStdErr(`test`)
	r.AppendStdErr(` OK`)
	if actual := r.GetStdErr(); actual != `test OK` {
		t.Fatalf(`expected stderr to be "test OK", was %q`, actual)
	}
}

func TestSetSuccess(t *testing.T) {
	r := NewResultsProxy()
	r.SetStatus(StatusNotRun)
	r.SetSuccess()
	if actual := r.GetStatus(); actual != StatusSucceeded {
		t.Fatalf(`expected status to be success; was %v`, actual)
	}
	r.SetStatus(StatusFailed)
	r.SetSuccess()
	if actual := r.GetStatus(); actual == StatusSucceeded {
		t.Fatalf(`expected status not to be success; was %v`, actual)
	}
}
