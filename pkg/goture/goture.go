package goture

import (
	"context"

	"github.com/pkg/errors"
)

type TaskWithResult[ResultType any] func(ctx context.Context) (ResultType, error)
type Task func(ctx context.Context) error

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

type GotureWithResult[ResultType any] struct {
	ctx context.Context
}

func (f GotureWithResult[ResultType]) Wait() (ResultType, error) {
	<-f.ctx.Done()
	if f.ctx.Err() != nil {
		var res = f.ctx.Err().(result[ResultType])
		return res.value, res.err
	}
	var zero ResultType
	return zero, nil
}

type Goture struct {
	ctx context.Context
}

func (f Goture) Wait() error {
	<-f.ctx.Done()
	return f.ctx.Err()
}

func recoverCancel(cancel context.CancelCauseFunc) {
	if r := recover(); r != nil {
		if err, ok := r.(error); ok {
			cancel(err)
			return
		}
		cancel(errors.Errorf("%v", r))
	}
}

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

func NewGoture(ctx context.Context, fn Task) Goture {
	var localCtx, cancel = context.WithCancelCause(ctx)
	go func() {
		defer recoverCancel(cancel)
		cancel(
			fn(localCtx),
		)
	}()
	return Goture{ctx: localCtx}
}

func NewParallelGoture(parentCtx context.Context, tasks ...Task) Goture {
	if len(tasks) == 0 {
		return Goture{ctx: parentCtx}
	}

	var localCtx, cancel = context.WithCancelCause(parentCtx)
	for _, fn := range tasks {
		go func(task Task) {
			defer recoverCancel(cancel)
			cancel(
				task(localCtx),
			)
		}(fn)
	}
	return Goture{ctx: localCtx}
}
