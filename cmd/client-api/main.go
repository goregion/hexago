package main

import (
	"context"
	"log/slog"
	"os"

	service_client_api "github.com/goregion/hexago/internal/service/client-api"
	"github.com/goregion/hexago/pkg/log"
	"github.com/goregion/hexago/pkg/tools"
)

func main() {
	var ctx = tools.MakeGrExitContext(context.Background())

	var logger = log.NewLogger(
		slog.NewTextHandler(os.Stdout,
			&slog.HandlerOptions{
				Level: slog.LevelInfo,
			},
		),
	)

	logger.LogIfError(
		service_client_api.Run(ctx, logger),
	)
}
