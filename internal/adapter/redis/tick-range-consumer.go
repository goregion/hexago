package adapter_redis

import (
	"context"
	"time"

	"github.com/goregion/hexago/internal/entity"
	"github.com/goregion/hexago/internal/port"
	"github.com/goregion/hexago/pkg/goter"
	"github.com/goregion/hexago/pkg/redis"
	"github.com/pkg/errors"
)

type TickRangeConsumer struct {
	redisClient *redis.Client
	consumer    port.TickRangeConsumer
}

func NewTickRangeConsumer(redisClient *redis.Client, consumer port.TickRangeConsumer) *TickRangeConsumer {
	return &TickRangeConsumer{
		redisClient: redisClient,
		consumer:    consumer,
	}
}

func (h *TickRangeConsumer) readNext(ctx context.Context, symbol string, from time.Time, to time.Time) error {
	slice, err := h.redisClient.XRange(ctx,
		makeTickStreamKey(symbol),
		makeTickID(from.UnixMilli(), "0"),
		makeTickID(to.UnixMilli(), "999"),
	).Result()
	if err != nil {
		return errors.Wrap(err, "failed to read tick from redis stream")
	}

	var tickRange = make([]*entity.Tick, 0, len(slice))
	for _, item := range slice {
		tickRange = append(tickRange,
			mustUnmarshalTick(item.Values),
		)
	}

	return h.consumer.ConsumeTickRange(ctx,
		&entity.TickRange{
			Symbol:    symbol,
			FromMs:    from.UnixMilli(),
			ToMs:      to.UnixMilli(),
			TickSlice: tickRange,
		},
	)
}

func (h *TickRangeConsumer) RunBlocked(ctx context.Context, startTime time.Time, timeframe time.Duration, symbols []string) error {
	for timestamp := range goter.DelayedTimeIteratorWithContext(ctx, startTime, timeframe) {
		for _, symbol := range symbols {
			// Add 1 millisecond to avoid re-reading the last tick of the previous range
			var fromTime = timestamp.Add(-timeframe)
			var toTime = timestamp.Add(-time.Millisecond)
			if err := h.readNext(ctx, symbol, fromTime, toTime); err != nil {
				return errors.Wrap(err, "failed to consume tick range")
			}
		}
	}
	return nil
}
