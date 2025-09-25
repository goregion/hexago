# Goture - Go Future Library

Goture is a lightweight Future pattern implementation for Go that allows executing tasks asynchronously and waiting for their completion.

## Features

- **Async Task Execution**: Execute tasks asynchronously without blocking the current goroutine
- **Parallel Execution**: Run multiple tasks in parallel and wait for all to complete
- **Result Handling**: Execute tasks that return values using generics
- **Error Handling**: Proper error propagation and panic recovery
- **Context Support**: Full context.Context support for cancellation and timeouts

## Installation

```bash
go get github.com/goregion/hexago/pkg/goture
```

## Usage

### Basic Task Execution

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/goregion/hexago/pkg/goture"
)

func main() {
    ctx := context.Background()
    
    task := func(ctx context.Context) error {
        time.Sleep(1 * time.Second)
        fmt.Println("Task completed!")
        return nil
    }
    
    future := goture.NewGoture(ctx, task)
    err := future.Wait()
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    }
}
```

### Parallel Task Execution

```go
func main() {
    ctx := context.Background()
    
    task1 := func(ctx context.Context) error {
        time.Sleep(100 * time.Millisecond)
        fmt.Println("Task 1 done")
        return nil
    }
    
    task2 := func(ctx context.Context) error {
        time.Sleep(200 * time.Millisecond)
        fmt.Println("Task 2 done")
        return nil
    }
    
    // Execute both tasks in parallel
    future := goture.NewParallelGoture(ctx, task1, task2)
    err := future.Wait() // Waits for both tasks to complete
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    }
}
```

### Tasks with Results

```go
func main() {
    ctx := context.Background()
    
    task := func(ctx context.Context) (string, error) {
        time.Sleep(100 * time.Millisecond)
        return "Hello from async task!", nil
    }
    
    future := goture.NewGotureWithResult(ctx, task)
    result, err := future.Wait()
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    fmt.Printf("Result: %s\n", result)
}
```

### Parallel Tasks with Results

```go
func main() {
    ctx := context.Background()
    
    task1 := func(ctx context.Context) (int, error) {
        time.Sleep(100 * time.Millisecond)
        return 10, nil
    }
    
    task2 := func(ctx context.Context) (int, error) {
        time.Sleep(200 * time.Millisecond)
        return 20, nil
    }
    
    task3 := func(ctx context.Context) (int, error) {
        time.Sleep(50 * time.Millisecond)
        return 30, nil
    }
    
    // Execute all tasks in parallel and collect results
    future := goture.NewParallelWithResult(ctx, task1, task2, task3)
    results, err := future.Wait() // Returns []int{10, 20, 30}
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    fmt.Printf("Results: %v\n", results)
    fmt.Printf("Sum: %d\n", results[0]+results[1]+results[2]) // Sum: 60
}
```

## API Reference

### Types

- `Task`: `func(ctx context.Context) error` - A function that can be executed asynchronously
- `TaskWithResult[T]`: `func(ctx context.Context) (T, error)` - A function that returns a result
- `Goture`: A future representing an async task execution
- `GotureWithResult[T]`: A future representing an async task execution with a result

### Functions

- `NewGoture(ctx context.Context, fn Task) Goture` - Creates a new future for single task
- `NewParallelGoture(ctx context.Context, tasks ...Task) Goture` - Creates a future for parallel tasks
- `NewGotureWithResult[T](ctx context.Context, fn TaskWithResult[T]) GotureWithResult[T]` - Creates a future with result
- `NewParallelWithResult[T](ctx context.Context, tasks ...TaskWithResult[T]) GotureWithResult[[]T]` - Creates a future for parallel tasks with results

### Methods

- `Wait() error` - Waits for task completion and returns any error
- `Wait() (T, error)` - (GotureWithResult) Waits for completion and returns result and error

## Error Handling

- Tasks that panic are automatically recovered and converted to errors
- For parallel execution, the first error encountered is returned
- All tasks continue executing even if one fails (results from successful tasks are still returned)
- Context cancellation is properly propagated
- For `NewParallelWithResult`, partial results are returned even when some tasks fail

## Testing

Run the tests with:

```bash
go test -v ./pkg/goture
```

## License

This library is part of the hexago project.