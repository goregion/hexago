package port

import (
	"context"

	"github.com/goregion/hexago/internal/entity"
)

type TickRangeConsumer interface {
	ConsumeTickRange(context.Context, *entity.TickRange) error
}
