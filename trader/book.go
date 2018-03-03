package trader

import "../gdax"

// Book track order book
type Book struct {
	Bids []gdax.SnapshotTuple
	Asks []gdax.SnapshotTuple
}

// NewBook constructor
func NewBook() *Book {
	book := Book{}
	return &book
}

// Snapshot initialize order book
func (book *Book) Snapshot(snapshot *gdax.Snapshot) {

}

// Update update order book
func (book *Book) Update(update *gdax.Update) {

}
