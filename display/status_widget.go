package display

import (
	"github.com/Unquabain/fac/task"
	"github.com/jroimartin/gocui"
)

// StatusWidget is a widget that displays the current status of
// a Task in the left-hand column.
type StatusWidget Widget

func newStatusWidget(t *task.Task, width int, yIter func() int) *StatusWidget {
	w := new(StatusWidget)
	w.Title = t.Name
	w.X = 0
	w.Y = yIter()
	w.H = 2
	w.W = width
	status := t.GetStatus()
	w.Stringer = status
	if !status.IsOK() {
		w.Attribute = gocui.ColorRed
	} else if status == task.StatusRunning {
		w.Attribute = gocui.ColorYellow
	} else if status == task.StatusSucceeded {
		w.Attribute = gocui.ColorGreen
	}
	w.Focus = false
	return w
}

func (sw *StatusWidget) viewName() string {
	return formatViewName(sw.Title, `status`)
}

// Layout satisfies the gocui.Manager interface, and
// contains the graphical logic.
func (sw *StatusWidget) Layout(g *gocui.Gui) error {
	return (*Widget)(sw).Layout(
		sw.viewName(),
		g,
		func(v *gocui.View) {
			if sw.Focus {
				v.BgColor = sw.Attribute
				v.FgColor = gocui.ColorBlack
			} else {
				v.BgColor = 0
				v.FgColor = sw.Attribute
			}
		},
	)
}
