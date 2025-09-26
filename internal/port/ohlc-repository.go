package port

import (
	"context"

	"github.com/goregion/hexago/internal/entity"
)

type OHLCRepository interface {
	StoreOHLC(context.Context, *entity.OHLC) error
}
