package main

import (
	"context"

	service_backoffice_api "github.com/goregion/hexago/internal/service/backoffice-api"
	service_binance_tick_consumer "github.com/goregion/hexago/internal/service/binance-tick-consumer"
	service_ohlc_generator "github.com/goregion/hexago/internal/service/ohlc-generator"
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

	// services to run in parallel
	tools.RunServicesAsync(ctx, cancel, logger,
		service_binance_tick_consumer.RunBlocked,
		service_ohlc_generator.RunBlocked,
		service_backoffice_api.RunBlocked,
	)

	// wait until context is done
	<-ctx.Done()
}
