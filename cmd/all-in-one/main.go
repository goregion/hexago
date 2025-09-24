package main

import (
	"context"

	app_backoffice_api "github.com/goregion/hexago/internal/app/backoffice-api"
	app_binance_tick_consumer "github.com/goregion/hexago/internal/app/binance-tick-consumer"
	app_ohlc_generator "github.com/goregion/hexago/internal/app/ohlc-generator"
	"github.com/goregion/hexago/pkg/launcher"
	"github.com/goregion/hexago/pkg/log"
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
	// configure logger
	logger := log.NewLogger(
		log.NewTextStdOutHandler(),
	)

	launcher.NewAppLauncher().
		// inject context that is canceled on SIGINT, SIGTERM
		WithGrexitContext().
		// inject logger into context
		WithLoggerContext(logger).
		// Run application, wait for it to finish
		WaitApplications(
			app_binance_tick_consumer.Launch,
			app_ohlc_generator.Launch,
			app_backoffice_api.Launch,
		).
		// log error if any
		LogIfError(logger)
}
