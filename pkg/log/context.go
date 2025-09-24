package log

import (
	"context"

	"github.com/pkg/errors"
)

const defaultLoggerContextKey = "logger"

type loggerContextKey string

var LoggerContextKey = loggerContextKey(defaultLoggerContextKey)

func WithLoggerContext(ctx context.Context, logger *Logger) context.Context {
	return context.WithValue(ctx, LoggerContextKey, logger)
}

func MustGetLoggerFromContext(ctx context.Context) *Logger {
	logger, err := GetLoggerFromContext(ctx)
	if err != nil {
		panic(err)
	}
	return logger
}

func GetLoggerFromContext(ctx context.Context) (*Logger, error) {
	if logger, ok := ctx.Value(LoggerContextKey).(*Logger); ok {
		return logger, nil
	}
	return nil, errors.Errorf("logger not found in context with key '%s'", LoggerContextKey)
}

func SetLoggerContextKey(key string) {
	LoggerContextKey = loggerContextKey(key)
}
