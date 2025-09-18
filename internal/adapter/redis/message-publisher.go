package adapter_redis

import (
	"context"
	"strconv"

	"github.com/goregion/hexago/internal/entity"
	"github.com/goregion/hexago/pkg/redis"
	"github.com/pkg/errors"
)

type MessagePublisher struct {
	redisClient *redis.Client
}

func NewMessagePublisher(redisClient *redis.Client) *MessagePublisher {
	return &MessagePublisher{
		redisClient: redisClient,
	}
}

func (p *MessagePublisher) PublishMessage(ctx context.Context, message *entity.Message) error {
	if err := p.redisClient.Publish(ctx,
		messagesChannel,
		[]string{
			strconv.Itoa(message.Key),
			message.Value,
		},
	).Err(); err != nil {
		return errors.Wrap(err, "failed to publish message to redis")
	}
	return nil
}
