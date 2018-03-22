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
	Product string
	List    []*Candle
}

// SortCandlesByTime sort list of candles by time
type SortCandlesByTime []*Candle

func (list SortCandlesByTime) Len() int {
	return len(list)
}
func (list SortCandlesByTime) Swap(left, right int) {
	list[left], list[right] = list[right], list[left]
}
func (list SortCandlesByTime) Less(left, right int) bool {
	return list[left].Time < list[right].Time
}
