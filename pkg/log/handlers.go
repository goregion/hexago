// Package log provides factory functions for creating various slog.Handler implementations.
// This file contains pre-configured handlers for common logging scenarios including
// text and JSON formatting to different output destinations.
package log

import (
	"io"
	"log/slog"
	"os"
)

// envEnableDebugLogLevel is the environment variable name that controls debug logging.
// Set this to "true" to enable DEBUG level logging, otherwise INFO level is used.
const envEnableDebugLogLevel = "ENABLE_DEBUG_LOG_LEVEL"

// makeOptions creates slog.HandlerOptions with appropriate log level based on environment.
// It checks the ENABLE_DEBUG_LOG_LEVEL environment variable to determine the log level.
func makeOptions() *slog.HandlerOptions {
	var logLevel = slog.LevelInfo
	if os.Getenv(envEnableDebugLogLevel) == "true" {
		logLevel = slog.LevelDebug
	}
	return &slog.HandlerOptions{
		Level: logLevel,
	}
}

// NewTextStdOutHandler creates a text-formatted handler that writes to stdout.
// This is useful for development environments where human-readable logs are preferred.
func NewTextStdOutHandler() slog.Handler {
	return slog.NewTextHandler(os.Stdout, makeOptions())
}

// NewJsonStdOutHandler creates a JSON-formatted handler that writes to stdout.
// This is ideal for production environments where structured logging is required
// for log aggregation and analysis tools.
func NewJsonStdOutHandler() slog.Handler {
	return slog.NewJSONHandler(os.Stdout, makeOptions())
}

// NewTextStdErrHandler creates a text-formatted handler that writes to stderr.
// This is useful for error logging or when you want to separate regular logs
// from error logs at the output level.
func NewTextStdErrHandler() slog.Handler {
	return slog.NewTextHandler(os.Stderr, makeOptions())
}

// NewJsonStdErrHandler creates a JSON-formatted handler that writes to stderr.
// Combines the benefits of structured JSON logging with stderr output separation.
func NewJsonStdErrHandler() slog.Handler {
	return slog.NewJSONHandler(os.Stderr, makeOptions())
}

// NewTextHandler creates a text-formatted handler that writes to the specified writer.
// This provides flexibility to log to any destination that implements io.Writer,
// such as files, buffers, or custom writers.
//
// Example:
//
//	file, _ := os.Create("app.log")
//	handler := NewTextHandler(file)
func NewTextHandler(writer io.Writer) slog.Handler {
	return slog.NewTextHandler(writer, makeOptions())
}

// NewJsonHandler creates a JSON-formatted handler that writes to the specified writer.
// This provides structured logging to any destination that implements io.Writer.
//
// Example:
//
//	file, _ := os.Create("app.json")
//	handler := NewJsonHandler(file)
func NewJsonHandler(writer io.Writer) slog.Handler {
	return slog.NewJSONHandler(writer, makeOptions())
}

// NewHandlerWithLevel creates a handler with custom log level configuration.
// This allows fine-grained control over logging levels, bypassing the environment-based configuration.
//
// Parameters:
//   - handler: A function that creates slog.Handler (e.g., slog.NewJSONHandler)
//   - writer: The io.Writer to write logs to
//   - level: The minimum log level to handle
//
// Example:
//
//	handler := NewHandlerWithLevel(slog.NewJSONHandler, os.Stdout, slog.LevelError)
func NewHandlerWithLevel(handler func(io.Writer, *slog.HandlerOptions) slog.Handler, writer io.Writer, level slog.Level) slog.Handler {
	options := &slog.HandlerOptions{
		Level: level,
	}
	return handler(writer, options)
}
