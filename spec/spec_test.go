package spec

import (
	"testing"
)

func newSuccessfulSpec() *Spec {
	return &Spec{
		Name:                `Test Spec`,
		Dependencies:        []string{`Dependency `, `Depencency 2`},
		Command:             `cat`,
		Args:                []string{`./test_data/success_data.txt`},
		Environment:         make(map[string]string),
		ExpectedReturnCode:  0,
		ExpectedStdOutRegex: `Successful Value`,
		ExpectedStdErrRegex: ``,
		results:             NewResultsProxy(),
	}
}

func newFailValidationSpec() *Spec {
	s := newSuccessfulSpec()
	s.Args[0] = `./test_data/fail_data.txt`
	return s
}

func TestRunWithSuccess(t *testing.T) {
	spec := newSuccessfulSpec()
	updatesCount := 0
	err := spec.Run(func(s *Spec) { updatesCount++ })
	if err != nil {
		t.Fatalf(`problem running command: %v`, err)
	}
	if actual := spec.GetStatus(); actual != StatusSucceeded {
		t.Fatalf(`spec should have succeeded; didn't: %v`, actual)
	}
	if actual := updatesCount; actual < 3 {
		t.Fatalf(`expected at least 3 updates: received %d`, actual)
	}
}

func TestRunWithFailure(t *testing.T) {
	spec := newFailValidationSpec()
	updatesCount := 0
	err := spec.Run(func(s *Spec) { updatesCount++ })
	if err != nil {
		t.Fatalf(`problem running command: %v`, err)
	}
	if actual := spec.GetStatus(); actual != StatusFailed {
		t.Fatalf(`spec should not have succeeded; did: %v`, actual)
	}
	if actual := updatesCount; actual < 3 {
		t.Fatalf(`expected at least 3 updates: received %d`, actual)
	}
}
