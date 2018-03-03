package gdax

// Snapshot snapshot of level 2 from exchange
type Snapshot struct {
	ProductID string
	Bids      []SnapshotTuple
	Asks      []SnapshotTuple
}

// SnapshotTuple tuple of price and size
type SnapshotTuple struct {
	Price float64
	Size  float64
}

// Update level 2 update from exchange
type Update struct {
	ProductID string
	Changes   []UpdateChange
}

// UpdateChange level 2 update details
type UpdateChange struct {
	Side  string
	Price float64
	Size  float64
}
