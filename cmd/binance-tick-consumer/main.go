package main

import (
	app_binance_tick_consumer "github.com/goregion/hexago/internal/app/binance-tick-consumer"
	"github.com/goregion/hexago/pkg/launcher"
	"github.com/goregion/hexago/pkg/log"
)

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
		WaitApplication(app_binance_tick_consumer.Launch).
		// log error if any
		LogIfError(logger)
}
