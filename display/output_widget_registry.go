package display

import (
	"fmt"

	"github.com/Unquabain/thing-doer/spec"
)

type registryKey struct {
	Title   string
	Channel OutputWidgetChannel
}

// OutputWidgetRegistry keeps track of all the widgets used,
// even when they don't have current views in memory, in order
// to preserve state.
type OutputWidgetRegistry map[registryKey]*OutputWidget

func (r OutputWidgetRegistry) makeOutputWidget(dims *layoutDims, task *spec.Spec, channel OutputWidgetChannel) *OutputWidget {
	key := registryKey{Title: task.Name, Channel: channel}
	sow, ok := r[key]
	if !ok {
		sow = new(OutputWidget)
		sow.Title = fmt.Sprintf(`%s (%s)`, task.Name, channel)
		sow.Channel = channel
		r[key] = sow
	}
	sow.W = dims.outputWidth
	sow.X = dims.taskGutter + 1
	switch channel {
	case OWCStdOut:
		sow.Stringer = stringerAdapter(task.GetStdOut())
	case OWCStdErr:
		sow.Stringer = stringerAdapter(task.GetStdErr())
		sow.X += dims.outputWidth + 1
	}
	return sow
}

func (r OutputWidgetRegistry) makeStdOutWidget(dims *layoutDims, task *spec.Spec) *OutputWidget {
	return r.makeOutputWidget(dims, task, OWCStdOut)
}

func (r OutputWidgetRegistry) makeStdErrWidget(dims *layoutDims, task *spec.Spec) *OutputWidget {
	return r.makeOutputWidget(dims, task, OWCStdErr)
}
