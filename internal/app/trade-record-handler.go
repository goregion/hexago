package app

import (
	"context"
	"feeder/internal/entity"
	"feeder/internal/port"
)

type TradeRecordHandler struct {
	port.TradeRecordHandler
	singleTradePublisher port.SingleTradePublisher
}

func NewTradeRecordHandler(singleTradePublisher port.SingleTradePublisher) *TradeRecordHandler {
	return &TradeRecordHandler{
		singleTradePublisher: singleTradePublisher,
	}
}

func (h *TradeRecordHandler) HandleTradeRecord(ctx context.Context, record *entity.TradeRecord) error {
	if record.Operation == entity.OperationType_REMOVE {
		return nil
	}

	var direction entity.SingleTrade_Direction
	switch record.Cmd {
	case entity.TradeRecord_BUY:
		direction = entity.SingleTrade_B
	case entity.TradeRecord_SELL:
		direction = entity.SingleTrade_S
	default:
		return nil
	}

	var price = record.OpenPrice
	var timestamp = record.OpenTimeMs
	if record.CloseTimeMs != 0 {
		price = record.ClosePrice
		timestamp = record.CloseTimeMs
		if direction == entity.SingleTrade_B {
			direction = entity.SingleTrade_S
		}
	}

	var singleTrade = &entity.SingleTrade{
		Direction:      direction,
		Price:          price,
		Timestamp:      timestamp,
		InstrumentName: record.Symbol,
		Amount:         int64(record.Volume * record.SymbolCsize),
		Currency1:      record.SymbolCurrency1,
		Currency2:      record.SymbolCurrency2,
	}

	return h.singleTradePublisher.PublishSingleTrade(ctx, singleTrade)
}
