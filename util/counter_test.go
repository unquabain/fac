package util

import "testing"

func TestCounter(t *testing.T) {
	c := new(Counter)
	expect := func(expected, actual int, message string) {
		if expected != actual {
			t.Fatalf(`%s: expected %d; received %d`, message, expected, actual)
		}
	}
	c.Add(4)
	expect(4, c.Val(), `Add 4 to 0`)
	c.Add(6)
	expect(10, c.Val(), `Add 6 to 4`)
	c.Sub(2)
	expect(8, c.Val(), `Sub 2 from 10`)
	c.Inc()
	expect(9, c.Val(), `Inc 8`)
	c.Dec()
	expect(8, c.Val(), `Dec 9`)
	c.Set(22)
	expect(22, c.Val(), `Set 22`)
}
