package adapter_redis

import (
	"strconv"

	"github.com/goregion/hexago/internal/entity"
	"github.com/goregion/must"
)

func makeOHLCStreamName(timeframeName string, symbol string) string {
	return "ohlc:" + timeframeName + ":" + symbol
}

func mustMarshalOHLC(ohlc *entity.OHLC) map[string]any {
	return map[string]any{
		"symbol":        ohlc.Symbol,
		"close_time_ms": ohlc.CloseTimeMs,
		"open":          ohlc.Open,
		"high":          ohlc.High,
		"low":           ohlc.Low,
		"close":         ohlc.Close,
	}
}

func mustUnmarshalOHLC(values map[string]any) *entity.OHLC {
	return &entity.OHLC{
		Symbol:      values["symbol"].(string),
		CloseTimeMs: must.Return(strconv.ParseInt(values["close_time_ms"].(string), 10, 64)),
		Open:        must.Return(strconv.ParseFloat(values["open"].(string), 64)),
		High:        must.Return(strconv.ParseFloat(values["high"].(string), 64)),
		Low:         must.Return(strconv.ParseFloat(values["low"].(string), 64)),
		Close:       must.Return(strconv.ParseFloat(values["close"].(string), 64)),
	}
}

func makeTickStreamKey(symbol string) string {
	return "tick:" + symbol
}

func makeTickID(timestamp int64) string {
	return strconv.FormatInt(timestamp, 10) + "-0"
}

func mustMarshalTick(message *entity.Tick) map[string]any {
	return map[string]any{
		"best_bid_price": message.BestBidPrice,
		"best_ask_price": message.BestAskPrice,
		"timestamp_ms":   message.TimestampMs,
	}
}

func mustUnmarshalTick(values map[string]any) *entity.Tick {
	return &entity.Tick{
		BestBidPrice: must.Return(
			strconv.ParseFloat(values["best_bid_price"].(string), 64),
		),
		BestAskPrice: must.Return(
			strconv.ParseFloat(values["best_ask_price"].(string), 64),
		),
		TimestampMs: must.Return(
			strconv.ParseInt(values["timestamp_ms"].(string), 10, 64),
		),
	}
}
