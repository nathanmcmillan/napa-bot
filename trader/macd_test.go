package trader

import "testing"

// TestMacd test macd
func TestMacd(t *testing.T) {
    var macd *Macd

    macd = NewMacd(12, 24, 230)
	expect(t, macd.Current, 230)
    expectString(t, macd.Signal, "wait")
	macd.Update(250.0)
	expect(t, macd.Current, 3.0)
    expectString(t, macd.Signal, "buy")
}

func expectString(t *testing.T, actual, expected string) {
	if actual != expected {
		t.Error("expected", expected, "got", actual)
	}
}
