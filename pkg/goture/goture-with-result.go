package goture

import (
	"context"
)

// TaskWithResult represents a function that can be executed asynchronously and returns a result
type TaskWithResult[ResultType any] func(ctx context.Context) (ResultType, error)

// result wraps the return value and error from a TaskWithResult
type result[ResultType any] struct {
	value ResultType
	err   error
}

func makeResult[ResultType any](value ResultType, err error) result[ResultType] {
	return result[ResultType]{value: value, err: err}
}

func (r result[ResultType]) Error() string {
	if r.err == nil {
		return ""
	}
	return r.err.Error()
}

// Implement error interface for better error handling
func (r result[ResultType]) Unwrap() error {
	return r.err
}

// GotureWithResult represents a future that will complete with a result when the associated task finishes
type GotureWithResult[ResultType any] struct {
	ctx context.Context
}

// Wait blocks until the task completes and returns the result and any error that occurred
func (f GotureWithResult[ResultType]) Wait() (ResultType, error) {
	<-f.ctx.Done()
	if cause := context.Cause(f.ctx); cause != nil {
		if res, ok := cause.(result[ResultType]); ok {
			return res.value, res.err
		}
		// Handle case where cause is not result[ResultType] (e.g., context cancellation)
		var zero ResultType
		return zero, cause
	}
	var zero ResultType
	return zero, nil
}

// NewGotureWithResult creates a new GotureWithResult that executes the given task asynchronously.
// The task will start executing immediately in a separate goroutine.
// Use this when you need to get a result value from the task.
func NewGotureWithResult[ResultType any](ctx context.Context, fn TaskWithResult[ResultType]) GotureWithResult[ResultType] {
	var localCtx, cancel = context.WithCancelCause(ctx)
	go func() {
		defer recoverCancel(cancel)
		cancel(
			makeResult(
				fn(localCtx),
			),
		)
	}()
	return GotureWithResult[ResultType]{ctx: localCtx}
}

// NewParallelWithResult creates a new GotureWithResult that executes all given tasks in parallel.
// It waits for all tasks to complete and returns a slice of results in the same order as the input tasks.
// If any task fails, the error from the first failing task is returned along with partial results.
// Results from tasks that completed successfully before the error will be included in the result slice.
func NewParallelWithResult[ResultType any](parentCtx context.Context, tasks ...TaskWithResult[ResultType]) GotureWithResult[[]ResultType] {
	if len(tasks) == 0 {
		// Return completed future for empty task list
		localCtx, cancel := context.WithCancelCause(parentCtx)
		cancel(makeResult([]ResultType{}, nil))
		return GotureWithResult[[]ResultType]{ctx: localCtx}
	}

	var localCtx, cancel = context.WithCancelCause(parentCtx)

	// Use sync mechanism to collect results from all tasks
	type taskResult struct {
		index  int
		result ResultType
		err    error
	}

	completed := make(chan taskResult, len(tasks))

	for i, fn := range tasks {
		go func(index int, task TaskWithResult[ResultType]) {
			defer func() {
				if r := recover(); r != nil {
					var zero ResultType
					if err, ok := r.(error); ok {
						completed <- taskResult{index: index, result: zero, err: err}
					} else {
						completed <- taskResult{index: index, result: zero, err: makeErrorFromPanic(r)}
					}
				}
			}()

			result, err := task(localCtx)
			completed <- taskResult{index: index, result: result, err: err}
		}(i, fn)
	}

	// Goroutine to wait for all tasks completion
	go func() {
		results := make([]ResultType, len(tasks))
		var firstError error

		for i := 0; i < len(tasks); i++ {
			taskRes := <-completed
			results[taskRes.index] = taskRes.result

			if taskRes.err != nil && firstError == nil {
				firstError = taskRes.err
			}
		}

		cancel(makeResult(results, firstError))
	}()

	return GotureWithResult[[]ResultType]{ctx: localCtx}
}
