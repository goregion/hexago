package port

import (
	"context"

	"github.com/goregion/hexago/internal/entity"
)

type LPTickConsumer interface {
	ConsumeLPTick(context.Context, *entity.LPTick) error
}
