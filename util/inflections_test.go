package util

import "testing"

func TestParameterize(t *testing.T) {
	expect := func(input, expected string) {
		actual := Parameterize(input)
		if actual != expected {
			t.Fatalf(`Parameterize %q should have been %q; was %q`, input, expected, actual)
		}
	}

	expect(`test`, `test`)
	expect(`Test`, `test`)
	expect(`Std Out`, `std-out`)
	expect(`Solve All Problems`, `solve-all-problems`)
}
