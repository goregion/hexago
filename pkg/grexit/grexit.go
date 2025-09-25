package grexit

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// DefaultShutdownTimeout is the default timeout for graceful shutdown.
const DefaultShutdownTimeout = 30 * time.Second

// WithGrexitContext returns a context that is canceled on SIGINT or SIGTERM.
func WithGrexitContext(ctx context.Context) context.Context {
	ctx, _ = WithGrexitCancelContext(ctx)
	return ctx
}

// WithGrexitCancelContext returns a context and cancel function that is canceled on SIGINT or SIGTERM.
func WithGrexitCancelContext(ctx context.Context) (context.Context, context.CancelFunc) {
	interrupt := make(chan os.Signal, 1)
	ctx, cancel := context.WithCancel(ctx)

	signal.Notify(interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	go func() {
		defer signal.Stop(interrupt) // Clean up signal notification
		defer close(interrupt)       // Close the channel to prevent goroutine leak

		select {
		case <-interrupt: // caught SIGINT or SIGTERM
			cancel()
		case <-ctx.Done(): // context was canceled elsewhere
		}
	}()

	return ctx, cancel
}

// WithGrexitTimeout returns a context that is canceled on SIGINT or SIGTERM with a default timeout.
func WithGrexitTimeout(ctx context.Context) context.Context {
	return WithGrexitTimeoutDuration(ctx, DefaultShutdownTimeout)
}

// WithGrexitTimeoutDuration returns a context that is canceled on SIGINT or SIGTERM with a specified timeout.
func WithGrexitTimeoutDuration(ctx context.Context, timeout time.Duration) context.Context {
	ctx, _ = WithGrexitTimeoutCancelContext(ctx, timeout)
	return ctx
}

// WithGrexitTimeoutCancel returns a context and cancel function with default timeout.
func WithGrexitTimeoutCancel(ctx context.Context) (context.Context, context.CancelFunc) {
	return WithGrexitTimeoutCancelContext(ctx, DefaultShutdownTimeout)
}

// WithGrexitTimeoutCancelContext returns a context and cancel function that is canceled on SIGINT or SIGTERM
// with a specified timeout. If the timeout expires before graceful shutdown completes, the context is canceled.
func WithGrexitTimeoutCancelContext(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	interrupt := make(chan os.Signal, 1)
	ctx, cancel := context.WithCancel(ctx)

	signal.Notify(interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	go func() {
		defer signal.Stop(interrupt)
		defer close(interrupt)

		var timeoutTimer *time.Timer
		var timeoutChan <-chan time.Time

		select {
		case <-interrupt: // caught SIGINT or SIGTERM
			// Start the timeout timer for graceful shutdown
			timeoutTimer = time.NewTimer(timeout)
			timeoutChan = timeoutTimer.C

			// Cancel the context to signal graceful shutdown
			cancel()

			// Wait for either completion or timeout
			select {
			case <-ctx.Done():
				// Graceful shutdown completed, stop the timer
				if timeoutTimer != nil && !timeoutTimer.Stop() {
					<-timeoutTimer.C // drain the channel if timer already fired
				}
			case <-timeoutChan:
				// Timeout expired, force shutdown
				cancel() // ensure context is canceled
			}

		case <-ctx.Done(): // context was canceled elsewhere
		}
	}()

	return ctx, cancel
}
