// Package goture provides a Future pattern implementation for Go.
// It allows executing tasks asynchronously and waiting for their completion.
package goture

import (
	"context"
)

// SuccessResult represents successful completion of a task
type SuccessResult struct{}

func (s SuccessResult) Error() string {
	return ""
}

// Task represents a function that can be executed asynchronously
type Task func(ctx context.Context) error

// Goture represents a future that will complete when the associated task finishes
type Goture struct {
	ctx context.Context
}

// Wait blocks until the task completes and returns any error that occurred
func (f Goture) Wait() error {
	<-f.ctx.Done()
	cause := context.Cause(f.ctx)
	if _, ok := cause.(SuccessResult); ok {
		return nil
	}
	return cause
}

// NewGoture creates a new Goture that executes the given task asynchronously.
// The task will start executing immediately in a separate goroutine.
func NewGoture(ctx context.Context, fn Task) Goture {
	var localCtx, cancel = context.WithCancelCause(ctx)
	go func() {
		defer recoverCancel(cancel)
		if err := fn(localCtx); err != nil {
			cancel(err)
		} else {
			cancel(SuccessResult{})
		}
	}()
	return Goture{ctx: localCtx}
}

// NewParallelGoture creates a new Goture that executes all given tasks in parallel.
// It waits for all tasks to complete and returns an error if any task fails.
// If multiple tasks fail, it returns the first error encountered.
func NewParallelGoture(parentCtx context.Context, tasks ...Task) Goture {
	if len(tasks) == 0 {
		// Return completed future for empty task list
		localCtx, cancel := context.WithCancelCause(parentCtx)
		cancel(SuccessResult{})
		return Goture{ctx: localCtx}
	}

	var localCtx, cancel = context.WithCancelCause(parentCtx)

	// Use sync mechanism to wait for all tasks
	completed := make(chan error, len(tasks))

	for _, fn := range tasks {
		go func(task Task) {
			defer recoverCancelForParallel(completed)
			completed <- task(localCtx)
		}(fn)
	}

	// Goroutine to wait for all tasks completion
	go func() {
		var firstError error
		for i := 0; i < len(tasks); i++ {
			if err := <-completed; err != nil && firstError == nil {
				firstError = err
			}
		}
		if firstError != nil {
			cancel(firstError)
		} else {
			cancel(SuccessResult{})
		}
	}()

	return Goture{ctx: localCtx}
}
