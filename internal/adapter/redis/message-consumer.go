package adapter_redis

import (
	"context"
	"strconv"

	"github.com/goregion/hexago/internal/entity"
	"github.com/goregion/hexago/internal/port"
	"github.com/goregion/hexago/pkg/redis"
	"github.com/pkg/errors"
)

type MessageConsumer struct {
	redisSub       *redis.PubSub
	messageHandler port.MessageConsumer
}

func NewMessageConsumer(ctx context.Context, redisClient *redis.Client, messageHandler port.MessageConsumer) *MessageConsumer {
	return &MessageConsumer{
		redisSub: redisClient.Subscribe(ctx,
			messagesChannel,
		),
		messageHandler: messageHandler,
	}
}

func (h *MessageConsumer) ReadMessage(ctx context.Context) error {
	msg, err := h.redisSub.ReceiveMessage(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to consume message from redis")
	}
	var key, errParse = strconv.Atoi(msg.PayloadSlice[0])
	if errParse != nil {
		return errors.Wrap(errParse, "failed to parse message key")
	}
	return h.messageHandler.ConsumeMessage(ctx, &entity.Message{
		Key:   key,
		Value: msg.PayloadSlice[0],
	})
}
