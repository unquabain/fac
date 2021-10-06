package display

import (
	"github.com/Unquabain/thing-doer/spec"
	"github.com/jroimartin/gocui"
)

type StatusWidget Widget

func newStatusWidget(task *spec.Spec, width int, yIter func() int) *StatusWidget {
	w := new(StatusWidget)
	w.Title = task.Name
	w.X = 0
	w.Y = yIter()
	w.H = 2
	w.W = width
	status := task.GetStatus()
	w.Stringer = status
	if !status.IsOK() {
		w.Attribute = gocui.ColorRed
	} else if status == spec.StatusRunning {
		w.Attribute = gocui.ColorYellow
	} else if status == spec.StatusSucceeded {
		w.Attribute = gocui.ColorGreen
	}
	w.Focus = false
	return w
}

func (sw *StatusWidget) viewName() string {
	return formatViewName(sw.Title, `status`)
}

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
