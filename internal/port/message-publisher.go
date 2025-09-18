package port

import (
	"context"

	"github.com/goregion/hexago/internal/entity"
)

type MessagePublisher interface {
	PublishMessage(context.Context, *entity.Message) error
}
