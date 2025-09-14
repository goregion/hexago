package port

import (
	"context"
	"feeder/internal/entity"
)

type TradeRecordHandler interface {
	HandleTradeRecord(context.Context, *entity.TradeRecord) error
}
