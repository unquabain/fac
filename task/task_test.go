package task

import (
	"testing"
)

func newSuccessfulTask() *Task {
	return &Task{
		Name:                `Test Task`,
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

func newFailValidationTask() *Task {
	s := newSuccessfulTask()
	s.Args[0] = `./test_data/fail_data.txt`
	return s
}

func TestRunWithSuccess(t *testing.T) {
	task := newSuccessfulTask()
	updatesCount := 0
	err := task.Run(func(s *Task) { updatesCount++ })
	if err != nil {
		t.Fatalf(`problem running command: %v`, err)
	}
	if actual := task.GetStatus(); actual != StatusSucceeded {
		t.Fatalf(`task should have succeeded; didn't: %v`, actual)
	}
	if actual := updatesCount; actual < 3 {
		t.Fatalf(`expected at least 3 updates: received %d`, actual)
	}
}

func TestRunWithFailure(t *testing.T) {
	task := newFailValidationTask()
	updatesCount := 0
	err := task.Run(func(s *Task) { updatesCount++ })
	if err != nil {
		t.Fatalf(`problem running command: %v`, err)
	}
	if actual := task.GetStatus(); actual != StatusFailed {
		t.Fatalf(`task should not have succeeded; did: %v`, actual)
	}
	if actual := updatesCount; actual < 3 {
		t.Fatalf(`expected at least 3 updates: received %d`, actual)
	}
}
