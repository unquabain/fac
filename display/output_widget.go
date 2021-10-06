package display

import "github.com/jroimartin/gocui"

type OutputWidgetChannel string

const (
	OWCStdOut OutputWidgetChannel = `stdout`
	OWCStdErr OutputWidgetChannel = `stderr`
)

func (owc OutputWidgetChannel) String() string {
	return string(owc)
}

type OutputWidget struct {
	Widget
	Channel OutputWidgetChannel
}

func (sw *OutputWidget) viewName() string {
	return formatViewName(sw.Title, sw.Channel.String())
}

func (ow *OutputWidget) Layout(g *gocui.Gui) error {
	return ow.Widget.Layout(
		ow.viewName(),
		g,
		func(v *gocui.View) {
			if ow.Focus {
				v.BgColor = gocui.ColorBlack
				v.FgColor = gocui.ColorWhite
			} else {
				v.BgColor = 0
				v.FgColor = 0
			}
			if v.Autoscroll && ow.Focus {
				_, h := v.Size()
				l := len(v.BufferLines())
				v.Autoscroll = false
				ow.OriginY = l - h
			}
		},
	)
}

func (sow *OutputWidget) Unlayout(g *gocui.Gui) error {
	viewName := sow.viewName()
	v, err := g.View(viewName)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	if err != gocui.ErrUnknownView {
		v.Clear()
		g.DeleteView(viewName)
	}
	return nil
}
