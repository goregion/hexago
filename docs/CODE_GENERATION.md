# Code Generation Tools

This project includes cross-platform Go-based tools for generating boilerplate code. These tools help maintain consistency and speed up development when working with hexagonal architecture.

## Available Tools

### 1. Template Initialization (`scripts/template-init/`)

Initializes a new project from this template.

```bash
# Run from template directory
go run ./scripts/template-init

# Or via Makefile
make template-init
```

**What it does:**
- Prompts for new project details (name, module path, description)
- Updates `go.mod` with new module path
- Replaces import paths throughout the codebase
- Updates README.md with project-specific information
- Removes template-specific files
- Initializes git repository

### 2. Adapter Generator (`scripts/create-adapter/`)

Creates a new adapter with the proper hexagonal architecture structure.

```bash
# Run from project directory
go run ./scripts/create-adapter

# Or via Makefile
make template-create-adapter
```

**What it does:**
- Prompts for adapter details (name, type, port interface)
- Creates adapter implementation in `internal/adapter/{name}/`
- Generates configuration structure
- Creates unit tests
- Generates port interface if it doesn't exist
- Provides guidance for next steps

**Generated files:**
```
internal/adapter/{name}/
├── adapter.go        # Main adapter implementation
├── config.go         # Configuration structure
└── adapter_test.go   # Unit tests

internal/port/
└── {name}.go         # Port interface (if new)
```

### 3. Service Generator (`scripts/create-service/`)

Creates domain and application services following hexagonal architecture patterns.

```bash
# Run from project directory
go run ./scripts/create-service

# Or via Makefile
make template-create-service
```

**What it does:**
- Prompts for service details (name, type)
- Creates domain service in `internal/service/{name}/`
- Creates application service in `internal/app/{name}/`
- Generates test files for both services
- Creates service port interface
- Provides architectural guidance

**Generated files:**
```
internal/service/{name}/
├── {name}.go         # Domain service
└── {name}_test.go    # Domain service tests

internal/app/{name}/
├── service.go        # Application service
└── service_test.go   # Application service tests

internal/port/
└── {name}-service.go # Service interface (if new)
```

## Why Go-based Tools?

### Cross-platform Compatibility
- **Works everywhere**: Windows, Linux, macOS without modification
- **No shell dependencies**: No need for bash, PowerShell, or specific shell features
- **Consistent behavior**: Same functionality across all platforms

### Maintainability
- **Type safety**: Compile-time error checking
- **IDE support**: Full Go IDE features (autocomplete, refactoring, debugging)
- **Testing**: Can write unit tests for the generation logic
- **Refactoring**: Easy to modify and extend functionality

### Integration
- **Module awareness**: Automatic `go.mod` parsing and module path handling
- **Import management**: Proper handling of Go import paths
- **Project structure**: Deep understanding of Go project conventions
- **Build integration**: Seamless integration with `go run`, Makefile, and CI/CD

## Tool Architecture

Each tool follows a similar structure:

```go
func main() {
    // 1. Collect user input
    data := getUserInput()
    
    // 2. Validate input
    if err := validate(data); err != nil {
        handleError(err)
    }
    
    // 3. Generate files from templates
    if err := generateFiles(data); err != nil {
        handleError(err)
    }
    
    // 4. Provide user guidance
    showNextSteps(data)
}
```

### Key Components:

- **User Input**: Interactive prompts with defaults
- **Template Processing**: Go `text/template` for file generation
- **File System Operations**: Cross-platform file/directory creation
- **Module Path Detection**: Automatic `go.mod` parsing
- **Color Output**: Informative console messages with colors
- **Error Handling**: Clear error messages and validation

## Customization

### Adding New Generators

1. **Create new script directory**: `scripts/create-{component}/`
2. **Implement main.go**: Follow existing patterns with embedded templates
3. **Update Makefile**: Add new target for easy access
4. **Document**: Update this file with new tool information

### Modifying Templates

Templates are embedded directly in the Go code using string literals. This approach:
- **Eliminates file dependencies**: No external template files to manage
- **Improves portability**: Everything is self-contained
- **Enables customization**: Easy to modify templates per project needs
- **Better version control**: Templates are versioned with the code

To modify a template:
1. Find the template string in the tool's source code
2. Update the template content
3. Test the generation with `go run`

## Best Practices

### When to Use Each Tool

- **Template Init**: Only once per project, when creating from template
- **Adapter Generator**: For each external integration (database, API, queue, etc.)
- **Service Generator**: For each business domain or major feature area

### Naming Conventions

- **Adapters**: Use descriptive names (`postgres`, `redis`, `stripe-api`)
- **Services**: Use domain names (`user`, `payment`, `notification`)
- **Interfaces**: Use clear, actionable names (`UserRepository`, `PaymentProcessor`)

### After Generation

1. **Implement interfaces**: Fill in the TODO comments
2. **Add dependencies**: Update constructors with required dependencies
3. **Write tests**: Expand the generated test stubs
4. **Wire up**: Add to dependency injection/application wiring
5. **Document**: Add specific documentation for your implementation

## Examples

### Creating a PostgreSQL Adapter

```bash
$ go run ./scripts/create-adapter
🔌 Creating New Adapter

Adapter name (e.g., postgres, kafka, http): postgres
Adapter type [driven]: 
Port interface name (e.g., UserRepository, EventPublisher): UserRepository

[INFO] Creating adapter: postgres
[INFO] Adapter struct: PostgresAdapter
[INFO] Port interface: UserRepository
[INFO] Directory: internal/adapter/postgres

✅ Adapter 'postgres' created successfully!
```

### Creating a User Service

```bash
$ go run ./scripts/create-service
⚙️ Creating New Service

Service name (e.g., user, notification, payment): user
Service type [domain]: 

[INFO] Creating service: user
[INFO] Domain service: UserService
[INFO] App service: UserApplicationService
[INFO] Service dir: internal/service/user
[INFO] App service dir: internal/app/user

✅ Service 'user' created successfully!
```

These tools help maintain consistency and reduce boilerplate while following hexagonal architecture principles.