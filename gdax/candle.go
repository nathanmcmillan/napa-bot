package gdax

// Candle product data
type Candle struct {
	Time    int64
	Low     float64
	High    float64
	Open    float64
	Closing float64
	Volume  float64
}

// CandleList list of candles
type CandleList struct {
	product string
	list []*Candle	
}