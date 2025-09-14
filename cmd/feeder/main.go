package main

import (
	"context"
	"feeder/internal/service"
	"feeder/pkg/log"
	"feeder/pkg/tools"
	"log/slog"
	"os"
)

func main() {
	ctx := tools.MakeGrexitWithContext(context.Background())

	var logger = log.NewLogger(
		slog.NewTextHandler(os.Stdout,
			&slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		),
	)

	service.Feeder(ctx, logger)

	return
}
