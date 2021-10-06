package spec

import (
	"fmt"
	"log"
	"strings"

	"github.com/Unquabain/thing-doer/util"
)

// SpecList represents all the Specs found in the spec file (YAML)
// and includes the methods to resolve their interdependencies and
// run them.
type SpecList map[string]*Spec

// UnmarshalYAML decorates the Specs found in the YAML spec file
// with some additional properties and initializes the Specs'
// internal structures.
func (sl SpecList) UnmarshalYAML(unmarshal func(interface{}) error) error {
	temp := make(map[string]*Spec)
	err := unmarshal(&temp)
	if err != nil {
		return err
	}
	count := 0
	for key, spec := range temp {
		spec.Name = key
		spec.Order = count
		count++
		spec.results = NewResultsProxy()
		sl[key] = spec
	}
	return nil
}

func parseDependencyName(depencencyName string) (key string, positive bool) {
	positive = true
	key = strings.TrimSpace(depencencyName)
	if strings.HasPrefix(key, `!`) || strings.HasPrefix(key, `-`) {
		positive = false
		key = strings.TrimSpace(key[1:])
	}
	return
}

// IsRunnable examines a Spec's dependency list and determines
// if it has been satisfied.
func (sl SpecList) IsRunnable(spec *Spec) (bool, error) {
	if spec.GetStatus() != StatusNotRun {
		return false, nil
	}
	if len(spec.Dependencies) == 0 {
		return true, nil
	}
	for _, dep := range spec.Dependencies {
		key, positive := parseDependencyName(dep)
		depSpec, ok := sl[key]
		if !ok {
			return false, fmt.Errorf(`dependency not found for %q: %q`, spec.Name, dep)
		}
		successStatus, failedStatus := StatusSucceeded, StatusFailed
		if !positive {
			successStatus, failedStatus = failedStatus, successStatus
		}
		dsStatus := depSpec.GetStatus()
		if dsStatus == failedStatus || dsStatus == StatusDependenciesNotMet {
			spec.results.SetStatus(StatusDependenciesNotMet)
			return false, nil
		}
		if dsStatus != successStatus {
			return false, nil
		}
	}
	return true, nil
}

// ReadyToRun returns a list of all the Specs that are currently
// ready to run because their dependencies have been satisified.
func (sl SpecList) ReadyToRun() ([]*Spec, error) {
	runnables := make([]*Spec, 0, len(sl))
	for _, spec := range sl {
		runnable, err := sl.IsRunnable(spec)
		if err != nil {
			return nil, fmt.Errorf(`invalid specification list: %w`, err)
		}
		if runnable {
			runnables = append(runnables, spec)
		}
	}
	return runnables, nil
}

// IsFinished tells the caller if all the Specs that can be
// run have been run (successfully or not).
func (sl SpecList) IsFinished() bool {
	for _, spec := range sl {
		if spec.GetStatus() == StatusNotRun {
			return false
		}
	}
	return true
}

// RunAll runs all the Specs, resolving their dependencies to
// run as many as it can in parallel. The function blocks until
// all Specs have been run, but the handler() callback will be
// called several times for each spec from different goroutines.
//
// First, all the specs that have no dependencies are run in
// parallel. After each spec finishes, the SpecList checks to
// see if any more specs have had their dependencies satisfied
// and launches those. The procedure runs until all specs have
// either run or been marked unrunnable (because their
// dependencies failed).
//
// If it ever finds that there are no currently running Specs,
// but no runnable Specs, but Specs that have not yet been run,
// it returns an error. It will also return an error if at least
// one Spec returns an error, though it may accumulate more errors,
// which are printed on STDERR.
func (sl SpecList) RunAll(handler func(*Spec)) error {
	runningTasks := new(util.Counter)
	errors := util.NewErrorList()
	gate := make(chan struct{})

	// Keep looping until all tasks report either finished,
	// skipped, or failed.
	for !sl.IsFinished() {
		rtr, err := sl.ReadyToRun() // All dependencies met successfully
		if err != nil {
			return fmt.Errorf(`failed determine runnable specs: %w`, err)
		}
		newTasks := len(rtr)
		if runningTasks.Val() == 0 && newTasks == 0 {
			// No running tasks, no new tasks, but some tasks are still
			// waiting to run. That means a dependency loop.
			return fmt.Errorf(`deadlock detected: not finished, but not ready to run`)
		}

		// Keep track of how many tasks are in-flight.
		runningTasks.Add(newTasks)

		// This loop may be empty if there are still
		// tasks running.
		for _, spec := range rtr {
			go func(s *Spec) {
				defer func() {
					runningTasks.Dec()
					// Gate is unbuffered. This will block until
					// the loop comes round again.
					gate <- struct{}{}
				}()
				errInner := s.Run(handler)
				if errInner != nil {
					errors.Appendf(`error running spec %q: %w`, s.Name, err)
					return
				}
			}(spec)
		}
		// Wait until one task finishes, then loop around
		// to re-evaluate if any tasks have had their
		// dependencies successfully run.
		// Because the above loop may not run, <-gate should
		// be called as many times as gate <- struct{}{} is
		// called.
		<-gate
		if errors.Len() > 0 {
			for _, err = range errors.Errors() {
				log.Println(err)
			}
			return fmt.Errorf(`received one or more errors running specs, the last of which is %w`, err)
		}
	}
	return nil
}
