package port

import (
	"context"

	"github.com/goregion/hexago/internal/entity"
)

type OHLCConsumer interface {
	ConsumeOHLC(context.Context, *entity.OHLC) error
}
