package main

import (
	"context"

	app_ohlc_generator "github.com/goregion/hexago/internal/app/ohlc-generator"
	"github.com/goregion/hexago/pkg/log"
	"github.com/goregion/hexago/pkg/tools"
)

func main() {
	tools.RunApp(
		// context with graceful exit on SIGINT, SIGTERM
		tools.MakeGrExitWithContext(
			context.Background(),
		),
		// logger
		log.NewLogger(
			log.NewTextStdOutHandler(),
		),
		// app to run
		app_ohlc_generator.RunBlocked,
	)
}
