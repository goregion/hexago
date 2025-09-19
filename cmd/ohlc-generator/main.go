package main

import (
	"context"

	service_ohlc_generator "github.com/goregion/hexago/internal/service/ohlc-generator"
	"github.com/goregion/hexago/pkg/log"
	"github.com/goregion/hexago/pkg/tools"
)

func main() {
	tools.RunService(
		// context with graceful exit on SIGINT, SIGTERM
		tools.MakeGrExitWithContext(
			context.Background(),
		),
		// logger
		log.NewLogger(
			log.NewTextStdOutHandler(),
		),
		// service to run
		service_ohlc_generator.RunBlocked,
	)
}
