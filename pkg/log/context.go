package log

import (
	"context"

	"github.com/pkg/errors"
)

type loggerContextKey string

var LoggerContextKey = loggerContextKey("logger")

func ContextWithLogger(ctx context.Context, logger *Logger) context.Context {
	return context.WithValue(ctx, LoggerContextKey, logger)
}

func LoggerFromContext(ctx context.Context) (*Logger, error) {
	if logger, ok := ctx.Value(LoggerContextKey).(*Logger); ok {
		return logger, nil
	}
	return nil, errors.New("logger not found in context")
}
