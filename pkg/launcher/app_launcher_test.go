package launcher

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/goregion/hexago/pkg/log"
)

// TestNewAppLauncher tests basic launcher creation
func TestNewAppLauncher(t *testing.T) {
	launcher := NewAppLauncher()
	if launcher == nil {
		t.Fatal("NewAppLauncher() returned nil")
	}

	// Test with nil context
	launcher2 := NewAppLauncherWithContext(nil)
	if launcher2 == nil {
		t.Fatal("NewAppLauncherWithContext(nil) returned nil")
	}

	// Test with valid context
	ctx := context.Background()
	launcher3 := NewAppLauncherWithContext(ctx)
	if launcher3 == nil {
		t.Fatal("NewAppLauncherWithContext(ctx) returned nil")
	}
}

// TestWithLoggerContext tests logger context enrichment
func TestWithLoggerContext(t *testing.T) {
	logger := log.NewLogger(log.NewTextStdOutHandler())

	launcher := NewAppLauncher().WithLoggerContext(logger)
	if launcher == nil {
		t.Fatal("WithLoggerContext returned nil")
	}

	// Test with nil logger - should not panic
	launcher2 := NewAppLauncher().WithLoggerContext(nil)
	if launcher2 == nil {
		t.Fatal("WithLoggerContext(nil) returned nil")
	}
}

// TestWithContext tests custom context enrichment
func TestWithContext(t *testing.T) {
	key := "test-key"
	value := "test-value"

	launcher := NewAppLauncher().WithContext(key, value)
	if launcher == nil {
		t.Fatal("WithContext returned nil")
	}

	// Test with nil key - should not panic
	launcher2 := NewAppLauncher().WithContext(nil, value)
	if launcher2 == nil {
		t.Fatal("WithContext(nil, value) returned nil")
	}
}

// TestWaitApplication tests single application execution
func TestWaitApplication(t *testing.T) {
	tests := []struct {
		name        string
		task        func(context.Context) error
		expectError bool
	}{
		{
			name: "successful task",
			task: func(ctx context.Context) error {
				return nil
			},
			expectError: false,
		},
		{
			name: "failing task",
			task: func(ctx context.Context) error {
				return errors.New("task failed")
			},
			expectError: true,
		},
		{
			name: "context-aware task",
			task: func(ctx context.Context) error {
				select {
				case <-time.After(10 * time.Millisecond):
					return nil
				case <-ctx.Done():
					return ctx.Err()
				}
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewAppLauncher().WaitApplication(tt.task)

			if tt.expectError && result.Error() == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && result.Error() != nil {
				t.Errorf("Expected no error but got: %v", result.Error())
			}
		})
	}
}

// TestWaitApplicationNilTask tests nil task handling
func TestWaitApplicationNilTask(t *testing.T) {
	result := NewAppLauncher().WaitApplication(nil)
	if result.Error() == nil {
		t.Error("Expected error for nil task")
	}
	if result.Error().Error() != "task cannot be nil" {
		t.Errorf("Expected 'task cannot be nil' error, got: %v", result.Error())
	}
}

// TestWaitApplications tests multiple application execution
func TestWaitApplications(t *testing.T) {
	task1 := func(ctx context.Context) error {
		time.Sleep(10 * time.Millisecond)
		return nil
	}

	task2 := func(ctx context.Context) error {
		time.Sleep(15 * time.Millisecond)
		return nil
	}

	start := time.Now()
	result := NewAppLauncher().WaitApplications(task1, task2)
	duration := time.Since(start)

	if result.Error() != nil {
		t.Errorf("Expected no error but got: %v", result.Error())
	}

	// Tasks should run in parallel, so total time should be ~15ms (max), not 25ms (sum)
	if duration > 50*time.Millisecond {
		t.Errorf("Tasks took too long: %v, expected ~15ms (parallel execution)", duration)
	}
}

// TestWaitApplicationsWithError tests error handling in parallel execution
func TestWaitApplicationsWithError(t *testing.T) {
	task1 := func(ctx context.Context) error {
		return nil
	}

	task2 := func(ctx context.Context) error {
		return errors.New("task2 failed")
	}

	result := NewAppLauncher().WaitApplications(task1, task2)
	if result.Error() == nil {
		t.Error("Expected error from failing task")
	}
}

// TestWaitApplicationsValidation tests input validation
func TestWaitApplicationsValidation(t *testing.T) {
	// Test empty tasks
	result := NewAppLauncher().WaitApplications()
	if result.Error() == nil {
		t.Error("Expected error for empty tasks")
	}

	// Test with nil task
	validTask := func(ctx context.Context) error { return nil }
	result2 := NewAppLauncher().WaitApplications(validTask, nil)
	if result2.Error() == nil {
		t.Error("Expected error for nil task in array")
	}
}

// TestWithTimeout tests timeout functionality
func TestWithTimeout(t *testing.T) {
	longTask := func(ctx context.Context) error {
		select {
		case <-time.After(100 * time.Millisecond):
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	start := time.Now()
	result := NewAppLauncher().
		WithTimeout(50 * time.Millisecond).
		WaitApplication(longTask)
	duration := time.Since(start)

	if result.Error() == nil {
		t.Error("Expected timeout error")
	}
	if !errors.Is(result.Error(), context.DeadlineExceeded) {
		t.Errorf("Expected context.DeadlineExceeded, got: %v", result.Error())
	}
	if duration > 70*time.Millisecond {
		t.Errorf("Task took too long: %v, expected ~50ms", duration)
	}
}

// TestWithTimeoutReplacement tests that multiple timeouts are handled correctly
func TestWithTimeoutReplacement(t *testing.T) {
	longTask := func(ctx context.Context) error {
		select {
		case <-time.After(200 * time.Millisecond):
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	start := time.Now()
	result := NewAppLauncher().
		WithTimeout(300 * time.Millisecond). // First timeout - should be replaced
		WithTimeout(100 * time.Millisecond). // Second timeout - should be used
		WaitApplication(longTask)
	duration := time.Since(start)

	if result.Error() == nil {
		t.Error("Expected timeout error")
	}
	// Should timeout after ~100ms, not 300ms
	if duration > 150*time.Millisecond {
		t.Errorf("Task took too long: %v, expected ~100ms (second timeout should be used)", duration)
	}
}

// TestWithoutTimeout tests timeout removal
func TestWithoutTimeout(t *testing.T) {
	quickTask := func(ctx context.Context) error {
		time.Sleep(20 * time.Millisecond)
		return nil
	}

	result := NewAppLauncher().
		WithTimeout(10 * time.Millisecond). // Set a short timeout
		WithoutTimeout().                   // Remove the timeout
		WaitApplication(quickTask)

	// Task should complete successfully without timeout
	if result.Error() != nil {
		t.Errorf("Expected no error after removing timeout, got: %v", result.Error())
	}
}

// TestFluentAPI tests method chaining
func TestFluentAPI(t *testing.T) {
	logger := log.NewLogger(log.NewTextStdOutHandler())

	task := func(ctx context.Context) error {
		return nil
	}

	// Test that all methods return the same interface for chaining
	result := NewAppLauncher().
		WithLoggerContext(logger).
		WithGrexitContext().
		WithContext("key", "value").
		WithTimeout(1 * time.Second).
		WithoutTimeout().
		WaitApplication(task)

	if result.Error() != nil {
		t.Errorf("Fluent API chain failed: %v", result.Error())
	}
}

// TestAppResult tests result interface
func TestAppResult(t *testing.T) {
	logger := log.NewLogger(log.NewTextStdOutHandler())

	// Test successful result
	successTask := func(ctx context.Context) error { return nil }
	result1 := NewAppLauncher().WaitApplication(successTask)

	if result1.Error() != nil {
		t.Error("Expected no error for successful task")
	}

	// LogIfError should not panic with nil or valid logger
	result1.LogIfError(nil)
	result1.LogIfError(logger)

	// Test error result
	errorTask := func(ctx context.Context) error { return errors.New("test error") }
	result2 := NewAppLauncher().WaitApplication(errorTask)

	if result2.Error() == nil {
		t.Error("Expected error for failing task")
	}

	// LogIfError should handle error case
	result2.LogIfError(logger, "additional", "context")
}

// TestConcurrentUsage tests that launcher is safe for concurrent use
func TestConcurrentUsage(t *testing.T) {
	const numGoroutines = 10
	var completed int32

	task := func(ctx context.Context) error {
		time.Sleep(10 * time.Millisecond)
		atomic.AddInt32(&completed, 1)
		return nil
	}

	// Launch multiple launchers concurrently
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			result := NewAppLauncher().
				WithContext("id", id).
				WaitApplication(task)

			if result.Error() != nil {
				t.Errorf("Goroutine %d failed: %v", id, result.Error())
			}
		}(i)
	}

	// Wait for all to complete
	time.Sleep(100 * time.Millisecond)

	if atomic.LoadInt32(&completed) != numGoroutines {
		t.Errorf("Expected %d completed tasks, got %d", numGoroutines, atomic.LoadInt32(&completed))
	}
}

// BenchmarkLauncherCreation benchmarks launcher creation
func BenchmarkLauncherCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewAppLauncher()
	}
}

// BenchmarkSimpleTask benchmarks simple task execution
func BenchmarkSimpleTask(b *testing.B) {
	task := func(ctx context.Context) error { return nil }

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := NewAppLauncher().WaitApplication(task)
		if result.Error() != nil {
			b.Fatal(result.Error())
		}
	}
}

// BenchmarkFluentAPI benchmarks fluent API chain
func BenchmarkFluentAPI(b *testing.B) {
	logger := log.NewLogger(log.NewTextStdOutHandler())
	task := func(ctx context.Context) error { return nil }

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := NewAppLauncher().
			WithLoggerContext(logger).
			WithGrexitContext().
			WithContext("key", i).
			WaitApplication(task)

		if result.Error() != nil {
			b.Fatal(result.Error())
		}
	}
}
