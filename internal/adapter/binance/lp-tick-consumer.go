package adapter_binance

import (
	"context"
	"strconv"

	"github.com/adshao/go-binance/v2"
	"github.com/goregion/hexago/internal/entity"
	"github.com/goregion/hexago/internal/port"
	"github.com/pkg/errors"
)

type LPTickConsumer struct {
	symbols               []string
	lpTickConsumer        port.LPTickConsumer
	tickEventErrorHandler func(err error)
}

func NewLPTickConsumer(symbols []string, lpTickConsumer port.LPTickConsumer, tickEventErrorHandler func(err error)) *LPTickConsumer {
	return &LPTickConsumer{
		symbols:               symbols,
		lpTickConsumer:        lpTickConsumer,
		tickEventErrorHandler: tickEventErrorHandler,
	}
}

func (c *LPTickConsumer) handleTickEvent(event *binance.WsBookTickerEvent) {
	var tick = &entity.LPTick{
		Symbol: event.Symbol,
	}

	var err error
	tick.BestBidPrice, err = strconv.ParseFloat(event.BestBidPrice, 64)
	if err != nil {
		c.tickEventErrorHandler(errors.Wrap(err, "failed to parse tick price"))
		return
	}
	tick.BestAskPrice, err = strconv.ParseFloat(event.BestAskPrice, 64)
	if err != nil {
		c.tickEventErrorHandler(errors.Wrap(err, "failed to parse tick price"))
		return
	}

	if err := c.lpTickConsumer.ConsumeLPTick(context.Background(), tick); err != nil {
		c.tickEventErrorHandler(errors.Wrap(err, "failed to consume tick"))
	}
}

// Launch starts the Binance LP tick consumer to listen for tick events for the configured symbols
func (c *LPTickConsumer) Launch(ctx context.Context) error {
	for {
		doneChan, stopChan, err := binance.WsCombinedBookTickerServe(c.symbols,
			c.handleTickEvent,
			c.tickEventErrorHandler,
		)
		if err != nil {
			c.tickEventErrorHandler(errors.Wrap(err, "failed to start binance tick consumer"))
		}

		select {
		case <-ctx.Done():
			stopChan <- struct{}{}
			return ctx.Err()
		case <-doneChan:
			// reconnect when connection is lost
		}
	}
}
