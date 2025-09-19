package tools

import "context"

func RunAsyncAndWait(ctx context.Context, fn func() error) error {
	var localCtx, cancel = context.WithCancelCause(ctx)
	defer cancel(nil)
	go func() {
		cancel(
			fn(),
		)
	}()

	<-localCtx.Done()
	return localCtx.Err()
}
