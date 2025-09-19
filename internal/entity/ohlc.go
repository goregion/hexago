package entity

type OHLC struct {
	Symbol      string  `db:"symbol"`
	Open        float64 `db:"open"`
	High        float64 `db:"high"`
	Low         float64 `db:"low"`
	Close       float64 `db:"close"`
	TimestampMs int64   `db:"timestamp_ms"`
}
