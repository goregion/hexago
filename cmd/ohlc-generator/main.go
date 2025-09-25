package main

import (
	"context"

	"github.com/goregion/goture"
	"github.com/goregion/grexit"
	app_ohlc_generator "github.com/goregion/hexago/internal/app/ohlc-generator"
	"github.com/goregion/hexago/pkg/log"
)

func main() {
	var logger = log.NewLogger(
		log.NewTextStdOutHandler(),
	)

	// create context
	var ctx = context.Background()
	// inject graceful exit on SIGINT, SIGTERM
	ctx = grexit.WithGrexitContext(ctx)
	// inject logger into context
	ctx = log.WithLoggerContext(ctx, logger)

	// run application and get future
	var future = goture.NewGoture(ctx, app_ohlc_generator.Launch)
	// wait for application to finish
	var err = future.Wait()

	// log error if any
	logger.LogIfError(err)
}
