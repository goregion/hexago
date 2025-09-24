package adapter_redis

import (
	"context"

	"github.com/goregion/hexago/internal/port"
	"github.com/goregion/hexago/pkg/goter"
	"github.com/goregion/hexago/pkg/redis"
	"github.com/pkg/errors"
)

type OHLCConsumer struct {
	redisClient   *redis.Client
	timeframeName string
	symbols       []string
	consumer      port.OHLCConsumer
}

func NewOHLCConsumer(redisClient *redis.Client, consumer port.OHLCConsumer, timeframeName string, symbols []string) *OHLCConsumer {
	return &OHLCConsumer{
		redisClient:   redisClient,
		timeframeName: timeframeName,
		consumer:      consumer,
		symbols:       symbols,
	}
}

func (h *OHLCConsumer) readNext(ctx context.Context, symbol string) error {
	streams, err := h.redisClient.XRead(ctx,
		&redis.XReadArgs{
			Streams: []string{makeOHLCStreamName(h.timeframeName, symbol), "$"},
			Count:   1,
			Block:   0,
		},
	).Result()
	if err != nil {
		return errors.Wrap(err, "failed to read OHLC from redis stream")
	}

	if len(streams) > 0 && len(streams[0].Messages) > 0 {
		if err := h.consumer.ConsumeOHLC(ctx, mustUnmarshalOHLC(streams[0].Messages[0].Values)); err != nil {
			return errors.Wrap(err, "failed to consume OHLC message")
		}
	}
	return nil
}

func (h *OHLCConsumer) RunBlocked(ctx context.Context) error {
	for range goter.Uint64IteratorWithContext(ctx) {
		for _, symbol := range h.symbols {
			if err := h.readNext(ctx, symbol); err != nil {
				return errors.Wrap(err, "failed to consume OHLC")
			}
		}
	}
	return nil
}
