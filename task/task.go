package task

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"sync"
)

// Task is the taskification of a task to run
// and its dependencies.
type Task struct {
	// Name is the name of the task, and how other tasks
	// will refer to it in their dependency list.
	Name string `yaml:"-"`

	// Dependencies is a list of other Task.Name values
	// that must succeed (or fail if they are prefixed with
	// either "!" or "-") before this task can be run.
	Dependencies []string

	// Command is the executable to run.
	Command string

	// Args are arguments to pass to Command
	Args []string

	// Environment is any shell environment variables
	// that Command will need.
	Environment map[string]string

	// ExpectedReturnCode is the return code that
	// Command should result in to consider this Task
	// successful.
	ExpectedReturnCode int `yaml:"expectedReturnCode"`

	// ExpectedStdOutRegex is a pattern that, if present,
	// will be checked against the full STDOUT to qualify
	// this Task run as successful.
	ExpectedStdOutRegex string `yaml:"expectedStdOutRegex"`

	// ExpectedStdErrRegex is a pattern that, if present,
	// will be checked against the full STDERR to qualify
	// this Task run as successful.
	ExpectedStdErrRegex string `yaml:"expectedStdErrRegex"`

	// Order is set in the YAML parser for consistency
	// in the interface. (Otherwise, the list reshuffles
	// whenever it updates.)
	Order int `yaml:"-"`

	results *ResultsProxy
}

// GetStatus gets the current status atomically.
func (s *Task) GetStatus() Status {
	return s.results.GetStatus()
}

// GetStdOut gets the accumulated STDOUT text atomically.
func (s *Task) GetStdOut() string {
	return s.results.GetStdOut()
}

// GetStdErr gets the accumulated STDERR text atomically.
func (s *Task) GetStdErr() string {
	return s.results.GetStdErr()
}

func (s *Task) evaluateSuccess() {
	if s.results.GetStatus() == StatusFailed {
		return
	}
	if s.results.GetReturnCode() != s.ExpectedReturnCode {
		s.results.SetStatus(StatusFailed)
		return
	}
	if s.ExpectedStdOutRegex != `` {
		pattern := regexp.MustCompile(s.ExpectedStdOutRegex)
		if !pattern.MatchString(s.results.GetStdOut()) {
			s.results.SetStatus(StatusFailed)
			return
		}
	}
	if s.ExpectedStdErrRegex != `` {
		pattern := regexp.MustCompile(s.ExpectedStdErrRegex)
		if !pattern.MatchString(s.results.GetStdErr()) {
			s.results.SetStatus(StatusFailed)
			return
		}
	}
	s.results.SetSuccess()
}

func (s *Task) env() []string {
	env := make([]string, 0, len(os.Environ())+len(s.Environment))
	copy(env, os.Environ())
	for key, val := range s.Environment {
		env = append(env, fmt.Sprintf(`%s=%s`, key, val))
	}
	return env
}

// Run runs the command defined by Task. It blocks until the
// command has finished, but updateHandler will be called
// several times from different go routines whenever a change
// has been made to the status of the Task.
func (s *Task) Run(updateHandler func(*Task)) error {
	var wg sync.WaitGroup
	s.results.SetStatus(StatusRunning)
	updateHandler(s)
	cmd := exec.Command(s.Command, s.Args...)
	cmd.Env = s.env()
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf(`couldn't open standard out for command %q %v: %w`, s.Command, s.Args, err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf(`couldn't open standard error for command %q %v: %w`, s.Command, s.Args, err)
	}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf(`couldn't start command %q %v: %w`, s.Command, s.Args, err)
	}
	wg.Add(2)
	go func() {
		defer stdout.Close()
		defer wg.Done()
		buff := make([]byte, 1024)
		for {
			n, err := stdout.Read(buff)
			if n > 0 {
				s.results.AppendStdOut(string(buff))
				updateHandler(s)
			}
			if err != nil {
				if err != io.EOF {
					fmt.Println(`Could not read std in from task`, s.Name, `read bytes`, n, err)
					s.results.SetStatus(StatusFailed)
				}
				return
			}
		}
	}()

	go func() {
		defer stderr.Close()
		defer wg.Done()
		buff := make([]byte, 1024)
		for {
			n, err := stderr.Read(buff)
			if n > 0 {
				s.results.AppendStdErr(string(buff))
				updateHandler(s)
			}
			if err != nil {
				if err != io.EOF {
					fmt.Println(err)
					s.results.SetStatus(StatusFailed)
				}
				return
			}
		}
	}()
	wg.Wait()
	if err := cmd.Wait(); err != nil {
		s.results.SetStatus(StatusFailed)
		s.results.AppendStdErr(fmt.Sprintf(`command failed %q %v: %v`, s.Command, s.Args, err))
	}
	s.results.SetReturnCode(cmd.ProcessState.ExitCode())
	s.evaluateSuccess()
	updateHandler(s)
	return nil
}
