package analyst

import (
	"../historian"
)

// Analyst market data
type Analyst struct {
	TimeInterval  int64
	RsiPeriods    int64
	Rsi           []float64
	EmaShort      int64
	EmaLong       int64
	Macd          []float64
	Ticker        []float64
	Book          []float64
	TickerAverage float64
	BuyAverage    float64
	SellAverage   float64
}

// PercentChange percent change in value
func PercentChange(previous, now float64) float64 {
	return (now - previous) / previous * 100
}

// UpdateMovingAverage updates a given average by one
func UpdateMovingAverage(current, update float64, periods int64) float64 {
	return (current*float64(periods) + update) / float64(periods+1)
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
func Sma(periods int64, candle []*historian.Candle) []float64 {
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

// UpdateSma update the average by one
func UpdateSma(sma float64, periods int64, candle historian.Candle) float64 {
	return 0.0
}

// Ema exponential moving average
func Ema(periods int64, candle []*historian.Candle) []float64 {
	size := int64(len(candle))
	ema := make([]float64, size)
	weight := 2.0 / (float64(periods) + 1.0)
	for i := int64(0); i < size; i++ {
		if i < periods {
			ema[i] = candle[i].Closing
			continue
		}
		previous := ema[i-1]
		ema[i] = (candle[i].Closing-previous)*weight + previous
	}
	return ema
}

// UpdateEma update the average by one
func UpdateEma(ema float64, periods int64, candle historian.Candle) float64 {
	weight := 2.0 / (float64(periods) + 1.0)
	return (candle.Closing-ema)*weight + ema
}

// Macd moving average convergence divergence
func Macd(periodsA int64, periodsB int64, candle []*historian.Candle) []float64 {
	emaA := Ema(periodsA, candle)
	emaB := Ema(periodsB, candle)
	size := len(candle)
	macd := make([]float64, size)
	for i := 0; i < size; i++ {
		macd[i] = emaA[i] - emaB[i]
	}
	return macd
}

// Rsi relative strength index
func Rsi(periods int64, candle []*historian.Candle) []float64 {
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
