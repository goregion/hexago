package tools

import (
	"context"

	"github.com/pkg/errors"

	"github.com/goregion/hexago/pkg/log"
)

func RunService(ctx context.Context,
	logger *log.Logger,
	service func(context.Context, *log.Logger),
) {
	defer func() {
		if r := recover(); r != nil {
			var logError error
			if err, ok := r.(error); ok {
				logError = err
			} else {
				logError = errors.Errorf("%v", r)
			}
			logger.Error("panic", "error", logError)
		}
	}()

	service(
		ctx,
		logger,
	)
}

func RunServicesAsync(ctx context.Context,
	cancel context.CancelFunc,
	logger *log.Logger,
	services ...func(context.Context, *log.Logger),
) {
	for _, service := range services {
		go func() {
			defer cancel()
			RunService(
				ctx,
				logger,
				service,
			)
		}()
	}
}
