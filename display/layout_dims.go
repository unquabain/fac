package display

type layoutDims struct {
	maxX, maxY,
	taskGutter,
	outputWidth int
}

func newLayoutDims(maxX, maxY int) *layoutDims {
	var ld layoutDims
	ld.maxX = maxX
	ld.maxY = maxY
	ld.taskGutter = maxX / 6
	restX := maxX - ld.taskGutter
	ld.outputWidth = restX/2 - 1
	return &ld
}

func (ld *layoutDims) widgetYIterator() func() int {
	memo := 0
	return func() int {
		m := memo
		memo += 3
		return m
	}
}

func (ld *layoutDims) outputYIterator(numWidgets int) func() (int, int) {
	memo := 0
	height := (ld.maxY - 1) / numWidgets
	return func() (int, int) {
		m := memo
		memo += height + 1
		return m, height
	}
}
