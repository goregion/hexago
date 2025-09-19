package app

import (
	"context"
	"sync"
	"time"

	"github.com/goregion/hexago/internal/entity"
	"github.com/goregion/hexago/internal/port"
	"github.com/pkg/errors"
)

type prices struct {
	BidPrice float64
	AskPrice float64
}

type TickProcessor struct {
	tickPublisher []port.TickPublisher

	quantityMap sync.Map // symbol, prices
}

func NewTickProcessor(tickPublisher ...port.TickPublisher) *TickProcessor {
	return &TickProcessor{
		tickPublisher: tickPublisher,
	}
}

func (p *TickProcessor) ConsumeLPTick(ctx context.Context, lpTick *entity.LPTick) error {
	var totalPrices = prices{
		BidPrice: lpTick.BestBidPrice,
		AskPrice: lpTick.BestAskPrice,
	}

	// LoadOrStore returns (actual_value, loaded)
	// loaded = true if the value already existed
	// loaded = false if the value was stored for the first time
	actual, loaded := p.quantityMap.LoadOrStore(lpTick.Symbol, totalPrices)

	// Publish if:
	// 1. Value didn't exist (loaded = false) - first time for this symbol
	// 2. Value changed (actual != totalPrices)
	actualPrices := actual.(prices)
	if !loaded || actualPrices.AskPrice != totalPrices.AskPrice || actualPrices.BidPrice != totalPrices.BidPrice {
		if loaded && actual != totalPrices {
			p.quantityMap.Store(lpTick.Symbol, totalPrices)
		}

		var tick = &entity.Tick{
			Symbol:       lpTick.Symbol,
			BestBidPrice: lpTick.BestBidPrice,
			BestAskPrice: lpTick.BestAskPrice,
			TimestampMs:  time.Now().UnixMilli(),
		}

		for _, p := range p.tickPublisher {
			if err := p.PublishTick(ctx, tick); err != nil {
				return errors.Wrap(err, "failed to publish tick")
			}
		}
	}
	return nil
}
