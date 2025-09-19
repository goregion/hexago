package app

import (
	"context"

	"github.com/goregion/hexago/internal/entity"
	"github.com/goregion/hexago/internal/port"
	"github.com/pkg/errors"
)

const (
	USE_BID_PRICE = iota
	USE_ASK_PRICE
)

type OHLCCreator struct {
	ohlcPublisher    []port.OHLCPublisher
	useBidOrAskPrice int
}

func NewOHLCCreator(useBidOrAskPrice int, ohlcPublisher ...port.OHLCPublisher) *OHLCCreator {
	return &OHLCCreator{
		ohlcPublisher:    ohlcPublisher,
		useBidOrAskPrice: useBidOrAskPrice,
	}
}

func (p *OHLCCreator) ConsumeTickRange(ctx context.Context, ticks *entity.TickRange) error {
	var ohlc = &entity.OHLC{
		Symbol:      ticks.Symbol,
		CloseTimeMs: ticks.ToMs,
	}

	for _, tick := range ticks.TickSlice {
		var price = tick.BestAskPrice
		if p.useBidOrAskPrice == USE_BID_PRICE {
			price = tick.BestBidPrice
		}

		if ohlc.Open == 0 {
			ohlc.Open = price
		}
		ohlc.Close = price
		if price > ohlc.High {
			ohlc.High = price
		}
		if price < ohlc.Low || ohlc.Low == 0 {
			ohlc.Low = price
		}
	}

	for _, p := range p.ohlcPublisher {
		if err := p.PublishOHLC(ctx, ohlc); err != nil {
			return errors.Wrap(err, "failed to publish OHLC")
		}
	}
	return nil
}
