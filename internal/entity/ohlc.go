package entity

type OHLC struct {
	Symbol      string  `db:"symbol"`
	Open        float64 `db:"open"`
	High        float64 `db:"high"`
	Low         float64 `db:"low"`
	Close       float64 `db:"close"`
	CloseTimeMs int64   `db:"close_time_ms"`
}
