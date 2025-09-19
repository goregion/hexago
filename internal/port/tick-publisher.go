package port

import (
	"context"

	"github.com/goregion/hexago/internal/entity"
)

type TickPublisher interface {
	PublishTick(context.Context, *entity.Tick) error
}
