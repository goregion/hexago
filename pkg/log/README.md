# Log Package

A wrapper around Go's standard `log/slog` library, providing convenient logging methods, support for multiple handlers, and context integration.

## Features

- 🚀 Simple API based on `log/slog`
- 🔄 Support for multiple handlers simultaneously
- 🎯 Integration with `context.Context`
- ⚙️ Flexible configuration via environment variables
- 📝 Convenient methods for various logging scenarios
- 🧪 Full test coverage

## Quick Start

```go
package main

import (
    "your-project/pkg/log"
)

func main() {
    // Create a simple logger
    logger := log.NewLogger(log.NewJsonStdOutHandler())
    
    logger.Info("Application started", "version", "1.0.0")
    logger.Error("Something went wrong", "error", err)
}
```

## Configuration

### Environment Variables

- `ENABLE_DEBUG_LOG_LEVEL=true` - enables DEBUG log level (default is INFO)

### Creating a logger

```go
// With default handler (JSON stdout)
logger := log.NewLogger()

// With specific handlers
logger := log.NewLogger(
    log.NewJsonStdOutHandler(),
    log.NewTextStdErrHandler(),
)

// With custom writer
logger := log.NewLogger(
    log.NewJsonHandler(customWriter),
)
```

## Available Handlers

- `NewTextStdOutHandler()` - text output to stdout
- `NewJsonStdOutHandler()` - JSON output to stdout  
- `NewTextStdErrHandler()` - text output to stderr
- `NewJsonStdErrHandler()` - JSON output to stderr
- `NewTextHandler(io.Writer)` - text output to custom writer
- `NewJsonHandler(io.Writer)` - JSON output to custom writer
- `NewHandlerWithLevel(handler, writer, level)` - handler with custom level

## Context Usage

```go
// Add logger to context
ctx := log.WithLoggerContext(context.Background(), logger)

// Get logger from context
logger := log.MustGetLoggerFromContext(ctx)

// Safe get with error check
logger, err := log.GetLoggerFromContext(ctx)
if err != nil {
    // handle error
}
```

## Service Lifecycle

```go
logger := log.NewLogger(log.NewTextStdOutHandler())

// Automatic logging of service start and stop
serviceLogger, stop := logger.StartService("user-service")
defer stop() // Will log "stop" on completion

serviceLogger.Info("Processing request", "user_id", 123)
```

## Convenience Methods

### Logging with pre-defined fields

```go
userLogger := logger.WithFields(map[string]any{
    "user_id":   12345,
    "tenant_id": "tenant-abc",
})

userLogger.Info("User logged in") // Will automatically include user_id and tenant_id
```

### Error logging

```go
// Logs only if error is not nil and not context.Canceled
logger.LogIfError(err, "Failed to process request", "retry_count", 3)

// Create logger with pre-defined error
errorLogger := logger.WithError(err)
errorLogger.Error("Operation failed")
```

### Logging with custom level

```go
logger.LogWithLevel(slog.LevelWarn, "Custom warning", "details", "something")
```

## Architecture

The package consists of the following components:

- **logger.go** - main `Logger` type and its methods
- **handlers.go** - factory functions for creating handlers
- **multi-handler.go** - `slog.Handler` implementation for multiple handlers support
- **context.go** - utilities for working with context

## Testing

```bash
go test ./pkg/log/...
go test -bench=. ./pkg/log/...
```

## Examples

More usage examples can be found in the `examples_test.go` file.