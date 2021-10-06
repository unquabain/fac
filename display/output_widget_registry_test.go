package display

import (
	"fmt"
	"testing"

	"github.com/Unquabain/thing-doer/spec"
	"gopkg.in/yaml.v2"
)

var sampleYAML = `---
Test 1:
  command: cat
Test 2:
  command: echo
`

func newSpecList() (spec.SpecList, error) {
	list := make(spec.SpecList)
	err := yaml.Unmarshal([]byte(sampleYAML), &list)
	if err != nil {
		return nil, fmt.Errorf(`could not create example SpecList: %w`, err)
	}
	return list, nil
}

func TestOutputWidgetRegistryTest(t *testing.T) {
	ld := newLayoutDims(240, 100)
	owr := make(OutputWidgetRegistry)
	list, err := newSpecList()
	if err != nil {
		t.Fatalf(`Could not initialize test: %v`, err)
	}
	task1, ok := list[`Test 1`]
	if !ok {
		t.Fatalf(`Could not find Test 1 in deserialized list: %v`, list)
	}
	task2, ok := list[`Test 2`]
	if !ok {
		t.Fatalf(`Could not find Test 2 in deserialized list: %v`, list)
	}
	var expected *OutputWidget
	var actual *OutputWidget
	expected = owr.makeStdOutWidget(ld, task1)
	actual = owr.makeStdOutWidget(ld, task1)

	if expected != actual {
		t.Fatal(`expected to receive the same widget back: received different`)
	}

	actual = owr.makeStdErrWidget(ld, task1)
	if expected == actual {
		t.Fatal(`expected to receive different widget back: received the same`)
	}

	actual = owr.makeStdOutWidget(ld, task2)
	if expected == actual {
		t.Fatal(`expected to receive different widget back: received the same`)
	}
}
