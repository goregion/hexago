package main

import (
	"context"

	app_backoffice_api "github.com/goregion/hexago/internal/app/backoffice-api"
	app_binance_tick_consumer "github.com/goregion/hexago/internal/app/binance-tick-consumer"
	app_ohlc_generator "github.com/goregion/hexago/internal/app/ohlc-generator"
	"github.com/goregion/hexago/pkg/log"
	"github.com/goregion/hexago/pkg/tools"
)

func RunAsync(ctx context.Context, logger *log.Logger, cancel context.CancelFunc, service func(context.Context, *log.Logger)) {
	go func() {
		defer cancel()
		service(
			ctx,
			logger,
		)
	}()
}

func main() {
	// context with graceful exit on SIGINT, SIGTERM
	var ctx, cancel = tools.MakeGrExitWithCancelContext(
		context.Background(),
	)
	defer cancel()

	// logger
	var logger = log.NewLogger(
		log.NewTextStdOutHandler(),
	)

	// apps to run in parallel
	tools.RunAppAsync(ctx, cancel, logger,
		app_binance_tick_consumer.RunBlocked,
		app_ohlc_generator.RunBlocked,
		app_backoffice_api.RunBlocked,
	)

	// wait until context is done
	<-ctx.Done()
}
