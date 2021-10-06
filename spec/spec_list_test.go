package spec

import (
	"fmt"
	"testing"

	yaml "gopkg.in/yaml.v2"
)

var exampleYAML = `---
Clear Logs:
  command: zsh
  args:
    - rm
    - logs/development.txt
    - logs/test.txt
  expectedReturnCode: 7
  environment:
    RAILS_ENV: development
Update Bundler:
  command: bin/bundle
  args:
    - update
  dependencies:
    - Clear Logs
Only On Fail:
  command: fortune
  dependencies:
    - "! Clear Logs"
`

func getSpecListFromYaml(serialized string) (SpecList, error) {
	list := make(SpecList)
	err := yaml.Unmarshal([]byte(serialized), &list)
	if err != nil {
		return nil, fmt.Errorf(`could not unmarshal spec from YAML: %v`, err)
	}
	return list, nil
}
func TestUnmarshalYAML(t *testing.T) {
	list, err := getSpecListFromYaml(exampleYAML)
	if err != nil {
		t.Fatalf(`could not test validity of unmarshaled list: %v`, err)
	}
	if actual := len(list); actual != 3 {
		t.Fatalf(`unexpected length of list: %d`, actual)
	}

	spec, ok := list[`Clear Logs`]
	if !ok {
		t.Fatalf(`unable to find "Clear Logs"`)
	}
	if spec == nil {
		t.Fatalf(`returned spec was nil`)
	}

	if actual := spec.Name; actual != `Clear Logs` {
		t.Fatalf(`unexpected name. Expected "Clear Logs"; found %q`, actual)
	}

	if actual := spec.Command; actual != `zsh` {
		t.Fatalf(`unexpected command. Expected "zsh", was %q`, actual)
	}

	if actual := spec.Args; len(actual) != 3 {
		t.Fatalf(`unexpected args. Expected ["rm", "logs/development.txt", "logs/test.txt"]. Received %v`, actual)
	}

	if actual := len(spec.Environment); actual != 1 {
		t.Fatalf(`unexpected environment length. Expected 1; received %v`, actual)
	}

	if actual, ok := spec.Environment[`RAILS_ENV`]; !ok || actual != `development` {
		t.Fatalf(`expected environment variable not present. Expected RAILS_ENV to equal "development"; was %q`, actual)
	}

	if actual := spec.ExpectedReturnCode; actual != 7 {
		t.Fatalf(`unexpected expected return code. Expected 7. Received %v`, actual)
	}

	if actual := spec.results; actual == nil {
		t.Fatalf(`expected results not to be nil; was`)
	}

	spec, ok = list[`Update Bundler`]
	if !ok {
		t.Fatalf(`unable to find "Update Bundler"`)
	}

	if actual := spec.Name; actual != `Update Bundler` {
		t.Fatalf(`unexpected name. Expected "Update Bundler"; found %q`, actual)
	}

	if actual := spec.Dependencies; len(actual) != 1 {
		t.Fatalf(`unexpected dependenceis. Expected ["Clear Logs"]; found %v`, actual)
	}

}

func TestIsRunnable(t *testing.T) {
	list, err := getSpecListFromYaml(exampleYAML)
	if err != nil {
		t.Fatalf(`could not test IsRunnable method: %v`, err)
	}

	clearLogs, ok := list[`Clear Logs`]
	if !ok {
		t.Fatalf(`did not find spec entry for "Clear Logs"`)
	}
	updateBundler, ok := list[`Update Bundler`]
	if !ok {
		t.Fatalf(`did not find spec entry for "Update Bundler"`)
	}
	onlyOnFail, ok := list[`Only On Fail`]
	if !ok {
		t.Fatalf(`did not find spec entry for "Only On Fail"`)
	}

	// "Clear Logs" has no dependencies and has not yet been run.
	// Should be runnable.
	if actual, err := list.IsRunnable(clearLogs); err != nil || !actual {
		if err != nil {
			t.Fatalf(`couldn't tell if "Clear Logs" was runnable: %v`, err)
		}
		t.Fatalf(`expected "Clear Logs" to be runnable; wasn't`)
	}

	// "Update Bundler" depends on "Clear Logs", which has not been run
	// so it is not runnable.
	if actual, err := list.IsRunnable(updateBundler); err != nil || actual {
		if err != nil {
			t.Fatalf(`couldn't tell if "Update Bundler" was runnable: %v`, err)
		}
		t.Fatalf(`expected "Update Bundler" not to be runnable; was`)
	}

	clearLogs.results.SetStatus(StatusFailed)

	// "Clear Logs" has been run, but failed. So "Update Bundler" cannot
	// be run.
	if actual, err := list.IsRunnable(updateBundler); err != nil || actual {
		if err != nil {
			t.Fatalf(`couldn't tell if "Update Bundler" was runnable: %v`, err)
		}
		t.Fatalf(`expected "Update Bundler" not to be runnable; was`)
	}

	clearLogs.results.SetStatus(StatusSucceeded)

	// "Clear Logs" has succeeded, but the last test marked "Update Bundler"
	// as having failed dependencies, so it's still not runnable.
	if actual, err := list.IsRunnable(updateBundler); err != nil || actual {
		if err != nil {
			t.Fatalf(`couldn't tell if "Update Bundler" was runnable: %v`, err)
		}
		t.Fatalf(`expected "Update Bundler" not to be runnable; was`)
	}

	// Reset "Update Bundler"
	updateBundler.results.SetStatus(StatusNotRun)

	// NOW, "Update Bundler" should be runnable, because it has not been
	// run, but its dependency, "Clear Logs", has run successfully.
	if actual, err := list.IsRunnable(updateBundler); err != nil || !actual {
		if err != nil {
			t.Fatalf(`couldn't tell if "Update Bundler" was runnable: %v`, err)
		}
		t.Fatalf(`expected "Update Bundler" to be runnable; wasn't`)
	}

	// Only On Fail is only runnable if Clear Logs failed.
	clearLogs.results.SetStatus(StatusNotRun)
	onlyOnFail.results.SetStatus(StatusNotRun)
	if actual, err := list.IsRunnable(onlyOnFail); err != nil || actual {
		if err != nil {
			t.Fatalf(`couldn't tell if "Only On Fail" was runnable: %v`, err)
		}
		t.Fatalf(`expected "Only On Fail" not to be runnable; was`)
	}
	clearLogs.results.SetStatus(StatusDependenciesNotMet)
	onlyOnFail.results.SetStatus(StatusNotRun)
	if actual, err := list.IsRunnable(onlyOnFail); err != nil || actual {
		if err != nil {
			t.Fatalf(`couldn't tell if "Only On Fail" was runnable: %v`, err)
		}
		t.Fatalf(`expected "Only On Fail" not to be runnable; was`)
	}
	clearLogs.results.SetStatus(StatusSucceeded)
	onlyOnFail.results.SetStatus(StatusNotRun)
	if actual, err := list.IsRunnable(onlyOnFail); err != nil || actual {
		if err != nil {
			t.Fatalf(`couldn't tell if "Only On Fail" was runnable: %v`, err)
		}
		t.Fatalf(`expected "Only On Fail" not to be runnable; was`)
	}
	clearLogs.results.SetStatus(StatusFailed)
	onlyOnFail.results.SetStatus(StatusNotRun)
	if actual, err := list.IsRunnable(onlyOnFail); err != nil || !actual {
		if err != nil {
			t.Fatalf(`couldn't tell if "Only On Fail" was runnable: %v`, err)
		}
		t.Fatalf(`expected "Only On Fail" to be runnable; wasn't`)
	}
}

func TestReadyToRun(t *testing.T) {
	list, err := getSpecListFromYaml(exampleYAML)
	if err != nil {
		t.Fatalf(`could not test ReadyToRun method: %v`, err)
	}
	rtr, err := list.ReadyToRun()
	if err != nil {
		t.Fatalf(`could not get ready-to-run specs: %v`, err)
	}
	if actual := len(rtr); actual != 1 {
		t.Fatalf(`expected one ready-to-run spec: Found %d`, actual)
	}
	if actual := rtr[0].Name; actual != `Clear Logs` {
		t.Fatalf(`expected "Clear Logs" to be ready, but %q was`, actual)
	}

	list[`Clear Logs`].results.SetStatus(StatusSucceeded)
	list[`Update Bundler`].results.SetStatus(StatusNotRun)

	rtr, err = list.ReadyToRun()
	if err != nil {
		t.Fatalf(`could not get ready-to-run specs: %v`, err)
	}
	if actual := len(rtr); actual != 1 {
		t.Fatalf(`expected one ready-to-run spec: Found %d`, actual)
	}
	if actual := rtr[0].Name; actual != `Update Bundler` {
		t.Fatalf(`expected "Update Bundler" to be ready, but %q was`, actual)
	}
}

func TestIsFinished(t *testing.T) {
	list, err := getSpecListFromYaml(exampleYAML)
	if err != nil {
		t.Fatalf(`could not test IsFinished method: %v`, err)
	}
	if list.IsFinished() {
		t.Fatalf(`expected list not to be finished, but was`)
	}
	list[`Clear Logs`].results.SetStatus(StatusFailed)
	if list.IsFinished() {
		t.Fatalf(`expected list not to be finished, but was`)
	}
	list[`Only On Fail`].results.SetStatus(StatusSucceeded)
	list[`Update Bundler`].results.SetStatus(StatusDependenciesNotMet)
	if !list.IsFinished() {
		t.Fatalf(`expected list to be finished, but was not`)
	}
}
