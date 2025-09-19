package entity

type Tick struct {
	Symbol       string
	BestBidPrice float64
	BestAskPrice float64
	TimestampMs  int64
}
