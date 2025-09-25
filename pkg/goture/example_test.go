package goture

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestGoture_Wait(t *testing.T) {
	ctx := context.Background()

	// Test successful task
	t.Run("successful task", func(t *testing.T) {
		task := func(ctx context.Context) error {
			time.Sleep(10 * time.Millisecond)
			return nil
		}

		future := NewGoture(ctx, task)
		err := future.Wait()

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	// Test task with error
	t.Run("task with error", func(t *testing.T) {
		expectedErr := errors.New("task failed")
		task := func(ctx context.Context) error {
			time.Sleep(10 * time.Millisecond)
			return expectedErr
		}

		future := NewGoture(ctx, task)
		err := future.Wait()

		if err == nil {
			t.Error("Expected error, got nil")
		}
		if err.Error() != expectedErr.Error() {
			t.Errorf("Expected error %v, got %v", expectedErr, err)
		}
	})
}

func TestParallelGoture_Wait(t *testing.T) {
	ctx := context.Background()

	// Test multiple successful tasks
	t.Run("multiple successful tasks", func(t *testing.T) {
		task1 := func(ctx context.Context) error {
			time.Sleep(20 * time.Millisecond)
			return nil
		}
		task2 := func(ctx context.Context) error {
			time.Sleep(30 * time.Millisecond)
			return nil
		}
		task3 := func(ctx context.Context) error {
			time.Sleep(10 * time.Millisecond)
			return nil
		}

		start := time.Now()
		future := NewParallelGoture(ctx, task1, task2, task3)
		err := future.Wait()
		duration := time.Since(start)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// Should take about 30ms (longest task), not 60ms (sum of all tasks)
		if duration > 100*time.Millisecond {
			t.Errorf("Tasks didn't run in parallel, took %v", duration)
		}
	})

	// Test with one failing task
	t.Run("one failing task", func(t *testing.T) {
		expectedErr := errors.New("task2 failed")
		task1 := func(ctx context.Context) error {
			time.Sleep(10 * time.Millisecond)
			return nil
		}
		task2 := func(ctx context.Context) error {
			time.Sleep(20 * time.Millisecond)
			return expectedErr
		}
		task3 := func(ctx context.Context) error {
			time.Sleep(30 * time.Millisecond)
			return nil
		}

		future := NewParallelGoture(ctx, task1, task2, task3)
		err := future.Wait()

		if err == nil {
			t.Error("Expected error, got nil")
		}
		if err.Error() != expectedErr.Error() {
			t.Errorf("Expected error %v, got %v", expectedErr, err)
		}
	})

	// Test empty task list
	t.Run("empty task list", func(t *testing.T) {
		future := NewParallelGoture(ctx)
		err := future.Wait()

		if err != nil {
			t.Errorf("Expected no error for empty task list, got %v", err)
		}
	})
}

func TestGotureWithResult_Wait(t *testing.T) {
	ctx := context.Background()

	// Test successful task with result
	t.Run("successful task with result", func(t *testing.T) {
		expectedResult := "hello world"
		task := func(ctx context.Context) (string, error) {
			time.Sleep(10 * time.Millisecond)
			return expectedResult, nil
		}

		future := NewGotureWithResult(ctx, task)
		result, err := future.Wait()

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if result != expectedResult {
			t.Errorf("Expected result %v, got %v", expectedResult, result)
		}
	})

	// Test task with error
	t.Run("task with error", func(t *testing.T) {
		expectedErr := errors.New("task failed")
		task := func(ctx context.Context) (int, error) {
			time.Sleep(10 * time.Millisecond)
			return 0, expectedErr
		}

		future := NewGotureWithResult(ctx, task)
		result, err := future.Wait()

		if err == nil {
			t.Error("Expected error, got nil")
		}
		if err.Error() != expectedErr.Error() {
			t.Errorf("Expected error %v, got %v", expectedErr, err)
		}
		if result != 0 {
			t.Errorf("Expected zero result, got %v", result)
		}
	})
}

func TestNewParallelWithResult(t *testing.T) {
	ctx := context.Background()

	// Test multiple successful tasks with results
	t.Run("multiple successful tasks with results", func(t *testing.T) {
		task1 := func(ctx context.Context) (string, error) {
			time.Sleep(20 * time.Millisecond)
			return "result1", nil
		}
		task2 := func(ctx context.Context) (string, error) {
			time.Sleep(30 * time.Millisecond)
			return "result2", nil
		}
		task3 := func(ctx context.Context) (string, error) {
			time.Sleep(10 * time.Millisecond)
			return "result3", nil
		}

		start := time.Now()
		future := NewParallelWithResult(ctx, task1, task2, task3)
		results, err := future.Wait()
		duration := time.Since(start)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		expectedResults := []string{"result1", "result2", "result3"}
		if len(results) != len(expectedResults) {
			t.Errorf("Expected %d results, got %d", len(expectedResults), len(results))
		}

		for i, expected := range expectedResults {
			if i < len(results) && results[i] != expected {
				t.Errorf("Expected result[%d] = %s, got %s", i, expected, results[i])
			}
		}

		// Should take about 30ms (longest task), not 60ms (sum of all tasks)
		if duration > 100*time.Millisecond {
			t.Errorf("Tasks didn't run in parallel, took %v", duration)
		}
	})

	// Test with one failing task
	t.Run("one failing task", func(t *testing.T) {
		expectedErr := errors.New("task2 failed")
		task1 := func(ctx context.Context) (int, error) {
			time.Sleep(10 * time.Millisecond)
			return 100, nil
		}
		task2 := func(ctx context.Context) (int, error) {
			time.Sleep(20 * time.Millisecond)
			return 0, expectedErr
		}
		task3 := func(ctx context.Context) (int, error) {
			time.Sleep(30 * time.Millisecond)
			return 300, nil
		}

		future := NewParallelWithResult(ctx, task1, task2, task3)
		results, err := future.Wait()

		if err == nil {
			t.Error("Expected error, got nil")
		}
		if err.Error() != expectedErr.Error() {
			t.Errorf("Expected error %v, got %v", expectedErr, err)
		}

		// Results should still be populated for all tasks
		if len(results) != 3 {
			t.Errorf("Expected 3 results even with error, got %d", len(results))
		}
		if len(results) >= 3 {
			// Task 1 should have succeeded
			if results[0] != 100 {
				t.Errorf("Expected results[0] = 100, got %d", results[0])
			}
			// Task 2 failed, should have zero value
			if results[1] != 0 {
				t.Errorf("Expected results[1] = 0 (zero value), got %d", results[1])
			}
			// Task 3 should have succeeded
			if results[2] != 300 {
				t.Errorf("Expected results[2] = 300, got %d", results[2])
			}
		}
	})

	// Test empty task list
	t.Run("empty task list", func(t *testing.T) {
		future := NewParallelWithResult[string](ctx)
		results, err := future.Wait()

		if err != nil {
			t.Errorf("Expected no error for empty task list, got %v", err)
		}
		if len(results) != 0 {
			t.Errorf("Expected empty results slice, got %v", results)
		}
	})

	// Test with different result types
	t.Run("different result types", func(t *testing.T) {
		task1 := func(ctx context.Context) (int, error) {
			return 42, nil
		}
		task2 := func(ctx context.Context) (int, error) {
			return 84, nil
		}

		future := NewParallelWithResult(ctx, task1, task2)
		results, err := future.Wait()

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(results) != 2 {
			t.Errorf("Expected 2 results, got %d", len(results))
		}
		if len(results) >= 2 {
			if results[0] != 42 {
				t.Errorf("Expected results[0] = 42, got %d", results[0])
			}
			if results[1] != 84 {
				t.Errorf("Expected results[1] = 84, got %d", results[1])
			}
		}
	})
}
