package gdax

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
