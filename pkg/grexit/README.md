# Grexit - Graceful Exit Library for Go

Grexit is a simple and lightweight library for handling graceful exits in Go applications by listening for system signals (SIGINT, SIGTERM) and providing context cancellation.

## Features

- ✅ Context-based graceful shutdown
- ✅ Handles SIGINT and SIGTERM signals
- ✅ Configurable shutdown timeout
- ✅ Force shutdown if timeout expires
- ✅ No goroutine leaks
- ✅ Simple and clean API
- ✅ Thread-safe

## Installation

```bash
go get github.com/goregion/hexa/pkg/grexit
```

## Usage

### Basic usage with context

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/goregion/hexa/pkg/grexit"
)

func main() {
    ctx := context.Background()
    
    // Create a context that will be canceled on SIGINT/SIGTERM
    grexitCtx := grexit.WithGrexitContext(ctx)
    
    // Your application logic
    select {
    case <-grexitCtx.Done():
        fmt.Println("Received shutdown signal, exiting gracefully...")
    case <-time.After(30 * time.Second):
        fmt.Println("Application completed normally")
    }
}
```

### Usage with timeout (recommended for production)

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/goregion/hexa/pkg/grexit"
)

func main() {
    ctx := context.Background()
    
    // Create a context with default 30s timeout
    grexitCtx := grexit.WithGrexitTimeout(ctx)
    
    // Or with custom timeout
    grexitCtx = grexit.WithGrexitTimeoutDuration(ctx, 10*time.Second)
    
    // Your application logic
    select {
    case <-grexitCtx.Done():
        fmt.Println("Received shutdown signal, exiting gracefully...")
        // Perform cleanup here...
        // If cleanup takes longer than timeout, the process will be forcefully terminated
    }
}
```

### Advanced usage with cancel function

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/goregion/hexa/pkg/grexit"
)

func main() {
    ctx := context.Background()
    
    // Create a context with cancel function
    grexitCtx, cancel := grexit.WithGrexitCancelContext(ctx)
    defer cancel() // Always call cancel to prevent goroutine leaks
    
    // Start your services
    go func() {
        // Simulate some work
        time.Sleep(5 * time.Second)
        cancel() // Programmatic cancellation
    }()
    
    // Wait for shutdown
    <-grexitCtx.Done()
    fmt.Println("Shutting down...")
}
```

### Advanced usage with timeout and cancel function

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/goregion/hexa/pkg/grexit"
)

func main() {
    ctx := context.Background()
    
    // Create a context with timeout and cancel function
    grexitCtx, cancel := grexit.WithGrexitTimeoutCancel(ctx)
    defer cancel() // Always call cancel to prevent goroutine leaks
    
    // Start your services
    go startHttpServer()
    go startBackgroundWorker()
    
    // Wait for shutdown signal
    <-grexitCtx.Done()
    
    fmt.Println("Shutting down services...")
    
    // Perform graceful shutdown here
    // If this takes longer than the timeout (default 30s),
    // the context will be canceled anyway to force shutdown
    shutdownServices()
    
    fmt.Println("Shutdown complete")
}

func startHttpServer() {
    // Your HTTP server logic
}

func startBackgroundWorker() {
    // Your background worker logic
}

func shutdownServices() {
    // Your cleanup logic
}
```

## API

### Basic Functions

### `WithGrexitContext(ctx context.Context) context.Context`

Returns a context that is canceled when SIGINT or SIGTERM is received.

### `WithGrexitCancelContext(ctx context.Context) (context.Context, context.CancelFunc)`

Returns a context and cancel function. The context is canceled when:
- SIGINT or SIGTERM is received
- The cancel function is called
- The parent context is canceled

**Important**: Always call the returned cancel function to prevent goroutine leaks.

### Timeout Functions

### `WithGrexitTimeout(ctx context.Context) context.Context`

Returns a context that is canceled when SIGINT or SIGTERM is received, with a default 30-second timeout for graceful shutdown.

### `WithGrexitTimeoutDuration(ctx context.Context, timeout time.Duration) context.Context`

Returns a context that is canceled when SIGINT or SIGTERM is received, with a custom timeout for graceful shutdown.

### `WithGrexitTimeoutCancel(ctx context.Context) (context.Context, context.CancelFunc)`

Returns a context and cancel function with default 30-second timeout.

### `WithGrexitTimeoutCancelContext(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc)`

Returns a context and cancel function with custom timeout. The context is canceled when:
- SIGINT or SIGTERM is received (starts the timeout period)
- The timeout expires (forces shutdown)
- The cancel function is called
- The parent context is canceled

**Important**: Always call the returned cancel function to prevent goroutine leaks.

## Best Practices

1. **Use timeout functions in production**: Prefer `WithGrexitTimeout*` functions to avoid hanging processes during shutdown.

2. **Always call cancel()**: When using functions that return a cancel function, always call it, preferably with `defer`.

3. **Choose appropriate timeout**: Default 30s is good for most cases, but adjust based on your application's cleanup requirements.

4. **Handle timeout gracefully**: Design your cleanup code to handle forced shutdown when timeout expires.

5. **Use in main goroutine**: These functions should typically be called in your main goroutine or early in your application startup.

6. **Combine with other contexts**: You can chain these contexts with other context types (timeout, deadline, etc.).

## Constants

### `DefaultShutdownTimeout`

Default timeout for graceful shutdown operations (30 seconds).

## Examples

See the test file for more usage examples.

## Contributing

Feel free to open issues or submit pull requests.

## License

This project is licensed under the MIT License.