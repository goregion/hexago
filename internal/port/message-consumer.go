package port

import (
	"context"

	"github.com/goregion/hexago/internal/entity"
)

type MessageConsumer interface {
	ConsumeMessage(context.Context, *entity.Message) error
}
