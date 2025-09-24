package service_ohlc

import (
	"context"

	"github.com/goregion/hexago/internal/entity"
	"github.com/goregion/hexago/internal/port"
	"github.com/pkg/errors"
)

// OHLC publisher service
type OHLCPublisher struct {
	ohlcPublisher    []port.OHLCPublisher
	useBidOrAskPrice int
}

func NewOHLCPublisher(ohlcPublisher ...port.OHLCPublisher) *OHLCPublisher {
	return &OHLCPublisher{
		ohlcPublisher: ohlcPublisher,
	}
}

// ConsumeOHLC publishes the given OHLC data using all configured publishers
func (p *OHLCPublisher) ConsumeOHLC(ctx context.Context, ohlc *entity.OHLC) error {
	for _, p := range p.ohlcPublisher {
		if err := p.PublishOHLC(ctx, ohlc); err != nil {
			return errors.Wrap(err, "failed to publish OHLC")
		}
	}
	return nil
}
