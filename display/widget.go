package display

import (
	"fmt"
	"strings"

	"github.com/Unquabain/fac/util"
	"github.com/jroimartin/gocui"
)

func formatViewName(title string, layoutType string) string {
	var builder strings.Builder

	builder.WriteString(util.Parameterize(title))
	builder.WriteRune('/')
	builder.WriteString(layoutType)
	return builder.String()
}

// Widget encapsulates a sub-view in the text interface.
// It's an abstract base for other, differentiated widget
// types.
type Widget struct {
	Title      string
	Stringer   fmt.Stringer
	Focus      bool
	Attribute  gocui.Attribute
	X, Y, H, W int
	OriginY    int
}

// Layout does NOT satisfy the gocui.Manager interface, but contains
// the common logic for other widget types that do.
func (w *Widget) Layout(viewName string, g *gocui.Gui, customize func(*gocui.View)) error {
	v, err := g.SetView(
		viewName,
		w.X, w.Y,
		w.X+w.W, w.Y+w.H,
	)
	if err != nil && err != gocui.ErrUnknownView {
		return fmt.Errorf(`cannot get view for %q (%q): %w`, w.Title, viewName, err)
	}
	if w.Focus {
		v.Title = fmt.Sprintf(`[%s]`, w.Title)
	} else {
		v.Title = fmt.Sprintf(` %s `, w.Title)
	}
	customize(v)
	v.Clear()
	fmt.Fprintf(v, ` %s`, w.Stringer)
	v.Frame = true
	ox, _ := v.Origin()
	err = v.SetOrigin(ox, w.OriginY)
	if err != nil {
		return fmt.Errorf(`couldn't set origin of %q: %w`, w.Title, err)
	}
	return nil
}

// For scrolling widgets, updates the internal scroll position.
func (w *Widget) CursorDown() {
	w.OriginY += 1
}

// For scrolling widgets, updates the internal scroll position.
func (w *Widget) PageDown() {
	w.OriginY += 10
}

// For scrolling widgets, updates the internal scroll position.
func (w *Widget) CursorUp() {
	w.OriginY -= 1
	if w.OriginY < 0 {
		w.OriginY = 0
	}
}

// For scrolling widgets, updates the internal scroll position.
func (w *Widget) PageUp() {
	w.OriginY -= 10
	if w.OriginY < 0 {
		w.OriginY = 0
	}
}

// For scrolling widgets, updates the internal scroll position.
func (w *Widget) Home() {
	w.OriginY = 0
}
