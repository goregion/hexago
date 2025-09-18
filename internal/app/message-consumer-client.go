package app

import (
	"context"

	"github.com/pkg/errors"

	"github.com/goregion/hexago/internal/entity"
	"github.com/goregion/hexago/internal/port"
)

type MessageConsumerClient struct {
	messagePublisher port.MessagePublisher
}

func NewMessageConsumerClient(messagePublisher port.MessagePublisher) *MessageConsumerClient {
	return &MessageConsumerClient{
		messagePublisher: messagePublisher,
	}
}

func (h *MessageConsumerClient) HandleMessage(ctx context.Context, message *entity.Message) error {
	if err := h.messagePublisher.PublishMessage(ctx, message); err != nil {
		return errors.Wrap(err, "failed to publish message")
	}
	return nil
}
