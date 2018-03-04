package datastore

var (
	taker = map[string]float64{
		"BTC-USD": 0.0025,
		"ETH-USD": 0.003,
		"LTC-USD": 0.003,
	}
)

// Order an order placed on the exchange
type Order struct {
	ID int64
	Product       string
	Price         float64
	Size          float64
}

// Profit price to make a profit
func (order *Order) Profit() float64 {
	return order.Price * (1.0 + taker[order.Product])
}
