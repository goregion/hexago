package tools

import (
	"context"
	"iter"
)

type integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

func IntegerIterator[Type integer]() iter.Seq[Type] {
	return func(yield func(Type) bool) {
		var i Type
		for {
			if !yield(i) {
				return
			}
			i++
		}
	}
}

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

func IteratorInt64WithContext(ctx context.Context) iter.Seq[int64] {
	return IteratorWithContext(ctx, IntegerIterator[int64]())
}
