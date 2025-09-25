package main

import (
	"context"

	"github.com/goregion/goture"
	"github.com/goregion/grexit"
	app_backoffice_api "github.com/goregion/hexago/internal/app/backoffice-api"
	"github.com/goregion/hexago/pkg/log"
)

func main() {
	var logger = log.NewLogger(
		log.NewTextStdOutHandler(),
	)

	logger.LogIfError(
		goture.NewGoture(
			log.WithLoggerContext(
				grexit.WithGrexitContext(
					context.Background(),
				),
				logger,
			),
			app_backoffice_api.Launch,
		).Wait(),
	)
}
