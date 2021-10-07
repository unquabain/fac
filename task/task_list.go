package task

import (
	"fmt"
	"log"
	"strings"

	"github.com/Unquabain/fac/util"
)

// TaskList represents all the Tasks found in the task file (YAML)
// and includes the methods to resolve their interdependencies and
// run them.
type TaskList map[string]*Task

// UnmarshalYAML decorates the Tasks found in the YAML task file
// with some additional properties and initializes the Tasks'
// internal structures.
func (sl TaskList) UnmarshalYAML(unmarshal func(interface{}) error) error {
	temp := make(map[string]*Task)
	err := unmarshal(&temp)
	if err != nil {
		return err
	}
	count := 0
	for key, task := range temp {
		task.Name = key
		task.Order = count
		count++
		task.results = NewResultsProxy()
		sl[key] = task
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

// IsRunnable examines a Task's dependency list and determines
// if it has been satisfied.
func (sl TaskList) IsRunnable(task *Task) (bool, error) {
	if task.GetStatus() != StatusNotRun {
		return false, nil
	}
	if len(task.Dependencies) == 0 {
		return true, nil
	}
	for _, dep := range task.Dependencies {
		key, positive := parseDependencyName(dep)
		depTask, ok := sl[key]
		if !ok {
			return false, fmt.Errorf(`dependency not found for %q: %q`, task.Name, dep)
		}
		successStatus, failedStatus := StatusSucceeded, StatusFailed
		if !positive {
			successStatus, failedStatus = failedStatus, successStatus
		}
		dsStatus := depTask.GetStatus()
		if dsStatus == failedStatus || dsStatus == StatusDependenciesNotMet {
			task.results.SetStatus(StatusDependenciesNotMet)
			return false, nil
		}
		if dsStatus != successStatus {
			return false, nil
		}
	}
	return true, nil
}

// ReadyToRun returns a list of all the Tasks that are currently
// ready to run because their dependencies have been satisified.
func (sl TaskList) ReadyToRun() ([]*Task, error) {
	runnables := make([]*Task, 0, len(sl))
	for _, task := range sl {
		runnable, err := sl.IsRunnable(task)
		if err != nil {
			return nil, fmt.Errorf(`invalid taskification list: %w`, err)
		}
		if runnable {
			runnables = append(runnables, task)
		}
	}
	return runnables, nil
}

// IsFinished tells the caller if all the Tasks that can be
// run have been run (successfully or not).
func (sl TaskList) IsFinished() bool {
	for _, task := range sl {
		if task.GetStatus() == StatusNotRun {
			return false
		}
	}
	return true
}

// RunAll runs all the Tasks, resolving their dependencies to
// run as many as it can in parallel. The function blocks until
// all Tasks have been run, but the handler() callback will be
// called several times for each task from different goroutines.
//
// First, all the tasks that have no dependencies are run in
// parallel. After each task finishes, the TaskList checks to
// see if any more tasks have had their dependencies satisfied
// and launches those. The procedure runs until all tasks have
// either run or been marked unrunnable (because their
// dependencies failed).
//
// If it ever finds that there are no currently running Tasks,
// but no runnable Tasks, but Tasks that have not yet been run,
// it returns an error. It will also return an error if at least
// one Task returns an error, though it may accumulate more errors,
// which are printed on STDERR.
func (sl TaskList) RunAll(handler func(*Task)) error {
	runningTasks := new(util.Counter)
	errors := util.NewErrorList()
	gate := make(chan struct{})

	// Keep looping until all tasks report either finished,
	// skipped, or failed.
	for !sl.IsFinished() {
		rtr, err := sl.ReadyToRun() // All dependencies met successfully
		if err != nil {
			return fmt.Errorf(`failed determine runnable tasks: %w`, err)
		}
		newTasks := len(rtr)
		if runningTasks.Val() == 0 && newTasks == 0 && !sl.IsFinished() {
			// No running tasks, no new tasks, but some tasks are still
			// waiting to run. That means a dependency loop.
			taskdump := new(strings.Builder)
			for name, task := range sl {
				fmt.Fprintf(taskdump, "%s: %s\n", name, task.GetStatus())
				for _, dep := range task.Dependencies {
					fmt.Fprintf(taskdump, "\t- %s\n", dep)
				}
			}
			return fmt.Errorf("deadlock detected: not finished, but not ready to run\n%s", taskdump.String())
		}

		// Keep track of how many tasks are in-flight.
		runningTasks.Add(newTasks)

		// This loop may be empty if there are still
		// tasks running.
		for _, task := range rtr {
			go func(s *Task) {
				defer func() {
					runningTasks.Dec()
					// Gate is unbuffered. This will block until
					// the loop comes round again.
					gate <- struct{}{}
				}()
				errInner := s.Run(handler)
				if errInner != nil {
					errors.Appendf(`error running task %q: %w`, s.Name, errInner)
					return
				}
			}(task)
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
			return fmt.Errorf(`received one or more errors running tasks, the last of which is %w`, err)
		}
	}
	return nil
}
