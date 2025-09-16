package main

import (
	"context"
	service_token_generator "feeder/internal/service/token-generator"
	"feeder/pkg/log"
	"log/slog"
	"os"
)

func main() {
	var logger = log.NewLogger(
		slog.NewTextHandler(os.Stdout,
			&slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		),
	)

	logger.LogIfError(
		service_token_generator.Run(context.Background(), logger),
	)
}
