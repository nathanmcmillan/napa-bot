package gdax

// Ticker last match from exchange
type Ticker struct {
	Time      int64
	ProductID string
	Price     float64
	Side      string
}
