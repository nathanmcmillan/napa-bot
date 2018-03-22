package main

type ema struct {
	periods int64
	weight  float64
	current float64
}

func newEma(periods int64, initial float64) *ema {
	e := &ema{}
	e.periods = periods
	e.weight = 2.0 / (float64(periods) + 1.0)
	e.current = initial
	return e
}

func (e *ema) update(value float64) {
	e.current = (value-e.current)*e.weight + e.current
}
