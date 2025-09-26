package service_ohlc

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

// OHLC creator from ticks
type OHLCCreator struct {
	ohlcPublisher                port.OHLCPublisher
	ohlcRepository               port.OHLCRepository
	repositoryTransactionManager port.TransactionManager
	useBidOrAskPrice             int
}

func NewOHLCCreator(useBidOrAskPrice int, ohlcPublisher port.OHLCPublisher, repositoryTransactionManager port.TransactionManager, ohlcRepository port.OHLCRepository) *OHLCCreator {
	return &OHLCCreator{
		ohlcPublisher:                ohlcPublisher,
		ohlcRepository:               ohlcRepository,
		repositoryTransactionManager: repositoryTransactionManager,
		useBidOrAskPrice:             useBidOrAskPrice,
	}
}

// ConsumeTickRange processes a range of ticks to create an OHLC and publishes it
func (p *OHLCCreator) ConsumeTickRange(ctx context.Context, ticks *entity.TickRange) error {
	var ohlc = &entity.OHLC{
		Symbol:      ticks.Symbol,
		TimestampMs: ticks.ToMs,
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

	if err := p.ohlcPublisher.PublishOHLC(ctx, ohlc); err != nil {
		return errors.Wrap(err, "failed to publish OHLC")
	}

	transactionCtx, commit, rollback, err := p.repositoryTransactionManager.WithTx(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to create transaction")
	}
	defer rollback()

	if err := p.ohlcRepository.StoreOHLC(transactionCtx, ohlc); err != nil {
		return errors.Wrap(err, "failed to store OHLC")
	}

	commit()
	return nil
}
