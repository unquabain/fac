package display

import (
	"fmt"
	"sort"

	"github.com/Unquabain/thing-doer/spec"
	"github.com/jroimartin/gocui"
)

type stringerAdapter string

func (sa stringerAdapter) String() string {
	return string(sa)
}

// FocusColumn is an enum for determining which of the
// three columns should receive keyboard input.
type FocusColumn int

const (
	FCTaskList FocusColumn = iota
	FCStdOut
	FCStdErr
)

// SpecLayoutManager is the main program layout manager.
// It divides the screen into three columns: the left for
// the Spec list, and the remaining are STDOUT and STDIN
// for the running specs. While specs are running the
// STDOUT and STDIN space is divided vertically for the
// running specs. Once it's done, you can use the arrow
// keys to navigate through the completed specs and
// examine their output.
type SpecLayoutManager struct {
	spec.SpecList
	IsFinished bool
	FocusColumn
	FocusRow      int
	outputWidgets OutputWidgetRegistry
}

func (slm *SpecLayoutManager) sorted() []*spec.Spec {
	slice := make([]*spec.Spec, len(slm.SpecList))
	i := 0
	for _, s := range slm.SpecList {
		slice[i] = s
		i++
	}
	sort.Slice(slice, func(i, j int) bool {
		return slice[i].Order < slice[j].Order
	})
	return slice
}

func (slm *SpecLayoutManager) showConsole(pos int, s *spec.Spec) bool {
	if slm.IsFinished {
		return slm.FocusRow == pos
	}
	return s.GetStatus() == spec.StatusRunning
}

// Update enques a re-lay-out the screen from a goroutine.
func (slm *SpecLayoutManager) Update(g *gocui.Gui) {
	g.Update(func(gg *gocui.Gui) error {
		slm.Layout(gg)
		return nil
	})
}

// ArrowUp updates the internal state in response to
// an Arrow Up keyboard event in the task list.
func (slm *SpecLayoutManager) ArrowUp() {
	l := len(slm.SpecList)
	slm.FocusRow = (l + slm.FocusRow - 1) % l
}

// ArrowDown updates the internal state in response to
// an Arrow Down keyboard event in the task list.
func (slm *SpecLayoutManager) ArrowDown() {
	l := len(slm.SpecList)
	slm.FocusRow = (slm.FocusRow + 1) % l
}

// SetFocusTaskList sets the internal state to focus on the
// individual tasks in the task list.
func (slm *SpecLayoutManager) SetFocusTaskList() {
	slm.FocusColumn = FCTaskList
}

// SetFocusStdOut sets the internal state to display the
// STDOUT pane as focused.
func (slm *SpecLayoutManager) SetFocusStdOut() {
	slm.FocusColumn = FCStdOut
}

// SetFocusStdErr sets the internal state to display the
// STDERR pane as focused.
func (slm *SpecLayoutManager) SetFocusStdErr() {
	slm.FocusColumn = FCStdErr
}

func (slm *SpecLayoutManager) setStatusKeybindings(w *StatusWidget, g *gocui.Gui) {
	g.DeleteKeybindings(w.viewName())
	g.SetKeybinding(
		w.viewName(),
		gocui.KeyArrowUp,
		gocui.ModNone,
		func(gg *gocui.Gui, v *gocui.View) error {
			slm.ArrowUp()
			slm.Update(gg)
			return nil
		},
	)
	g.SetKeybinding(
		w.viewName(),
		gocui.KeyArrowDown,
		gocui.ModNone,
		func(gg *gocui.Gui, v *gocui.View) error {
			slm.ArrowDown()
			slm.Update(gg)
			return nil
		},
	)
	g.SetKeybinding(
		w.viewName(),
		gocui.KeyArrowRight,
		gocui.ModNone,
		func(gg *gocui.Gui, v *gocui.View) error {
			slm.SetFocusStdOut()
			slm.Update(gg)
			return nil
		},
	)
}

func (slm *SpecLayoutManager) setStdoutKeybindings(sow *OutputWidget, g *gocui.Gui) {
	g.DeleteKeybindings(sow.viewName())
	g.SetKeybinding(
		sow.viewName(),
		gocui.KeyArrowLeft,
		gocui.ModNone,
		func(gg *gocui.Gui, v *gocui.View) error {
			slm.SetFocusTaskList()
			slm.Update(gg)
			return nil
		},
	)
	g.SetKeybinding(
		sow.viewName(),
		gocui.KeyArrowRight,
		gocui.ModNone,
		func(gg *gocui.Gui, v *gocui.View) error {
			slm.SetFocusStdErr()
			slm.Update(gg)
			return nil
		},
	)
	g.SetKeybinding(
		sow.viewName(),
		gocui.KeyArrowDown,
		gocui.ModNone,
		func(gg *gocui.Gui, v *gocui.View) error {
			sow.CursorDown()
			slm.Update(gg)
			return nil
		},
	)
	g.SetKeybinding(
		sow.viewName(),
		gocui.KeyPgdn,
		gocui.ModNone,
		func(gg *gocui.Gui, v *gocui.View) error {
			sow.PageDown()
			slm.Update(gg)
			return nil
		},
	)
	g.SetKeybinding(
		sow.viewName(),
		gocui.KeyArrowUp,
		gocui.ModNone,
		func(gg *gocui.Gui, v *gocui.View) error {
			sow.CursorUp()
			slm.Update(gg)
			return nil
		},
	)
	g.SetKeybinding(
		sow.viewName(),
		gocui.KeyPgup,
		gocui.ModNone,
		func(gg *gocui.Gui, v *gocui.View) error {
			sow.PageUp()
			slm.Update(gg)
			return nil
		},
	)
	g.SetKeybinding(
		sow.viewName(),
		gocui.KeyHome,
		gocui.ModNone,
		func(gg *gocui.Gui, v *gocui.View) error {
			sow.Home()
			slm.Update(gg)
			return nil
		},
	)
}

func (slm *SpecLayoutManager) setStderrKeybindings(sew *OutputWidget, g *gocui.Gui) {
	g.DeleteKeybindings(sew.viewName())
	g.SetKeybinding(
		sew.viewName(),
		gocui.KeyArrowLeft,
		gocui.ModNone,
		func(gg *gocui.Gui, v *gocui.View) error {
			slm.SetFocusStdOut()
			slm.Update(gg)
			return nil
		},
	)
	g.SetKeybinding(
		sew.viewName(),
		gocui.KeyArrowDown,
		gocui.ModNone,
		func(gg *gocui.Gui, v *gocui.View) error {
			sew.CursorDown()
			slm.Update(gg)
			return nil
		},
	)
	g.SetKeybinding(
		sew.viewName(),
		gocui.KeyPgdn,
		gocui.ModNone,
		func(gg *gocui.Gui, v *gocui.View) error {
			sew.PageDown()
			slm.Update(gg)
			return nil
		},
	)
	g.SetKeybinding(
		sew.viewName(),
		gocui.KeyArrowUp,
		gocui.ModNone,
		func(gg *gocui.Gui, v *gocui.View) error {
			sew.CursorUp()
			slm.Update(gg)
			return nil
		},
	)
	g.SetKeybinding(
		sew.viewName(),
		gocui.KeyPgup,
		gocui.ModNone,
		func(gg *gocui.Gui, v *gocui.View) error {
			sew.PageUp()
			slm.Update(gg)
			return nil
		},
	)
	g.SetKeybinding(
		sew.viewName(),
		gocui.KeyHome,
		gocui.ModNone,
		func(gg *gocui.Gui, v *gocui.View) error {
			sew.Home()
			slm.Update(gg)
			return nil
		},
	)
}

// Layout satisfies the gocui.Manager interface.
// The main drawing logic of the manager.
func (slm *SpecLayoutManager) Layout(g *gocui.Gui) error {
	debugger.Reset()
	if !slm.IsFinished {
		slm.IsFinished = slm.SpecList.IsFinished()
	}
	dims := newLayoutDims(g.Size())
	sorted := slm.sorted()

	stdoutWidgets := make([]*OutputWidget, 0, len(sorted))
	stderrWidgets := make([]*OutputWidget, 0, len(sorted))
	if slm.outputWidgets == nil {
		slm.outputWidgets = make(OutputWidgetRegistry)
	}
	taskGutterY := dims.widgetYIterator()
	for pos, task := range slm.sorted() {
		w := newStatusWidget(task, dims.taskGutter, taskGutterY)
		g.DeleteKeybindings(w.viewName())
		var (
			sow *OutputWidget
			sew *OutputWidget
		)
		sow = slm.outputWidgets.makeStdOutWidget(dims, task)
		sew = slm.outputWidgets.makeStdErrWidget(dims, task)
		if slm.showConsole(pos, task) {
			stdoutWidgets = append(stdoutWidgets, sow)
			stderrWidgets = append(stderrWidgets, sew)
			sow.Attribute = 0
			sew.Attribute = 0
			if slm.IsFinished {
				w.Focus = true
				if slm.FocusColumn == FCStdOut {
					sow.Focus = true
				} else {
					sow.Focus = false
				}
				if slm.FocusColumn == FCStdErr {
					sew.Focus = true
				} else {
					sew.Focus = false
				}
				defer func(w *StatusWidget, sow *OutputWidget, sew *OutputWidget) {
					slm.setStatusKeybindings(w, g)
					slm.setStdoutKeybindings(sow, g)
					slm.setStderrKeybindings(sew, g)
					viewName := ``
					switch slm.FocusColumn {
					case FCTaskList:
						viewName = w.viewName()
					case FCStdOut:
						viewName = sow.viewName()
					case FCStdErr:
						viewName = sew.viewName()
					}
					g.SetCurrentView(viewName)
				}(w, sow, sew)
			} else {
				w.Focus = false
			}
		} else {
			defer func(sow *OutputWidget, sew *OutputWidget) {
				sow.Unlayout(g)
				sew.Unlayout(g)
			}(sow, sew)
		}
		w.Layout(g)
	}

	if len(stdoutWidgets) == 0 {
		return nil
	}

	outputY := dims.outputYIterator(len(stdoutWidgets))
	for _, sow := range stdoutWidgets {
		sow.Y, sow.H = outputY()
		err := sow.Layout(g)
		if err != nil {
			return fmt.Errorf(`couldn't layout stdout for %q: %w`, sow.Title, err)
		}
	}

	outputY = dims.outputYIterator(len(stdoutWidgets))
	for _, sew := range stderrWidgets {
		sew.Y, sew.H = outputY()
		err := sew.Layout(g)
		if err != nil {
			return fmt.Errorf(`couldn't layout stderr for %q: %w`, sew.Title, err)
		}
	}

	if debugger.Len() > 0 {
		debugger.Layout(g)
	}

	return nil
}
