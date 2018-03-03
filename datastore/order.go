package datastore

const (
	taker = float64(0.003)
)

// Order an order placed on the exchange
type Order struct {
	ID            string
	Price         float64
	Size          float64
	Product       string
	Side          string
	Stp           string
	Type          string
	TimeInForce   string
	PostOnly      bool
	CreatedAt     string
	FillFees      float64
	FilledSize    float64
	ExecutedValue float64
	Status        string
	Settled       bool
}

// MinPrice price to make a profit
func MinPrice(price float64) float64 {
	return price * (1.0 + taker)
}
