package app

import (
	"context"
	"time"

	"github.com/goregion/hexago/internal/entity"
	"github.com/goregion/hexago/internal/port"
	"github.com/pkg/errors"
)

type TickProcessor struct {
	tickPublisher []port.TickPublisher
}

func NewTickProcessor(tickPublisher ...port.TickPublisher) *TickProcessor {
	return &TickProcessor{
		tickPublisher: tickPublisher,
	}
}

func (p *TickProcessor) ConsumeLPTick(ctx context.Context, lpTick *entity.LPTick) error {
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
	return nil
}
