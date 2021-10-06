package spec

import (
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"sync"
)

type Spec struct {
	Name                string `yaml:"-"`
	Dependencies        []string
	Command             string
	Args                []string
	Environment         map[string]string
	ExpectedReturnCode  int    `yaml:"expectedReturnCode"`
	ExpectedStdOutRegex string `yaml:"expectedStdOutRegex"`
	ExpectedStdErrRegex string `yaml:"expectedStdErrRegex"`
	Order               int    `yaml:"-"`
	results             *ResultsProxy
}

func (s *Spec) GetStatus() Status {
	return s.results.GetStatus()
}

func (s *Spec) GetStdOut() string {
	return s.results.GetStdOut()
}

func (s *Spec) GetStdErr() string {
	return s.results.GetStdErr()
}

func (s *Spec) evaluateSuccess() {
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

func (s *Spec) Run(updateHandler func(*Spec)) error {
	var wg sync.WaitGroup
	s.results.SetStatus(StatusRunning)
	updateHandler(s)
	cmd := exec.Command(s.Command, s.Args...)
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
					fmt.Println(`Could not read std in from spec`, s.Name, `read bytes`, n, err)
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
		return fmt.Errorf(`command failed %q %v: %w`, s.Command, s.Args, err)
	}
	s.results.SetReturnCode(cmd.ProcessState.ExitCode())
	s.evaluateSuccess()
	updateHandler(s)
	return nil
}
