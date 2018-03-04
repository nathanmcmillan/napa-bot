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
	Price         float64
	Size          float64
	ProfitPrice float64
}

func NewOrder(id int64, product string, price, size float64) *Order {
	order := &Order{}
	order.ID = id
	order.Price = price 
	order.Size = size
	order.ProfitPrice = ProfitPrice(product, price)
	return order
}

// ProfitPrice price to make a profit
func ProfitPrice(product string, price float64) float64 {
	return price * (1.0 + taker[product])
}
