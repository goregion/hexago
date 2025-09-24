package tools

import "context"

type AbstractAdapter interface {
	RunBlocked(ctx context.Context) error
}
