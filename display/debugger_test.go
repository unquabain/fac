package display

import (
	"fmt"
	"testing"
)

func TestDims(t *testing.T) {
	d := new(Debugger)
	fmt.Fprintln(d, `ABC`)
	fmt.Fprintln(d, `ABCDE`)
	fmt.Fprintln(d, `AB`)
	fmt.Fprintln(d, `ABCDEF`)
	fmt.Fprintln(d, `ABC`)
	expectedX := 6
	expectedY := 5
	actualX, actualY := d.dims()
	if expectedX != actualX {
		t.Fatalf(`unexected width: expected %d; got %d`, expectedX, actualX)
	}
	if expectedY != actualY {
		t.Fatalf(`unexected height: expected %d; got %d`, expectedY, actualY)
	}
}

func TestDebug(t *testing.T) {
	debugger.Reset()
	debug(`Argument 1 %d`, 42)
	debug(`Argument 2 %q`, `kinder egg`)
	debug(`Argument 3 %v`, fmt.Errorf(`1202 program alarm`))
	expected := `Argument 1 42
Argument 2 "kinder egg"
Argument 3 1202 program alarm
`
	actual := debugger.String()
	if expected != actual {
		t.Fatalf("expected:\n%s\nactual:\n%s", expected, actual)
	}
}
