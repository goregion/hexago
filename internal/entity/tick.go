package entity

type Tick struct {
	Symbol       string
	BestBidPrice float64
	BestAskPrice float64
	TimestampMs  int64
}

type TickRange struct {
	Symbol    string
	FromMs    int64
	ToMs      int64
	TickSlice []*Tick
}
