package tools

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func MakeGrexitWithContext(ctx context.Context) context.Context {
	var interrupt = make(chan os.Signal, 1)
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		defer cancel()

		signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
		select {
		case <-interrupt:
			cancel()
		case <-ctx.Done():
		}
	}()
	return ctx
}
