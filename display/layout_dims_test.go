package display

import "testing"

func TestNewLayoutDims(t *testing.T) {
	ld := newLayoutDims(240, 100)
	expect := func(property string, expected, actual int) {
		if expected != actual {
			t.Fatalf(`expected %s to be %d: was %d`, property, expected, actual)
		}
	}
	expect(`maxX`, 240, ld.maxX)
	expect(`maxY`, 100, ld.maxY)
	expect(`taskGutter`, 40, ld.taskGutter)
	expect(`outputWidth`, 99, ld.outputWidth)
}

func TestWidgetYIterator(t *testing.T) {
	ld := newLayoutDims(240, 100)
	yi := ld.widgetYIterator()
	expect := func(expected int) {
		actual := yi()
		if expected != actual {
			t.Fatalf(`expected widgetYIterator()() to return %d; returned %d`, expected, actual)
		}
	}
	expect(0)
	expect(3)
	expect(6)
	expect(9)
}

func TestOutputYIterator(t *testing.T) {
	ld := newLayoutDims(240, 101)
	yi := ld.outputYIterator(10)
	expect := func(expectedY, expectedH int) {
		actualY, actualH := yi()
		if expectedY != actualY {
			t.Fatalf(`expected Y to be %d; was %d`, expectedY, actualY)
		}
		if expectedH != actualH {
			t.Fatalf(`expected H to be %d; was %d`, expectedH, actualH)
		}
	}
	expect(0, 10)
	expect(11, 10)
	expect(22, 10)
	expect(33, 10)
}
