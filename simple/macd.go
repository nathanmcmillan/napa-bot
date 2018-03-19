package main

type macd struct {
	short   *ema
	long    *ema
	current float64
	signal  string
}

func newMacd(short int64, long int64, initial float64) *macd {
	m := &macd{}
	m.short = newEma(short, initial)
	m.long = newEma(long, initial)
	m.signal = "wait"
	return m
}

func (m *macd) update(closing float64) {
	m.short.update(closing)
	m.long.update(closing)
	before := m.current
	m.current = m.short.current - m.long.current
	if before < 0 && m.current > 0 {
		m.signal = "buy"
	} else if before > 0 && m.current < 0 {
		m.signal = "sell"
	} else {
		m.signal = "wait"
	}
	m.signal = "buy"
}
