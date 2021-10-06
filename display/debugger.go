package display

import (
	"fmt"
	"strings"

	"github.com/jroimartin/gocui"
)

type Debugger struct{ strings.Builder }

func (d *Debugger) dims() (int, int) {
	msg := d.String()
	lines := strings.Split(msg, "\n")
	y, x := 0, 0
	for i, line := range lines {
		y = i
		xl := len(line)
		if xl > x {
			x = xl
		}
	}
	return x, y
}

func (d *Debugger) Layout(g *gocui.Gui) {
	windowW, windowH := g.Size()
	widgetW, widgetH := d.dims()
	v, err := g.SetView(
		`debug`,
		windowW-widgetW-1, windowH-widgetH-1,
		windowW-1, windowH-1,
	)
	if err != nil && err != gocui.ErrUnknownView {
		fmt.Println(d.String())
		return
	}
	v.Clear()
	v.Title = `Debug`
	fmt.Fprint(v, d.String())
	g.SetViewOnTop(`debug`)
}

var debugger *Debugger = new(Debugger)

func debug(format string, args ...interface{}) {
	if debugger == nil {
		debugger = new(Debugger)
	}
	fmt.Fprintln(debugger, fmt.Sprintf(format, args...))
}
