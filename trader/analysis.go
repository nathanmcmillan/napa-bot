package trader

import (
	"../gdax"
)

// PercentChange percent change in value
func PercentChange(previous, now float64) float64 {
	return (now - previous) / previous * 100
}

// Average average of multiple values
func Average(values []float64) float64 {
	count := len(values)
	average := float64(0.0)
	for i := 0; i < count; i++ {
		average += values[i]
	}
	return average / float64(count)
}

// Sma simple moving average
func CreateSma(periods int64, candle []*gdax.Candle) []float64 {
	size := int64(len(candle))
	sma := make([]float64, size)
	for i := int64(0); i < size; i++ {
		if i < periods {
			sma[i] = candle[i].Closing
			continue
		}
		sum := float64(0.0)
		for j := i - periods; j < i; j++ {
			sum += candle[j].Closing
		}
		sma[i] = sum / float64(periods)
	}
	return sma
}

// Rsi relative strength index
func CreateRsi(periods int64, candle []*gdax.Candle) []float64 {
	size := int64(len(candle))
	u := make([]float64, size)
	d := make([]float64, size)
	rsi := make([]float64, size)
	if size == 0 {
		return rsi
	}
	u[0] = 0.0
	d[0] = 0.0
	rsi[0] = 0.0
	for i := int64(1); i < size; i++ {
		prev := candle[i-1].Closing
		now := candle[i].Closing
		if now > prev {
			u[i] = now - prev
			d[i] = 0.0
		} else {
			u[i] = 0.0
			d[i] = prev - now
		}
		if i < periods {
			rsi[i] = 0.0
			continue
		}
		smaU := float64(0.0)
		smaD := float64(0.0)
		for j := i - periods; j < i; j++ {
			smaU += u[j]
			smaD += d[j]
		}
		smaU /= float64(periods)
		smaD /= float64(periods)
		var rs float64
		if smaU == 0.0 || smaD == 0.0 {
			rs = 0.0
		} else {
			rs = smaU / smaD
		}
		rsi[i] = 100.0 - (100.0 / (1.0 + rs))
	}
	return rsi
}
