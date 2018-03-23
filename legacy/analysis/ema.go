package analysis

// ExponentialMovingAverage exponential moving average
type ExponentialMovingAverage struct {
	Periods int64
	Weight  float64
	Current float64
}

// NewEma constructor
func NewEma(periods int64, initial float64) *ExponentialMovingAverage {
	ema := &ExponentialMovingAverage{}
	ema.Periods = periods
	ema.Weight = 2.0 / (float64(periods) + 1.0)
	ema.Current = initial
	return ema
}

// Update update exponential moving average
func (ema *ExponentialMovingAverage) Update(value float64) {
	ema.Current = (value-ema.Current)*ema.Weight + ema.Current
}
