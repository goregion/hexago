package adapter_grpc_api

import (
	"github.com/goregion/hexago/internal/adapter/grpc-api/gen"
	"github.com/goregion/hexago/internal/entity"
)

func mustMarshalOHLC(ohlc *entity.OHLC) *gen.OHLC {
	if ohlc == nil {
		return nil
	}
	return &gen.OHLC{
		Open:         ohlc.Open,
		Close:        ohlc.Close,
		High:         ohlc.High,
		Low:          ohlc.Low,
		ClosesTimeMs: ohlc.CloseTimeMs,
	}
}
