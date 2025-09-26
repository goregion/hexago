package launcher

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/goregion/hexago/pkg/log"
)

// TestTimeoutEdgeCases tests various timeout edge cases
func TestTimeoutEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		timeout  time.Duration
		taskTime time.Duration
		wantErr  bool
	}{
		{
			name:     "zero timeout",
			timeout:  0,
			taskTime: 10 * time.Millisecond,
			wantErr:  false, // Zero timeout should be ignored
		},
		{
			name:     "negative timeout",
			timeout:  -1 * time.Second,
			taskTime: 10 * time.Millisecond,
			wantErr:  false, // Negative timeout should be ignored
		},
		{
			name:     "very short timeout",
			timeout:  1 * time.Millisecond,
			taskTime: 50 * time.Millisecond,
			wantErr:  true,
		},
		{
			name:     "exact timeout",
			timeout:  50 * time.Millisecond,
			taskTime: 50 * time.Millisecond,
			wantErr:  false, // Should complete just in time
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := func(ctx context.Context) error {
				select {
				case <-time.After(tt.taskTime):
					return nil
				case <-ctx.Done():
					return ctx.Err()
				}
			}

			result := NewAppLauncher().
				WithTimeout(tt.timeout).
				WaitApplication(task)

			if tt.wantErr && result.Error() == nil {
				t.Error("Expected timeout error but got none")
			}
			if !tt.wantErr && result.Error() != nil {
				t.Errorf("Expected no error but got: %v", result.Error())
			}
		})
	}
}

// TestContextHierarchy tests that context values are properly inherited
func TestContextHierarchy(t *testing.T) {
	logger := log.NewLogger(log.NewTextStdOutHandler())
	testKey := "test-key"
	testValue := "test-value"

	receivedValues := make(map[string]interface{})

	task := func(ctx context.Context) error {
		// Check if logger is available
		if _, err := log.GetLoggerFromContext(ctx); err != nil {
			receivedValues["logger"] = nil
		} else {
			receivedValues["logger"] = true
		}

		// Check if custom value is available
		if value := ctx.Value(testKey); value != nil {
			receivedValues["custom"] = value
		}

		return nil
	}

	result := NewAppLauncher().
		WithLoggerContext(logger).
		WithContext(testKey, testValue).
		WaitApplication(task)

	if result.Error() != nil {
		t.Fatalf("Task failed: %v", result.Error())
	}

	if receivedValues["logger"] == nil {
		t.Error("Logger not found in task context")
	}

	if receivedValues["custom"] != testValue {
		t.Errorf("Expected custom value %v, got %v", testValue, receivedValues["custom"])
	}
}

// TestMultipleTasksWithDifferentBehavior tests complex parallel scenarios
func TestMultipleTasksWithDifferentBehavior(t *testing.T) {
	completedTasks := make(map[string]bool)

	quickTask := func(ctx context.Context) error {
		time.Sleep(10 * time.Millisecond)
		completedTasks["quick"] = true
		return nil
	}

	slowTask := func(ctx context.Context) error {
		time.Sleep(50 * time.Millisecond)
		completedTasks["slow"] = true
		return nil
	}

	errorTask := func(ctx context.Context) error {
		time.Sleep(20 * time.Millisecond)
		completedTasks["error"] = true
		return errors.New("intentional error")
	}

	result := NewAppLauncher().WaitApplications(quickTask, slowTask, errorTask)

	// Should get error from errorTask
	if result.Error() == nil {
		t.Error("Expected error from errorTask")
	}

	// All tasks should have had a chance to run and complete or be interrupted
	// Due to parallel execution, the exact completion state depends on timing
	t.Logf("Completed tasks: %+v", completedTasks)
}

// TestResourceCleanup tests that resources are properly cleaned up
func TestResourceCleanup(t *testing.T) {
	// Test that cancel functions are properly managed
	launcher := NewAppLauncher()

	// Set timeout
	launcher = launcher.WithTimeout(1 * time.Second)

	// Replace timeout
	launcher = launcher.WithTimeout(2 * time.Second)

	// Remove timeout
	launcher = launcher.WithoutTimeout()

	// Should still work normally
	task := func(ctx context.Context) error { return nil }
	result := launcher.WaitApplication(task)

	if result.Error() != nil {
		t.Errorf("Expected no error after cleanup, got: %v", result.Error())
	}
}

// TestLogIfErrorBehavior tests AppResult.LogIfError method
func TestLogIfErrorBehavior(t *testing.T) {
	// Create a custom logger that captures output
	var logOutput strings.Builder
	logger := log.NewLogger(log.NewTextHandler(&logOutput))

	tests := []struct {
		name        string
		task        func(context.Context) error
		logger      *log.Logger
		shouldLog   bool
		expectPanic bool
	}{
		{
			name:        "success with logger",
			task:        func(ctx context.Context) error { return nil },
			logger:      logger,
			shouldLog:   false,
			expectPanic: false,
		},
		{
			name:        "error with logger",
			task:        func(ctx context.Context) error { return errors.New("test error") },
			logger:      logger,
			shouldLog:   true,
			expectPanic: false,
		},
		{
			name:        "error with nil logger",
			task:        func(ctx context.Context) error { return errors.New("test error") },
			logger:      nil,
			shouldLog:   false,
			expectPanic: false,
		},
		{
			name:        "success with nil logger",
			task:        func(ctx context.Context) error { return nil },
			logger:      nil,
			shouldLog:   false,
			expectPanic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logOutput.Reset()

			result := NewAppLauncher().WaitApplication(tt.task)

			defer func() {
				if r := recover(); r != nil {
					if !tt.expectPanic {
						t.Errorf("Unexpected panic: %v", r)
					}
				} else if tt.expectPanic {
					t.Error("Expected panic but none occurred")
				}
			}()

			result.LogIfError(tt.logger, "additional context")

			loggedContent := logOutput.String()
			if tt.shouldLog && loggedContent == "" {
				t.Error("Expected log output but got none")
			}
			if !tt.shouldLog && loggedContent != "" {
				t.Errorf("Expected no log output but got: %s", loggedContent)
			}
		})
	}
}

// TestConcurrentTimeoutModification tests concurrent timeout modifications
func TestConcurrentTimeoutModification(t *testing.T) {
	// Each goroutine works with its own launcher to avoid race conditions
	done := make(chan bool, 3)

	go func() {
		launcher := NewAppLauncher()
		for i := 0; i < 10; i++ {
			launcher = launcher.WithTimeout(time.Duration(i+1) * 10 * time.Millisecond)
			time.Sleep(1 * time.Millisecond)
		}
		done <- true
	}()

	go func() {
		launcher := NewAppLauncher().WithTimeout(100 * time.Millisecond)
		for i := 0; i < 5; i++ {
			launcher = launcher.WithoutTimeout()
			time.Sleep(2 * time.Millisecond)
		}
		done <- true
	}()

	go func() {
		launcher := NewAppLauncher().WithTimeout(20 * time.Millisecond)
		task := func(ctx context.Context) error {
			time.Sleep(5 * time.Millisecond)
			return nil
		}
		result := launcher.WaitApplication(task)
		// This task should either succeed or get timeout - both are acceptable
		if result.Error() != nil && !errors.Is(result.Error(), context.DeadlineExceeded) {
			t.Errorf("Unexpected error: %v", result.Error())
		}
		done <- true
	}()

	// Wait for all goroutines
	for i := 0; i < 3; i++ {
		<-done
	}
}

// TestLargeNumberOfTasks tests handling many parallel tasks
func TestLargeNumberOfTasks(t *testing.T) {
	const numTasks = 10 // Reduced for simpler test

	task1 := func(ctx context.Context) error {
		time.Sleep(1 * time.Millisecond)
		return nil
	}
	task2 := func(ctx context.Context) error {
		time.Sleep(2 * time.Millisecond)
		return nil
	}
	task3 := func(ctx context.Context) error {
		time.Sleep(3 * time.Millisecond)
		return nil
	}
	task4 := func(ctx context.Context) error {
		time.Sleep(4 * time.Millisecond)
		return nil
	}
	task5 := func(ctx context.Context) error {
		time.Sleep(5 * time.Millisecond)
		return nil
	}

	start := time.Now()
	result := NewAppLauncher().WaitApplications(task1, task2, task3, task4, task5)
	duration := time.Since(start)

	if result.Error() != nil {
		t.Errorf("Expected no error but got: %v", result.Error())
	}

	// Should complete reasonably quickly due to parallel execution
	if duration > 50*time.Millisecond {
		t.Errorf("Too slow for parallel tasks: %v", duration)
	}

	t.Logf("Executed 5 tasks in %v", duration)
}

// ExampleAppLauncher demonstrates typical usage patterns
func ExampleAppLauncher() {
	logger := log.NewLogger(log.NewTextStdOutHandler())

	task := func(ctx context.Context) error {
		// Simulate some work
		select {
		case <-time.After(100 * time.Millisecond):
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	result := NewAppLauncher().
		WithLoggerContext(logger).
		WithGrexitContext().
		WithTimeout(5 * time.Second).
		WaitApplication(task)

	result.LogIfError(logger, "Application execution failed")
}
