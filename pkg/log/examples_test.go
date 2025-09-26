// Package log examples and tests demonstrate the usage patterns and functionality
// of the structured logging library. These examples serve as both documentation
// and validation of the library's features.
package log

import (
	"bytes"
	"context"
	"errors"
	"os"
	"testing"
)

// ExampleLogger_basic demonstrates the most basic usage of the logger
// with JSON formatting output to stdout.
func ExampleLogger_basic() {
	// Create logger with JSON output
	logger := NewLogger(NewJsonStdOutHandler())

	logger.Info("Application started", "version", "1.0.0")
	logger.Error("Something went wrong", "error", errors.New("test error"))
}

// ExampleLogger_multipleHandlers shows how to configure a logger
// that outputs to multiple destinations simultaneously.
func ExampleLogger_multipleHandlers() {
	// Create custom stderr handler using existing functions
	var stderrBuf bytes.Buffer
	stderrHandler := NewTextHandler(&stderrBuf) // This simulates stderr

	// Logger with multiple handlers
	logger := NewLogger(
		NewJsonStdOutHandler(),
		stderrHandler,
	)

	logger.Warn("This will appear in both stdout (JSON) and stderr (text)")
}

// ExampleLogger_withContext demonstrates how to store and retrieve
// logger instances from Go context for propagation through call stacks.
func ExampleLogger_withContext() {
	logger := NewLogger(NewTextStdOutHandler())
	ctx := WithLoggerContext(context.Background(), logger)

	// Retrieve logger from context
	loggerFromCtx := MustGetLoggerFromContext(ctx)
	loggerFromCtx.Info("Message from context logger")
}

// ExampleLogger_serviceLifecycle demonstrates automatic service lifecycle logging
// with start and stop messages for service boundary tracking.
func ExampleLogger_serviceLifecycle() {
	logger := NewLogger(NewTextStdOutHandler())

	// Start service with automatic start/stop logging
	serviceLogger, stop := logger.StartService("user-service")
	defer stop() // Automatically logs service stop

	serviceLogger.Info("Processing user request", "user_id", 12345)
}

// ExampleLogger_withFields shows how to create loggers with pre-configured fields
// for consistent context across multiple log messages.
func ExampleLogger_withFields() {
	logger := NewLogger(NewTextStdOutHandler())

	// Create logger with pre-configured fields
	userLogger := logger.WithFields(map[string]any{
		"user_id":   12345,
		"tenant_id": "tenant-abc",
	})

	userLogger.Info("User action performed", "action", "login")
}

// ExampleLogger_errorHandling demonstrates various error logging patterns
// including conditional logging and error context preservation.
func ExampleLogger_errorHandling() {
	logger := NewLogger(NewTextStdOutHandler())

	err := errors.New("database connection failed")

	// Logs only if error exists and is not context cancellation
	logger.LogIfError(err, "Failed to connect to database", "retry_count", 3)

	// Logger with pre-configured error
	errorLogger := logger.WithError(err)
	errorLogger.Error("Operation failed")
}

// TestLogger_GetLoggerFromContext verifies successful logger retrieval from context.
func TestLogger_GetLoggerFromContext(t *testing.T) {
	logger := NewLogger(NewTextStdOutHandler())
	ctx := WithLoggerContext(context.Background(), logger)

	retrievedLogger, err := GetLoggerFromContext(ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if retrievedLogger != logger {
		t.Error("Retrieved logger doesn't match the original")
	}
}

// TestLogger_GetLoggerFromContext_NotFound verifies proper error handling
// when no logger exists in context.
func TestLogger_GetLoggerFromContext_NotFound(t *testing.T) {
	ctx := context.Background()

	_, err := GetLoggerFromContext(ctx)
	if err == nil {
		t.Error("Expected error when logger not in context")
	}
}

// TestMultiHandler_WithNoHandlers verifies that logger works correctly
// even when created without explicit handlers (should use default).
func TestMultiHandler_WithNoHandlers(t *testing.T) {
	// Test logger without handlers - should create default
	logger := NewLogger()

	// This test verifies that logger doesn't panic
	logger.Info("Test message")
}

// BenchmarkLogger_Info measures the performance of basic info logging
// to help identify performance regressions.
func BenchmarkLogger_Info(b *testing.B) {
	logger := NewLogger(NewJsonHandler(os.Stdout))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("Benchmark message", "iteration", i)
	}
}
