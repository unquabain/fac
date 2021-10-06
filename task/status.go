package task

// Status is an enum for the lifecycle of a Task.
type Status uint32

const (
	StatusNotRun Status = iota
	StatusDependenciesNotMet
	StatusRunning
	StatusFailed
	StatusSucceeded
)

func (s Status) String() string {
	switch s {
	case StatusNotRun:
		return `Waiting`
	case StatusDependenciesNotMet:
		return `Dependencies Not Met`
	case StatusRunning:
		return `Running`
	case StatusFailed:
		return `Failed`
	case StatusSucceeded:
		return `Succeeded`
	default:
		return `Unknown`
	}
}

// IsOK boils the status down to a simple boolean
// meaning whether or not the outcome is expected
// or exceptional.
func (s Status) IsOK() bool {
	switch s {
	case StatusNotRun:
		return true
	case StatusDependenciesNotMet:
		return false
	case StatusRunning:
		return true
	case StatusFailed:
		return false
	case StatusSucceeded:
		return true
	default:
		return false
	}
}
