package analysis

import (
	"fmt"
	"testing"
)

// TestEma test ema
func TestEma(t *testing.T) {
	var ema *ExponentialMovingAverage

	ema = NewEma(12, 200)
	ema.Update(210.0)
	expectString(t, fmt.Sprintf("%.4f", ema.Current), "201.5385")
	ema.Update(250.0)
	expectString(t, fmt.Sprintf("%.4f", ema.Current), "208.9941")
	ema.Update(230.0)
	expectString(t, fmt.Sprintf("%.4f", ema.Current), "212.2258")
	ema.Update(190.0)
	expectString(t, fmt.Sprintf("%.4f", ema.Current), "208.8064")
	ema.Update(110.0)
	expectString(t, fmt.Sprintf("%.4f", ema.Current), "193.6054")
}
