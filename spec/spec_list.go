package spec

import (
	"fmt"
	"log"

	"github.com/Unquabain/thing-doer/util"
)

type SpecList map[string]*Spec

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

func (sl SpecList) IsRunnable(spec *Spec) (bool, error) {
	if spec.GetStatus() != StatusNotRun {
		return false, nil
	}
	if len(spec.Dependencies) == 0 {
		return true, nil
	}
	for _, dep := range spec.Dependencies {
		depSpec, ok := sl[dep]
		if !ok {
			return false, fmt.Errorf(`dependency not found for %q: %q`, spec.Name, dep)
		}
		dsStatus := depSpec.GetStatus()
		if dsStatus == StatusFailed || dsStatus == StatusDependenciesNotMet {
			spec.results.SetStatus(StatusDependenciesNotMet)
			return false, nil
		}
		if dsStatus != StatusSucceeded {
			return false, nil
		}
	}
	return true, nil
}

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

func (sl SpecList) IsFinished() bool {
	for _, spec := range sl {
		if spec.GetStatus() == StatusNotRun {
			return false
		}
	}
	return true
}

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
