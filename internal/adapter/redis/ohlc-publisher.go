package adapter_redis

import (
	"context"

	"github.com/goregion/hexago/internal/entity"
	"github.com/goregion/hexago/pkg/redis"
	"github.com/pkg/errors"
)

type OHLCPublisher struct {
	redisClient   *redis.Client
	timeframeName string
}

func NewOHLCPublisher(redisClient *redis.Client, timeframeName string) *OHLCPublisher {
	return &OHLCPublisher{
		redisClient:   redisClient,
		timeframeName: timeframeName,
	}
}

// PublishOHLC publishes the given OHLC to the appropriate Redis stream
func (p *OHLCPublisher) PublishOHLC(ctx context.Context, ohlc *entity.OHLC) error {
	if err := p.redisClient.XAdd(ctx,
		&redis.XAddArgs{
			Stream: makeOHLCStreamName(p.timeframeName, ohlc.Symbol),
			Values: mustMarshalOHLC(ohlc),
		},
	).Err(); err != nil {
		return errors.Wrap(err, "failed to publish OHLC to redis stream")
	}
	return nil
}
