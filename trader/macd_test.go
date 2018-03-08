package trader

import (
	"fmt"
	"testing"
)

// TestMacd test macd
func TestMacd(t *testing.T) {
	var macd *Macd

	/*
		200	200.0000 - 200.0000
		250	207.6923 - 203.7037
		230 211.1243 - 205.6516
		190	207.8744 - 204.4922
		110 192.8168 - 197.4928
		210 195.4604 - 198.4193
		240 202.3126 - 201.4994
	*/

	macd = NewMacd(12, 26)
	macd.Update(200.0)
	expectString(t, macd.Signal, "wait")
	expectString(t, fmt.Sprintf("%.4f", macd.Current), "0.0000")

	macd.Update(250.0)
	expectString(t, macd.Signal, "wait")
	expectString(t, fmt.Sprintf("%.4f", macd.Current), "3.9886")

	macd.Update(230.0)
	expectString(t, macd.Signal, "wait")
	expectString(t, fmt.Sprintf("%.4f", macd.Current), "5.4727")

	macd.Update(190.0)
	expectString(t, macd.Signal, "wait")
	expectString(t, fmt.Sprintf("%.4f", macd.Current), "3.3822")

	macd.Update(110.0)
	expectString(t, macd.Signal, "sell")
	expectString(t, fmt.Sprintf("%.4f", macd.Current), "-4.6760")

	macd.Update(210.0)
	expectString(t, macd.Signal, "wait")
	expectString(t, fmt.Sprintf("%.4f", macd.Current), "-2.9589")

	macd.Update(240.0)
	expectString(t, macd.Signal, "buy")
	expectString(t, fmt.Sprintf("%.4f", macd.Current), "0.8133")
}

func expectString(t *testing.T, actual, expected string) {
	if actual != expected {
		t.Error("expected", expected, "got", actual)
	}
}
