package gdax

var (
	taker = map[string]float64{
		"BTC-USD": 0.0025,
		"ETH-USD": 0.003,
		"LTC-USD": 0.003,
	}
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

// ProfitPrice calculates the price needed to make a profit
func ProfitPrice(product string, price float64) float64 {
	return price * (1.0 + taker[product])
}
