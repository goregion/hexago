package goture_test

import (
	"context"
	"fmt"
	"time"

	"github.com/goregion/hexago/pkg/goture"
)

// Example demonstrates basic usage of Goture
func Example() {
	ctx := context.Background()

	// Simple task execution
	task := func(ctx context.Context) error {
		fmt.Println("Task executed")
		return nil
	}

	future := goture.NewGoture(ctx, task)
	err := future.Wait()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	// Output: Task executed
}

// ExampleNewParallelGoture demonstrates parallel task execution
func ExampleNewParallelGoture() {
	ctx := context.Background()

	task1 := func(ctx context.Context) error {
		time.Sleep(10 * time.Millisecond)
		fmt.Println("Task 1 completed")
		return nil
	}

	task2 := func(ctx context.Context) error {
		time.Sleep(20 * time.Millisecond)
		fmt.Println("Task 2 completed")
		return nil
	}

	task3 := func(ctx context.Context) error {
		time.Sleep(5 * time.Millisecond)
		fmt.Println("Task 3 completed")
		return nil
	}

	// Execute all tasks in parallel
	future := goture.NewParallelGoture(ctx, task1, task2, task3)
	err := future.Wait()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	fmt.Println("All tasks completed")

	// Unordered output:
	// Task 3 completed
	// Task 1 completed
	// Task 2 completed
	// All tasks completed
}

// ExampleNewGotureWithResult demonstrates task execution with result
func ExampleNewGotureWithResult() {
	ctx := context.Background()

	task := func(ctx context.Context) (string, error) {
		time.Sleep(10 * time.Millisecond)
		return "Hello, Goture!", nil
	}

	future := goture.NewGotureWithResult(ctx, task)
	result, err := future.Wait()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Result: %s\n", result)
	// Output: Result: Hello, Goture!
}

// ExampleNewParallelWithResult demonstrates parallel task execution with results
func ExampleNewParallelWithResult() {
	ctx := context.Background()

	// Tasks that return different results
	task1 := func(ctx context.Context) (int, error) {
		time.Sleep(20 * time.Millisecond)
		return 10, nil
	}

	task2 := func(ctx context.Context) (int, error) {
		time.Sleep(10 * time.Millisecond)
		return 20, nil
	}

	task3 := func(ctx context.Context) (int, error) {
		time.Sleep(30 * time.Millisecond)
		return 30, nil
	}

	// Execute all tasks in parallel and collect results
	future := goture.NewParallelWithResult(ctx, task1, task2, task3)
	results, err := future.Wait()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Results: %v\n", results)
	fmt.Printf("Sum: %d\n", results[0]+results[1]+results[2])
	// Output: Results: [10 20 30]
	// Sum: 60
}
