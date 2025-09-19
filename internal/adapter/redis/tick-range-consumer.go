package adapter_redis

import (
	"context"
	"time"

	"github.com/goregion/hexago/internal/entity"
	"github.com/goregion/hexago/internal/port"
	"github.com/goregion/hexago/pkg/redis"
	"github.com/goregion/hexago/pkg/tools"
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
		makeTickID(from.UnixMilli()),
		makeTickID(to.UnixMilli()),
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

	return h.consumer.ConsumeTickRange(ctx, tickRange)
}

func (h *TickRangeConsumer) RunBlocked(ctx context.Context, startTime time.Time, timeframeName time.Duration, symbols []string) error {
	for timestamp := range tools.DelayedTimeIteratorWithContext(ctx, startTime, timeframeName) {
		for _, symbol := range symbols {
			if err := h.readNext(ctx, symbol, timestamp.Add(-timeframeName), timestamp); err != nil {
				return errors.Wrap(err, "failed to consume tick range")
			}
		}
	}
	return nil
}
