package port

import (
	"context"

	"github.com/goregion/hexago/internal/entity"
)

type OHLCPublisher interface {
	PublishOHLC(context.Context, *entity.OHLC) error
}
