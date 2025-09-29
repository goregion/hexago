# Commands

This directory contains all executable commands and entry points for the application. Each subdirectory represents a separate executable that can be built and run.

## What to create here:

### Application Entry Points
- **Main application** - Primary application server (`main.go`)
- **CLI tools** - Command-line utilities and tools
- **Migration tools** - Database migration utilities
- **Seed tools** - Data seeding utilities
- **Background workers** - Standalone worker processes
- **Admin tools** - Administrative utilities

### Command Structure
Each command should be in its own subdirectory with a `main.go` file:

```
cmd/
├── server/         # Main HTTP/gRPC server
│   └── main.go
├── cli/            # Command-line interface
│   └── main.go
├── migrate/        # Database migrations
│   └── main.go
├── worker/         # Background worker
│   └── main.go
├── seed/           # Data seeding
│   └── main.go
└── admin/          # Admin tools
    └── main.go
```

### Example Structure:
```go
// cmd/server/main.go
package main

import (
    "context"
    "log"
    
    "github.com/goregion/hexago/pkg/launcher"
    "github.com/goregion/hexago/internal/adapter"
    "github.com/goregion/hexago/internal/app"
)

func main() {
    ctx := context.Background()
    
    // Initialize dependencies
    app := launcher.NewApp()
    
    // Setup adapters
    httpAdapter := adapter.NewHTTPAdapter(app)
    
    // Start server
    if err := httpAdapter.Start(ctx); err != nil {
        log.Fatal(err)
    }
}
```

### CLI Command Example:
```go
// cmd/cli/main.go
package main

import (
    "flag"
    "fmt"
    "os"
)

func main() {
    var command = flag.String("cmd", "", "Command to execute")
    flag.Parse()
    
    switch *command {
    case "create-user":
        createUser()
    case "list-users":
        listUsers()
    default:
        fmt.Println("Available commands: create-user, list-users")
        os.Exit(1)
    }
}
```

## Key Principles:
- Keep main functions thin - delegate to application layer
- Each command should have a single responsibility
- Use dependency injection to wire up components
- Handle configuration and environment setup here
- Graceful shutdown handling for long-running processes
- Proper error handling and logging setup

## Building Commands:
```bash
# Build specific command
go build -o bin/server ./cmd/server

# Build all commands
make build-all
```