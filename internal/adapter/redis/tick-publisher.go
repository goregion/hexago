package adapter_redis

import (
	"context"

	"github.com/goregion/hexago/internal/entity"
	"github.com/goregion/hexago/pkg/redis"
	"github.com/pkg/errors"
)

type TickPublisher struct {
	redisClient *redis.Client
}

func NewTickPublisher(redisClient *redis.Client) *TickPublisher {
	return &TickPublisher{
		redisClient: redisClient,
	}
}

func (p *TickPublisher) PublishTick(ctx context.Context, tick *entity.Tick) error {
	if err := p.redisClient.XAdd(ctx,
		&redis.XAddArgs{
			ID:     makeTickID(tick.TimestampMs, "*"),
			Stream: makeTickStreamKey(tick.Symbol),
			Values: mustMarshalTick(tick),
		},
	).Err(); err != nil {
		return errors.Wrap(err, "failed to publish tick to redis stream")
	}
	return nil
}
