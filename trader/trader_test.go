package trader

import "testing"

// TestTicker test move average
func TestMovingAverage(t *testing.T) {
	var move *MovingAverage

	move = NewMovingAverage(2)
	move.Rolling(3.0)
	expect(t, move.Current, 3.0)
	move.Rolling(5.0)
	expect(t, move.Current, 4.0)
	move.Rolling(7.0)
	expect(t, move.Current, 6.0)

	move = NewMovingAverage(4)
	move.Rolling(1.0)
	expect(t, move.Current, 1.0)
	move.Rolling(2.0)
	expect(t, move.Current, 1.5)
	move.Rolling(3.0)
	expect(t, move.Current, 2.0)
	move.Rolling(9.0)
	expect(t, move.Current, 3.75)
	move.Rolling(17.0)
	expect(t, move.Current, 7.75)
	move.Rolling(13.0)
	expect(t, move.Current, 10.5)
	move.Rolling(101.0)
	expect(t, move.Current, 35.0)
}

func expect(t *testing.T, actual, expected float64) {
	if actual != expected {
		t.Error("expected", expected, "got", actual)
	}
}
