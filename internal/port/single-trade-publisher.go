package port

import (
	"context"
	"feeder/internal/entity"
)

type SingleTradePublisher interface {
	PublishSingleTrade(ctx context.Context, trade *entity.SingleTrade) error
}
