package tools

import "context"

func RunAsyncBlocked(ctx context.Context, fn func(ctx context.Context) error) error {
	var localCtx, cancel = context.WithCancelCause(ctx)
	defer cancel(nil)
	go func() {
		cancel(
			fn(localCtx),
		)
	}()

	<-localCtx.Done()
	return localCtx.Err()
}
