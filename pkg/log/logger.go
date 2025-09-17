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

func (l *Logger) Error(msg string, err ...error) {
	var erros []any
	for _, err := range err {
		if err != nil && !errors.Is(err, context.Canceled) {
			erros = append(erros, "error", err)
		}
	}
	if len(erros) > 0 {
		l.Logger.Error(msg, "error", err)
	}
}

func (l *Logger) LogIfError(err error, messages ...any) {
	if err != nil && !errors.Is(err, context.Canceled) {
		msg, args := formatMessage(err, messages...)
		l.Logger.Error(msg, args...)
	}
}

func formatMessage(err error, messages ...any) (string, []any) {
	if len(messages) > 0 {
		var msg = messages[0].(string)
		if len(messages) > 1 {
			var args = []any{err}
			return msg, append(args, messages[1:]...)
		}
		return msg, []any{"error", err}
	}
	return err.Error(), nil
}
