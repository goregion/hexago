package launcher

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/goregion/hexago/pkg/log"
)

// TestIntegrationScenarios tests real-world usage scenarios
func TestIntegrationScenarios(t *testing.T) {
	logger := log.NewLogger(log.NewTextStdOutHandler())

	t.Run("microservices_simulation", func(t *testing.T) {
		var servicesCompleted sync.Map

		createService := func(name string, duration time.Duration) func(context.Context) error {
			return func(ctx context.Context) error {
				select {
				case <-time.After(duration):
					servicesCompleted.Store(name, true)
					return nil
				case <-ctx.Done():
					servicesCompleted.Store(name, false)
					return ctx.Err()
				}
			}
		}

		authService := createService("auth", 20*time.Millisecond)
		userService := createService("user", 30*time.Millisecond)
		orderService := createService("order", 25*time.Millisecond)

		result := NewAppLauncher().
			WithLoggerContext(logger).
			WithTimeout(100*time.Millisecond). // Enough time for all
			WaitApplications(authService, userService, orderService)

		if result.Error() != nil {
			t.Errorf("Expected no error, got: %v", result.Error())
		}

		// Check all services completed
		services := []string{"auth", "user", "order"}
		for _, service := range services {
			if completed, ok := servicesCompleted.Load(service); !ok || !completed.(bool) {
				t.Errorf("Service %s did not complete successfully", service)
			}
		}
	})

	t.Run("timeout_behavior", func(t *testing.T) {
		longRunningTask := func(ctx context.Context) error {
			select {
			case <-time.After(200 * time.Millisecond):
				return nil
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		start := time.Now()
		result := NewAppLauncher().
			WithLoggerContext(logger).
			WithTimeout(50 * time.Millisecond).
			WaitApplication(longRunningTask)
		duration := time.Since(start)

		if result.Error() == nil {
			t.Error("Expected timeout error")
		}
		if duration > 70*time.Millisecond {
			t.Errorf("Task should have been cancelled around 50ms, took: %v", duration)
		}
	})
}

// TestErrorPropagation tests how errors are handled and propagated
func TestErrorPropagation(t *testing.T) {
	customError := errors.New("custom business logic error")

	t.Run("single_task_error", func(t *testing.T) {
		task := func(ctx context.Context) error {
			return customError
		}

		result := NewAppLauncher().WaitApplication(task)
		if !errors.Is(result.Error(), customError) {
			t.Errorf("Expected custom error, got: %v", result.Error())
		}
	})

	t.Run("multiple_tasks_first_error_wins", func(t *testing.T) {
		fastFailTask := func(ctx context.Context) error {
			return customError
		}

		slowSuccessTask := func(ctx context.Context) error {
			time.Sleep(100 * time.Millisecond)
			return nil
		}

		result := NewAppLauncher().WaitApplications(fastFailTask, slowSuccessTask)
		if result.Error() == nil {
			t.Error("Expected error from fast failing task")
		}
	})
}

// TestMemoryAndResourceUsage tests resource management
func TestMemoryAndResourceUsage(t *testing.T) {
	t.Run("no_goroutine_leak", func(t *testing.T) {
		// This test would ideally check for goroutine leaks
		// For now, we'll just ensure that multiple operations don't cause issues

		for i := 0; i < 100; i++ {
			task := func(ctx context.Context) error {
				return nil
			}

			result := NewAppLauncher().
				WithTimeout(10 * time.Millisecond).
				WithoutTimeout().
				WaitApplication(task)

			if result.Error() != nil {
				t.Errorf("Iteration %d failed: %v", i, result.Error())
			}
		}
	})

	t.Run("context_cleanup", func(t *testing.T) {
		launcher := NewAppLauncher()

		// Set and replace timeout multiple times
		for i := 0; i < 10; i++ {
			launcher = launcher.WithTimeout(time.Duration(i+1) * 10 * time.Millisecond)
		}

		// Remove timeout
		launcher = launcher.WithoutTimeout()

		// Should still work normally
		task := func(ctx context.Context) error { return nil }
		result := launcher.WaitApplication(task)

		if result.Error() != nil {
			t.Errorf("Expected no error after cleanup, got: %v", result.Error())
		}
	})
}

// TestChainedOperations tests complex fluent API usage
func TestChainedOperations(t *testing.T) {
	logger := log.NewLogger(log.NewTextStdOutHandler())

	t.Run("long_chain", func(t *testing.T) {
		task := func(ctx context.Context) error {
			// Verify all context values are available
			if _, err := log.GetLoggerFromContext(ctx); err != nil {
				return errors.New("logger not found in context")
			}
			if ctx.Value("config") == nil {
				return errors.New("config not found in context")
			}
			if ctx.Value("version") == nil {
				return errors.New("version not found in context")
			}
			return nil
		}

		result := NewAppLauncher().
			WithLoggerContext(logger).
			WithGrexitContext().
			WithContext("config", map[string]string{"env": "test"}).
			WithContext("version", "1.0.0").
			WithTimeout(1 * time.Second).
			WithTimeout(2 * time.Second). // Replace previous timeout
			WithoutTimeout().             // Remove timeout
			WithTimeout(5 * time.Second). // Set final timeout
			WaitApplication(task)

		if result.Error() != nil {
			t.Errorf("Long chain failed: %v", result.Error())
		}
	})

	t.Run("mixed_operations", func(t *testing.T) {
		task1 := func(ctx context.Context) error {
			if ctx.Value("task-config") == nil {
				return errors.New("task config missing")
			}
			time.Sleep(10 * time.Millisecond)
			return nil
		}

		task2 := func(ctx context.Context) error {
			if ctx.Value("task-config") == nil {
				return errors.New("task config missing")
			}
			time.Sleep(20 * time.Millisecond)
			return nil
		}

		task3 := func(ctx context.Context) error {
			if ctx.Value("task-config") == nil {
				return errors.New("task config missing")
			}
			time.Sleep(30 * time.Millisecond)
			return nil
		}

		result := NewAppLauncher().
			WithLoggerContext(logger).
			WithContext("task-config", "parallel-execution").
			WithTimeout(200*time.Millisecond).
			WaitApplications(task1, task2, task3)

		if result.Error() != nil {
			t.Errorf("Mixed operations failed: %v", result.Error())
		}
	})
}

// TestEdgeCasesAndBoundaryConditions tests unusual but valid scenarios
func TestEdgeCasesAndBoundaryConditions(t *testing.T) {
	t.Run("immediate_task_completion", func(t *testing.T) {
		task := func(ctx context.Context) error {
			// Complete immediately
			return nil
		}

		start := time.Now()
		result := NewAppLauncher().
			WithTimeout(1 * time.Second).
			WaitApplication(task)
		duration := time.Since(start)

		if result.Error() != nil {
			t.Errorf("Expected no error, got: %v", result.Error())
		}
		if duration > 10*time.Millisecond {
			t.Errorf("Task took too long for immediate completion: %v", duration)
		}
	})

	t.Run("very_short_timeout", func(t *testing.T) {
		task := func(ctx context.Context) error {
			time.Sleep(10 * time.Millisecond)
			return nil
		}

		result := NewAppLauncher().
			WithTimeout(1 * time.Nanosecond). // Extremely short
			WaitApplication(task)

		if result.Error() == nil {
			t.Error("Expected timeout error with very short timeout")
		}
	})

	t.Run("empty_context_values", func(t *testing.T) {
		task := func(ctx context.Context) error { return nil }

		result := NewAppLauncher().
			WithContext("", "empty-key").
			WithContext("empty-value", "").
			WithContext("nil-value", nil).
			WaitApplication(task)

		if result.Error() != nil {
			t.Errorf("Expected no error with empty context values, got: %v", result.Error())
		}
	})
}

// BenchmarkComplexScenarios benchmarks realistic usage patterns
func BenchmarkComplexScenarios(b *testing.B) {
	logger := log.NewLogger(log.NewTextStdOutHandler())

	b.Run("single_task_with_context", func(b *testing.B) {
		task := func(ctx context.Context) error {
			return nil
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			result := NewAppLauncher().
				WithLoggerContext(logger).
				WithContext("iteration", i).
				WaitApplication(task)

			if result.Error() != nil {
				b.Fatal(result.Error())
			}
		}
	})

	b.Run("parallel_tasks", func(b *testing.B) {
		task1 := func(ctx context.Context) error { return nil }
		task2 := func(ctx context.Context) error { return nil }
		task3 := func(ctx context.Context) error { return nil }
		task4 := func(ctx context.Context) error { return nil }
		task5 := func(ctx context.Context) error { return nil }

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			result := NewAppLauncher().
				WithLoggerContext(logger).
				WaitApplications(task1, task2, task3, task4, task5)

			if result.Error() != nil {
				b.Fatal(result.Error())
			}
		}
	})

	b.Run("timeout_operations", func(b *testing.B) {
		task := func(ctx context.Context) error {
			return nil
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			result := NewAppLauncher().
				WithTimeout(1 * time.Second).
				WithoutTimeout().
				WithTimeout(2 * time.Second).
				WaitApplication(task)

			if result.Error() != nil {
				b.Fatal(result.Error())
			}
		}
	})
}
