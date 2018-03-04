package gdax

// Settings settings for market analysis
type Settings struct {
	Products []string
	Channels []string
	TimeInterval int64
	RsiPeriods   int64
	EmaShort     int64
	EmaLong      int64
}
