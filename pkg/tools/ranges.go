package tools

import (
	"context"
	"iter"
	"time"
)

func IteratorWithContext[ValueType any](ctx context.Context, iterator iter.Seq[ValueType]) iter.Seq[ValueType] {
	return func(yield func(ValueType) bool) {
		for v := range iterator {
			select {
			case <-ctx.Done():
				return
			default:
				if !yield(v) {
					return
				}
			}
		}
	}
}

func Iterator2WithContext[KeyType any, ValueType any](ctx context.Context, iterator iter.Seq2[KeyType, ValueType]) iter.Seq2[KeyType, ValueType] {
	return func(yield func(KeyType, ValueType) bool) {
		for k, v := range iterator {
			select {
			case <-ctx.Done():
				return
			default:
				if !yield(k, v) {
					return
				}
			}
		}
	}
}

type integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

func IntegerIteratorWithContext[Type integer](ctx context.Context) iter.Seq[Type] {
	return func(yield func(Type) bool) {
		var i Type
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if !yield(i) {
					return
				}
			}
			i++
		}
	}
}

func Uint64IteratorWithContext(ctx context.Context) iter.Seq[uint64] {
	return IntegerIteratorWithContext[uint64](ctx)
}

func DelayedTimeIteratorWithContext(ctx context.Context, startTime time.Time, duration time.Duration) iter.Seq[time.Time] {
	return func(yield func(time.Time) bool) {
		for {
			var toTime = startTime.Truncate(duration).Add(duration)
			var timer = time.NewTimer(time.Until(toTime))
			select {
			case <-ctx.Done():
				return
			case timestamp := <-timer.C:
				if !yield(timestamp) {
					return
				}
			}
			startTime = toTime
		}
	}
}
