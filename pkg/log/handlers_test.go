package log

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"strings"
	"testing"
	"time"
)

// TestNewTextHandler verifies text handler creates proper formatted output
func TestNewTextHandler(t *testing.T) {
	var buf bytes.Buffer
	handler := NewTextHandler(&buf)
	logger := slog.New(handler)

	logger.Info("test message", "key", "value")

	output := buf.String()
	if !strings.Contains(output, "test message") {
		t.Errorf("Expected log message to contain 'test message', got: %s", output)
	}
	if !strings.Contains(output, "key=value") {
		t.Errorf("Expected log message to contain 'key=value', got: %s", output)
	}
}

// TestNewJsonHandler verifies JSON handler creates valid JSON output
func TestNewJsonHandler(t *testing.T) {
	var buf bytes.Buffer
	handler := NewJsonHandler(&buf)
	logger := slog.New(handler)

	logger.Info("test message", "key", "value", "number", 42)

	// Parse JSON to verify it's valid
	var logEntry map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
		t.Fatalf("Failed to parse JSON output: %v", err)
	}

	// Verify expected fields
	if logEntry["msg"] != "test message" {
		t.Errorf("Expected msg to be 'test message', got: %v", logEntry["msg"])
	}
	if logEntry["key"] != "value" {
		t.Errorf("Expected key to be 'value', got: %v", logEntry["key"])
	}
	if logEntry["number"] != float64(42) { // JSON numbers are float64
		t.Errorf("Expected number to be 42, got: %v", logEntry["number"])
	}
}

// TestMakeOptions verifies log level configuration based on environment
func TestMakeOptions(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		expected slog.Level
	}{
		{"default level", "", slog.LevelInfo},
		{"debug enabled", "true", slog.LevelDebug},
		{"debug disabled", "false", slog.LevelInfo},
		{"invalid value", "invalid", slog.LevelInfo},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variable
			if tt.envValue != "" {
				os.Setenv(envEnableDebugLogLevel, tt.envValue)
			} else {
				os.Unsetenv(envEnableDebugLogLevel)
			}
			defer os.Unsetenv(envEnableDebugLogLevel)

			options := makeOptions()
			if options.Level != tt.expected {
				t.Errorf("Expected level %v, got %v", tt.expected, options.Level)
			}
		})
	}
}

// TestCustomLevelHandler verifies custom level configuration works correctly
func TestCustomLevelHandler(t *testing.T) {
	var buf bytes.Buffer

	// Use slog directly with custom level
	options := &slog.HandlerOptions{Level: slog.LevelError}
	handler := slog.NewJSONHandler(&buf, options)
	logger := slog.New(handler)

	// This should not be logged (below ERROR level)
	logger.Info("info message")
	logger.Warn("warn message")

	// This should be logged (ERROR level)
	logger.Error("error message")

	output := buf.String()

	// Should not contain info or warn messages
	if strings.Contains(output, "info message") {
		t.Error("Info message should not be logged at ERROR level")
	}
	if strings.Contains(output, "warn message") {
		t.Error("Warn message should not be logged at ERROR level")
	}

	// Should contain error message
	if !strings.Contains(output, "error message") {
		t.Error("Error message should be logged at ERROR level")
	}
}

// TestNewTextStdOutHandler verifies stdout text handler works
func TestNewTextStdOutHandler(t *testing.T) {
	handler := NewTextStdOutHandler()
	if handler == nil {
		t.Error("NewTextStdOutHandler should not return nil")
	}

	// Test that it's enabled for info level
	ctx := context.Background()
	if !handler.Enabled(ctx, slog.LevelInfo) {
		t.Error("Handler should be enabled for INFO level")
	}
}

// TestNewJsonStdOutHandler verifies stdout JSON handler works
func TestNewJsonStdOutHandler(t *testing.T) {
	handler := NewJsonStdOutHandler()
	if handler == nil {
		t.Error("NewJsonStdOutHandler should not return nil")
	}

	// Test that it's enabled for info level
	ctx := context.Background()
	if !handler.Enabled(ctx, slog.LevelInfo) {
		t.Error("Handler should be enabled for INFO level")
	}
}

// TestLogger_WithFields verifies logger creates structured context correctly
func TestLogger_WithFields(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(NewJsonHandler(&buf))

	userLogger := logger.WithFields(map[string]any{
		"user_id":   12345,
		"tenant_id": "tenant-abc",
	})

	userLogger.Info("user action", "action", "login")

	var logEntry map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
		t.Fatalf("Failed to parse JSON output: %v", err)
	}

	// Verify pre-configured fields are present
	if logEntry["user_id"] != float64(12345) {
		t.Errorf("Expected user_id to be 12345, got: %v", logEntry["user_id"])
	}
	if logEntry["tenant_id"] != "tenant-abc" {
		t.Errorf("Expected tenant_id to be 'tenant-abc', got: %v", logEntry["tenant_id"])
	}
	if logEntry["action"] != "login" {
		t.Errorf("Expected action to be 'login', got: %v", logEntry["action"])
	}
}

// TestLogger_StartService verifies service lifecycle logging
func TestLogger_StartService(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(NewJsonHandler(&buf))

	serviceLogger, stop := logger.StartService("test-service")

	// Log something in between
	serviceLogger.Info("processing request")

	stop() // Should log stop message

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")

	if len(lines) != 3 {
		t.Fatalf("Expected 3 log lines, got %d", len(lines))
	}

	// Parse each log line
	var startEntry, processEntry, stopEntry map[string]interface{}

	json.Unmarshal([]byte(lines[0]), &startEntry)
	json.Unmarshal([]byte(lines[1]), &processEntry)
	json.Unmarshal([]byte(lines[2]), &stopEntry)

	// Verify start message
	if startEntry["msg"] != "start" {
		t.Errorf("Expected first message to be 'start', got: %v", startEntry["msg"])
	}
	if startEntry["service"] != "test-service" {
		t.Errorf("Expected service to be 'test-service', got: %v", startEntry["service"])
	}

	// Verify process message
	if processEntry["msg"] != "processing request" {
		t.Errorf("Expected second message to be 'processing request', got: %v", processEntry["msg"])
	}
	if processEntry["service"] != "test-service" {
		t.Errorf("Expected service to be 'test-service', got: %v", processEntry["service"])
	}

	// Verify stop message
	if stopEntry["msg"] != "stop" {
		t.Errorf("Expected third message to be 'stop', got: %v", stopEntry["msg"])
	}
	if stopEntry["service"] != "test-service" {
		t.Errorf("Expected service to be 'test-service', got: %v", stopEntry["service"])
	}
}

// TestLogger_LogIfError verifies conditional error logging
func TestLogger_LogIfError(t *testing.T) {
	tests := []struct {
		name        string
		err         error
		shouldLog   bool
		expectedMsg string
	}{
		{"nil error", nil, false, ""},
		{"regular error", errors.New("test error"), true, "operation failed"},
		{"context canceled", context.Canceled, false, ""},
		{"wrapped context canceled", errors.New("wrapped: context canceled"), true, "operation failed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := NewLogger(NewJsonHandler(&buf))

			logger.LogIfError(tt.err, "operation failed", "retry", 3)

			output := buf.String()

			if tt.shouldLog {
				if output == "" {
					t.Error("Expected log output, got none")
				}

				var logEntry map[string]interface{}
				if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
					t.Fatalf("Failed to parse JSON output: %v", err)
				}

				if logEntry["msg"] != tt.expectedMsg {
					t.Errorf("Expected msg to be '%s', got: %v", tt.expectedMsg, logEntry["msg"])
				}
				if logEntry["retry"] != float64(3) {
					t.Errorf("Expected retry to be 3, got: %v", logEntry["retry"])
				}
				if logEntry["error"] == nil {
					t.Error("Expected error field to be present")
				}
			} else {
				if output != "" {
					t.Errorf("Expected no log output, got: %s", output)
				}
			}
		})
	}
}

// TestMultiHandler_MultipleOutputs verifies multi-handler functionality
func TestMultiHandler_MultipleOutputs(t *testing.T) {
	var buf1, buf2 bytes.Buffer

	handler1 := NewJsonHandler(&buf1)
	handler2 := NewTextHandler(&buf2)

	logger := NewLogger(handler1, handler2)
	logger.Info("test message", "key", "value")

	// Verify JSON output
	var logEntry map[string]interface{}
	if err := json.Unmarshal(buf1.Bytes(), &logEntry); err != nil {
		t.Fatalf("Failed to parse JSON output: %v", err)
	}
	if logEntry["msg"] != "test message" {
		t.Errorf("JSON: Expected msg to be 'test message', got: %v", logEntry["msg"])
	}

	// Verify text output
	textOutput := buf2.String()
	if !strings.Contains(textOutput, "test message") {
		t.Errorf("Text: Expected output to contain 'test message', got: %s", textOutput)
	}
	if !strings.Contains(textOutput, "key=value") {
		t.Errorf("Text: Expected output to contain 'key=value', got: %s", textOutput)
	}
}

// TestContext_LoggerStorage verifies context logger storage and retrieval
func TestContext_LoggerStorage(t *testing.T) {
	logger := NewLogger(NewTextHandler(&bytes.Buffer{}))
	ctx := WithLoggerContext(context.Background(), logger)

	// Test successful retrieval
	retrievedLogger, err := GetLoggerFromContext(ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if retrievedLogger != logger {
		t.Error("Retrieved logger doesn't match original")
	}

	// Test MustGet (should not panic)
	mustLogger := MustGetLoggerFromContext(ctx)
	if mustLogger != logger {
		t.Error("MustGet logger doesn't match original")
	}

	// Test missing logger
	emptyCtx := context.Background()
	_, err = GetLoggerFromContext(emptyCtx)
	if err == nil {
		t.Error("Expected error when logger not in context")
	}
}

// TestFormatMessage verifies error message formatting
func TestFormatMessage(t *testing.T) {
	err := errors.New("test error")

	tests := []struct {
		name        string
		err         error
		messages    []any
		expectedMsg string
		expectedLen int
	}{
		{"string message with args", err, []any{"operation failed", "retry", 3}, "operation failed", 4},
		{"string message only", err, []any{"operation failed"}, "operation failed", 2},
		{"non-string first arg", err, []any{123, "retry", 3}, "test error", 0},
		{"no messages", err, []any{}, "test error", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, args := formatMessage(tt.err, tt.messages...)

			if msg != tt.expectedMsg {
				t.Errorf("Expected message '%s', got '%s'", tt.expectedMsg, msg)
			}
			if len(args) != tt.expectedLen {
				t.Errorf("Expected %d args, got %d", tt.expectedLen, len(args))
			}
		})
	}
}

// BenchmarkLogger_JSONVsText compares performance of JSON vs Text handlers
func BenchmarkLogger_JSONVsText(b *testing.B) {
	tests := []struct {
		name    string
		handler slog.Handler
	}{
		{"JSON", NewJsonHandler(&bytes.Buffer{})},
		{"Text", NewTextHandler(&bytes.Buffer{})},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			logger := NewLogger(tt.handler)
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				logger.Info("benchmark message",
					"iteration", i,
					"timestamp", time.Now(),
					"service", "benchmark-service",
				)
			}
		})
	}
}

// BenchmarkMultiHandler_Performance measures multi-handler overhead
func BenchmarkMultiHandler_Performance(b *testing.B) {
	singleHandler := NewLogger(NewJsonHandler(&bytes.Buffer{}))
	multiHandler := NewLogger(
		NewJsonHandler(&bytes.Buffer{}),
		NewTextHandler(&bytes.Buffer{}),
	)

	b.Run("Single", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			singleHandler.Info("message", "key", "value")
		}
	})

	b.Run("Multi", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			multiHandler.Info("message", "key", "value")
		}
	})
}
