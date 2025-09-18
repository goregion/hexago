package app

import (
	"context"

	"github.com/goregion/hexago/internal/entity"
	"github.com/goregion/hexago/internal/port"
)

type MessageProcessor struct {
	messagePublisher port.MessagePublisher
}

func NewMessageProcessor(messagePublisher port.MessagePublisher) *MessageProcessor {
	return &MessageProcessor{
		messagePublisher: messagePublisher,
	}
}

func (h *MessageProcessor) ConsumeMessage(ctx context.Context, message *entity.Message) error {
	var newMessage = &entity.Message{
		Value: message.Value + "_processed",
	}
	if err := h.messagePublisher.PublishMessage(ctx, newMessage); err != nil {
		return err
	}
	return nil
}
