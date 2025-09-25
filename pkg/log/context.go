// Package log provides context integration utilities for logger instances.
// This file contains functions for storing and retrieving loggers from Go contexts,
// enabling logger propagation throughout the application call stack.
package log

import (
	"context"
	"fmt"
)

const defaultLoggerContextKey = "logger"

type loggerContextKey string

// LoggerContextKey is the key used to store logger instances in context.
// It can be customized using SetLoggerContextKey function.
var LoggerContextKey = loggerContextKey(defaultLoggerContextKey)

// WithLoggerContext adds a logger instance to the given context.
// This allows the logger to be passed through the application call stack
// without explicit parameter passing.
//
// Example:
//
//	ctx := WithLoggerContext(context.Background(), logger)
//	// Pass ctx to other functions that need logging
func WithLoggerContext(ctx context.Context, logger *Logger) context.Context {
	return context.WithValue(ctx, LoggerContextKey, logger)
}

// MustGetLoggerFromContext retrieves a logger from the context and panics if not found.
// Use this only when you are certain that a logger exists in the context.
// For safer retrieval, use GetLoggerFromContext instead.
//
// Example:
//
//	logger := MustGetLoggerFromContext(ctx)
//	logger.Info("Operation completed")
func MustGetLoggerFromContext(ctx context.Context) *Logger {
	logger, err := GetLoggerFromContext(ctx)
	if err != nil {
		panic(err)
	}
	return logger
}

// GetLoggerFromContext safely retrieves a logger from the context.
// Returns an error if no logger is found with the current LoggerContextKey.
// This is the preferred method for logger retrieval when the presence
// of a logger in context is not guaranteed.
//
// Example:
//
//	logger, err := GetLoggerFromContext(ctx)
//	if err != nil {
//	    // Handle missing logger case
//	    return err
//	}
func GetLoggerFromContext(ctx context.Context) (*Logger, error) {
	if logger, ok := ctx.Value(LoggerContextKey).(*Logger); ok {
		return logger, nil
	}
	return nil, fmt.Errorf("logger not found in context with key '%s'", LoggerContextKey)
}

// SetLoggerContextKey allows customization of the context key used for logger storage.
// This is useful when you need to avoid conflicts with other packages that might
// use the same context key. Call this function during application initialization.
//
// Example:
//
//	SetLoggerContextKey("my-app-logger")
func SetLoggerContextKey(key string) {
	LoggerContextKey = loggerContextKey(key)
}
