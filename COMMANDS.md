# Available Commands Reference

This document lists all available commands in the Hexago template.

## 📋 Make Commands

Run `make help` to see all available targets.

### Development Setup
```bash
make setup-dev     # Setup development environment (installs tools)
make install       # Install Go dependencies
make deps          # Update dependencies
```

### Code Generation
```bash
make generate      # Generate protobuf code from .proto files
```

### Testing
```bash
make test          # Run all tests
make test-unit     # Run unit tests only
make test-integration  # Run integration tests only
make test-coverage # Run tests with coverage report
make test-race     # Run tests with race detection
```

### Code Quality
```bash
make fmt           # Format Go code
make lint          # Run golangci-lint
make clean         # Clean build artifacts
```

### Building
```bash
make build         # Build all applications
```

### Running Applications
```bash
make run-all       # Run all-in-one application
make run-consumer  # Run binance tick consumer
make run-api       # Run backoffice API  
make run-generator # Run OHLC generator
make dev-run       # Run in development mode
make dev-watch     # Run with hot reload (requires air)
```

### Docker Commands
```bash
make docker-build     # Build Docker images
make docker-up        # Start services with docker-compose
make docker-down      # Stop docker-compose services
make docker-logs      # Show docker logs
make docker-restart   # Restart docker services
make docker-clean     # Clean docker resources
make docker-monitoring # Start with monitoring stack
make docker-shell     # Open shell in container
```

### Template Commands
```bash
make template-init           # Initialize new project from template
make template-create-adapter # Create new adapter boilerplate
make template-create-service # Create new service boilerplate
```

## 🎯 VS Code Tasks

Access via `Ctrl+Shift+P` → "Tasks: Run Task"

### Build Tasks
- **Generate gRPC code from proto** - Generate protobuf code
- **Build Application** - Build all binaries
- **Format Code** - Run gofmt on all files
- **Lint Code** - Run golangci-lint

### Test Tasks
- **Run All Tests** - Execute complete test suite
- **Run Unit Tests** - Execute unit tests only
- **Run Integration Tests** - Execute integration tests only

### Development Tasks
- **Start Development Server** - Run with hot reload
- **Docker: Build** - Build Docker images
- **Docker: Start Services** - Start docker-compose
- **Docker: Stop Services** - Stop docker-compose

### Template Tasks
- **Template: Initialize New Project** - Interactive project setup
- **Template: Create New Adapter** - Generate adapter boilerplate
- **Template: Create New Service** - Generate service boilerplate

## 🐛 Debug Configurations

Available in VS Code Debug panel:

### Application Debugging
- **Debug All-in-One** - Debug main application
- **Debug Binance Consumer** - Debug tick consumer
- **Debug Backoffice API** - Debug gRPC API
- **Debug OHLC Generator** - Debug OHLC service
- **Debug Current File** - Debug currently open Go file

### Testing & Tools
- **Debug Tests** - Debug test files
- **Debug Template Init** - Debug template initialization
- **Attach to Process** - Attach debugger to running process

## 📦 Go Commands

Standard Go commands work as expected:

```bash
# Build specific application
go build ./cmd/all-in-one

# Run with arguments
go run ./cmd/all-in-one -port 8080

# Install tools
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

# Run specific tests
go test ./tests/unit/... -v

# Generate mocks (if mockery is set up)
go generate ./...
```

## 🐳 Docker Commands

Direct Docker commands:

```bash
# Build from Dockerfile
docker build -f docker/Dockerfile -t hexago .

# Run container
docker run -p 8080:8080 -p 9090:9090 hexago

# View logs
docker logs hexago-app

# Execute shell in container
docker exec -it hexago-app sh
```

## 🔧 Development Workflow

### Initial Setup
1. `make setup-dev` - Install development tools
2. `make generate` - Generate protobuf code
3. `make test` - Verify everything works

### Daily Development
1. `make dev-run` or `make dev-watch` - Start development server
2. Make changes to code
3. `make test-unit` - Run fast unit tests
4. `make fmt && make lint` - Format and lint code
5. `make test` - Run all tests before commit

### Docker Development
1. `make docker-up` - Start services (DB, Redis, etc.)
2. `make dev-run` - Run application locally against containers
3. `make docker-logs` - Monitor service logs
4. `make docker-down` - Clean up when done

### Creating New Components
1. `make template-create-adapter` - Generate new adapter
2. `make template-create-service` - Generate new service
3. Follow prompts for customization
4. Add tests and implement logic

### Testing Strategy
- **Unit Tests** - Test business logic in isolation
- **Integration Tests** - Test with real external dependencies
- **End-to-End Tests** - Test complete workflows

Use `make test-unit` for rapid feedback during development, `make test-integration` when changes affect external integrations.

## 🎨 Customization Commands

When using as template:

1. `make template-init` - Initialize new project
2. Update module path in `go.mod`
3. Replace import paths throughout codebase
4. Customize entities and business logic
5. Update documentation and README
6. Configure CI/CD for your repository

## 💡 Tips

- Use `make help` to see all available targets
- VS Code tasks are fastest for development
- Docker commands are great for production-like testing
- Template commands help maintain consistency
- Always run `make test` before committing changes

## 🔍 Troubleshooting

### Common Issues
- **Proto generation fails**: Ensure `protoc` and plugins are installed
- **Tests fail**: Check if external dependencies (Redis, DB) are running
- **Docker build fails**: Verify Dockerfile paths and Go module setup
- **VS Code tasks not working**: Check tasks.json configuration

### Getting Help
- Check documentation in `/docs` directory
- Review examples in the codebase
- Use GitHub Issues for bug reports
- Refer to architecture documentation for design decisions