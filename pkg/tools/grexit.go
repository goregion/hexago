package tools

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// MakeGrExitContext returns a context that is canceled on SIGINT or SIGTERM.
func MakeGrExitContext(ctx context.Context) context.Context {
	var interrupt = make(chan os.Signal, 1)
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		defer cancel()

		signal.Notify(interrupt,
			syscall.SIGINT,
			syscall.SIGTERM,
		)

		select {
		case <-interrupt: // caught SIGINT or SIGTERM
			cancel()
		case <-ctx.Done(): // context was canceled elsewhere
		}
	}()
	return ctx
}
