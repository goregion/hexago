package log

import (
	"context"
	"errors"
	"log/slog"
)

type loggerContextKeyType string

var loggerContextKey loggerContextKeyType = "logger"

type Logger struct {
	*slog.Logger
}

func NewLogger(
	handlers ...slog.Handler,
) *Logger {
	return &Logger{
		Logger: slog.New(
			&multiHandler{
				handlers: handlers,
			},
		),
	}
}

func (l *Logger) StartService(ctx context.Context, serviceName string) (*Logger, context.Context, func()) {
	var logger = l.Logger.With(
		"service",
		serviceName,
	)

	logger.Info("start service")
	var result = &Logger{
		Logger: logger,
	}

	return result,
		context.WithValue(ctx, loggerContextKey, result),
		func() {
			logger.Info("stop service")
		}
}

func GetLoggerContextKey() string {
	return string(loggerContextKey)
}

func SetLoggerContextKey(key string) {
	loggerContextKey = loggerContextKeyType(key)
}

func GetLogger(ctx context.Context) *Logger {
	var log = ctx.Value(loggerContextKey)
	if log == nil {
		panic("fatal error, logger not found in context")
	}
	return log.(*Logger)
}

func (l *Logger) Error(msg string, err error) {
	if err != nil && !errors.Is(err, context.Canceled) {
		l.Logger.Error(msg, "error", err)
	}
}

func (l *Logger) LogIfError(msg string, err error) {
	if err != nil && !errors.Is(err, context.Canceled) {
		l.Logger.Error(msg, "error", err)
	}
}
