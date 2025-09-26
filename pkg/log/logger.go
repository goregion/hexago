// Package log provides a structured logging wrapper around the standard log/slog package.
// It offers enhanced functionality including multi-handler support, context integration,
// and convenient utility methods for common logging patterns.
package log

import (
	"context"
	"errors"
	"log/slog"
)

// Logger wraps the standard slog.Logger with additional functionality
// and provides convenient methods for structured logging.
type Logger struct {
	*slog.Logger
}

// NewLogger creates a new Logger instance with the specified handlers.
// If no handlers are provided, it defaults to using JSON stdout handler.
// Multiple handlers can be used simultaneously through the internal multiHandler.
func NewLogger(handlers ...slog.Handler) *Logger {
	if len(handlers) == 0 {
		// По умолчанию используем JSON handler для stdout
		handlers = append(handlers, NewJsonStdOutHandler())
	}

	return &Logger{
		Logger: slog.New(
			&multiHandler{
				handlers: handlers,
			},
		),
	}
}

// StartService creates a new logger instance for a specific service and returns
// a cleanup function. It automatically logs service start and provides a function
// that logs service stop when called. This is useful for tracking service lifecycle.
//
// Example:
//
//	serviceLogger, stop := logger.StartService("user-service")
//	defer stop() // Will log "stop" message
//	serviceLogger.Info("Processing request")
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

// LogIfError logs an error message only if the provided error is not nil
// and is not a context.Canceled error. This helps avoid logging expected
// cancellation errors while still capturing actual problems.
//
// The first message parameter should be a format string, followed by any
// additional arguments. The error will be automatically included in the log.
func (l *Logger) LogIfError(err error, messages ...any) {
	if err != nil && !errors.Is(err, context.Canceled) {
		var msg, args = formatMessage(err, messages...)
		l.Logger.Error(msg, args...)
	}
}

// WithFields creates a new Logger with pre-configured structured fields.
// This is useful for creating contextual loggers that automatically include
// common fields like user_id, request_id, etc. in all log messages.
//
// Example:
//
//	userLogger := logger.WithFields(map[string]any{
//		"user_id": 12345,
//		"tenant_id": "tenant-abc",
//	})
//	userLogger.Info("User performed action", "action", "login")
func (l *Logger) WithFields(fields map[string]any) *Logger {
	args := make([]any, 0, len(fields)*2)
	for k, v := range fields {
		args = append(args, k, v)
	}
	return &Logger{
		Logger: l.Logger.With(args...),
	}
}

// WithError creates a new Logger with a pre-configured error field.
// This is useful when you need to log multiple messages related to the same error.
//
// Example:
//
//	errorLogger := logger.WithError(dbErr)
//	errorLogger.Error("Failed to save user")
//	errorLogger.Warn("Falling back to cache")
func (l *Logger) WithError(err error) *Logger {
	return &Logger{
		Logger: l.Logger.With("error", err),
	}
}

// formatMessage formats the error and additional messages into a log message and arguments.
// It safely handles the case where the first message might not be a string and ensures
// the error is properly included in the structured log output.
func formatMessage(err error, messages ...any) (string, []any) {
	if len(messages) > 0 {
		// Проверяем, что первый элемент - строка
		if msg, ok := messages[0].(string); ok {
			if len(messages) > 1 {
				var args = []any{"error", err}
				return msg, append(args, messages[1:]...)
			}
			return msg, []any{"error", err}
		}
	}
	return err.Error(), []any{}
}
