package log

import (
	"context"
	"errors"
	"log/slog"
)

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

func (l *Logger) StartService(serviceName string) (*Logger, func()) {
	var logger = l.Logger.With(
		"service",
		serviceName,
	)

	logger.Info("start")
	var result = &Logger{
		Logger: logger,
	}

	return result,
		func() {
			logger.Info("stop")
		}
}

func (l *Logger) LogIfError(err error, messages ...any) {
	if err != nil && !errors.Is(err, context.Canceled) {
		var msg, args = formatMessage(err, messages...)
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
