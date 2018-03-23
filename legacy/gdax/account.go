package gdax

// Account exchange account product data
type Account struct {
	ID        string
	Currency  string
	Balance   float64
	Available float64
	Hold      float64
	ProfileID string
}
