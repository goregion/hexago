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
	startTime   time.Time
	timeframe   time.Duration
	symbols     []string
	consumer    port.TickRangeConsumer
}

func NewTickRangeConsumer(
	redisClient *redis.Client,
	startTime time.Time,
	timeframe time.Duration,
	symbols []string,
	consumer port.TickRangeConsumer,
) *TickRangeConsumer {
	return &TickRangeConsumer{
		redisClient: redisClient,
		startTime:   startTime,
		timeframe:   timeframe,
		symbols:     symbols,
		consumer:    consumer,
	}
}

// readNext reads the next range of ticks for a given symbol from Redis and passes them to the consumer
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

// Launch starts the consumer to read tick ranges for all symbols at intervals defined by the timeframe
func (h *TickRangeConsumer) Launch(ctx context.Context) error {
	for timestamp := range goter.DelayedTimeIteratorWithContext(ctx, h.startTime, h.timeframe) {
		for _, symbol := range h.symbols {
			// Add 1 millisecond to avoid re-reading the last tick of the previous range
			var fromTime = timestamp.Add(-h.timeframe)
			var toTime = timestamp.Add(-time.Millisecond)
			if err := h.readNext(ctx, symbol, fromTime, toTime); err != nil {
				return errors.Wrap(err, "failed to consume tick range")
			}
		}
	}
	return nil
}
