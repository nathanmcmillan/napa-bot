package trader

import "testing"

// TestEma test ema
func TestEma(t *testing.T) {
	var ema *ExponentialMovingAverage

	ema = NewEma(12, 200)
	ema.Update(210.0)
	expect(t, ema.Current, 199.0)
}