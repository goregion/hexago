package goture

import (
	"context"

	"github.com/pkg/errors"
)

func recoverCancel(cancel context.CancelCauseFunc) {
	if r := recover(); r != nil {
		if err, ok := r.(error); ok {
			cancel(err)
			return
		}
		cancel(errors.Errorf("%v", r))
	}
}

func recoverCancelForParallel(ch chan<- error) {
	if r := recover(); r != nil {
		if err, ok := r.(error); ok {
			ch <- err
			return
		}
		ch <- errors.Errorf("%v", r)
	}
}

func makeErrorFromPanic(r interface{}) error {
	if err, ok := r.(error); ok {
		return err
	}
	return errors.Errorf("%v", r)
}
