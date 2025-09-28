# 🏗️ Hexago - Hexagonal Architecture Go Template

[![CI Status](https://github.com/goregion/hexago/workflows/CI%2FCD%20Pipeline/badge.svg)](https://github.com/goregion/hexago/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/goregion/hexago)](https://goreportcard.com/report/github.com/goregion/hexago)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/github/go-mod/go-version/goregion/hexago)](https://golang.org)
[![Release](https://img.shields.io/github/release/goregion/hexago.svg)](https://github.com/goregion/hexago/releases)

**Production-ready Go template** implementing **Hexagonal Architecture** (Ports & Adapters pattern) with comprehensive tooling, testing, and documentation.

> 💡 **Perfect for**: Microservices, APIs, CLI applications, and any Go project requiring clean architecture principles

This template provides a complete foundation for building maintainable, testable, and scalable Go applications following hexagonal architecture patterns with real-world examples and best practices.

## 🎯 Template Features

- ✅ **Complete Hexagonal Architecture** implementation with clear separation of concerns
- ✅ **Code generators** for adapters and services with interactive CLI
- ✅ **Comprehensive testing** strategy (unit + integration + GitHub Actions)
- ✅ **Docker & Docker Compose** support with multi-stage builds
- ✅ **VS Code integration** with tasks, settings, and debugging configuration
- ✅ **GitHub Actions CI/CD** pipeline with testing, linting, and security scanning
- ✅ **Production-ready logging** with structured output and context propagation
- ✅ **gRPC API** with Protocol Buffers and auto-generated documentation
- ✅ **Template initialization** script for creating new projects
- ✅ **Cross-platform support** (Windows, macOS, Linux)
- ✅ **Best practices** and comprehensive documentation

## 🚀 Quick Start

### Option 1: Use GitHub Template (Recommended)

1. **Click "Use this template" button** above or [click here](https://github.com/goregion/hexago/generate)
2. **Create your repository** from the template
3. **Clone your new repository**
4. **Initialize the template** for your project:

```bash
cd your-project
go run ./scripts/template-init
```

### Option 2: Clone and Initialize

```bash
# Clone the repository
git clone https://github.com/goregion/hexago.git my-project
cd my-project

# Initialize template with your project details
go run ./scripts/template-init

# Or manually customize (see Customization section below)
```

## ⚡ Quick Start

📋 **New to Hexago?** Follow our [Quick Start Guide](QUICKSTART.md) for a step-by-step setup process.

## 📋 Available Commands

### 🛠️ Development Workflow
```bash
# Setup development environment
make setup-dev              # Install tools and dependencies
make generate               # Generate protobuf code

# Development server  
make dev-run                # Run with hot reload
make dev-watch              # Run with file watching (requires air)

# Testing
make test                   # Run all tests
make test-unit              # Run unit tests only  
make test-integration       # Run integration tests only
make test-coverage          # Generate coverage report

# Code quality
make fmt                    # Format code
make lint                   # Run linter
make clean                  # Clean build artifacts
```

### 🐳 Docker Development
```bash
make docker-up              # Start services (PostgreSQL, Redis)
make docker-down            # Stop services
make docker-logs            # View logs
make docker-clean           # Clean resources
make docker-monitoring      # Start with Prometheus/Grafana
```

### 🎨 Code Generation
```bash
make template-init          # Initialize new project from template
make template-create-adapter # Generate new adapter boilerplate
make template-create-service # Generate new service boilerplate
```

**💡 See [COMMANDS.md](COMMANDS.md) for complete command reference including VS Code tasks and debug configurations.**

## ⚙️ Prerequisites

- **Go 1.23+** - Latest version recommended
- **Protocol Buffers compiler** (`protoc`) - For gRPC code generation  
- **Docker & Docker Compose** - For local development environment
- **Make** - For build automation (Windows users: install via chocolatey or use PowerShell alternatives)
- **VS Code** (recommended) - With Go extension for optimal development experience

### Development Commands

```bash
# Setup development environment (installs tools and dependencies)
make setup-dev

# Install dependencies
make install

# Generate code (protobuf, mocks, etc.)
make generate

# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run tests with race detection
make test-race

# Format code
make fmt

# Run linter
make lint

# Build applications
make build

# Clean artifacts
make clean
```

## 🎨 Using the Template

### 🔧 For Existing Projects (Manual Customization)

1. **Update Module Path**: Change `github.com/goregion/hexago` in `go.mod`
2. **Replace Imports**: Update all import statements throughout the codebase
3. **Customize Entities**: Modify entities in `internal/entity/` for your domain
4. **Define Ports**: Create interfaces in `internal/port/` for your use cases
5. **Implement Services**: Add business logic in `internal/service/`
6. **Build Adapters**: Create external integrations in `internal/adapter/`
7. **Wire Dependencies**: Update dependency injection in `internal/app/`

### 🎯 Architecture Overview

This template demonstrates **Hexagonal Architecture** with a real cryptocurrency trading system:

![Hexagonal Architecture](./docs/hexagonal.png)

#### Key Components:

- **🏛️ Domain Core** (`internal/entity/`, `internal/service/`) - Pure business logic
- **🔌 Ports** (`internal/port/`) - Interfaces defining contracts
- **🔧 Adapters** (`internal/adapter/`) - External system integrations  
- **🎯 Application Services** (`internal/app/`) - Workflow orchestration
- **� API Layer** (`cmd/`, `api/`) - External interfaces (gRPC, HTTP, CLI)

## 📚 Documentation

- **[🚀 Quick Start Guide](QUICKSTART.md)** - Get up and running in minutes
- **[🏗️ Architecture Guide](docs/ARCHITECTURE.md)** - Detailed hexagonal architecture explanation
- **[💻 Development Guide](docs/DEVELOPMENT.md)** - Setup and development workflows  
- **[📋 Commands Reference](COMMANDS.md)** - All available make commands and VS Code tasks
- **[🤝 Contributing Guide](CONTRIBUTING.md)** - How to contribute to this project
- **[🔐 Security Policy](SECURITY.md)** - Security guidelines and vulnerability reporting
- **[📝 Changelog](CHANGELOG.md)** - Version history and release notes

## 🎯 Real-World Example

This template includes a working **cryptocurrency trading system** as an example implementation:

### 🏗️ **Domain**: OHLC (Open, High, Low, Close) Data Processing
- **Entities**: `Tick`, `OHLC`, `LiquidityProviderTick`
- **Services**: Price aggregation, OHLC generation, data validation
- **Adapters**: Binance WebSocket, PostgreSQL, Redis, gRPC API

### 🔄 **Data Flow**:
```
Binance WebSocket → Tick Consumer → Redis → OHLC Generator → Database → gRPC API
```

### 🎯 **Architectural Benefits Demonstrated**:
- **Business Logic Isolation**: Core OHLC logic independent of data sources
- **Easy Testing**: Mock adapters for reliable unit tests
- **Flexible Integration**: Add new exchanges without changing core logic
- **Technology Independence**: Swap databases or APIs without core changes

## 📁 Project Structure

```
hexago/
├── 📁 .github/              # GitHub configuration and workflows
│   ├── workflows/           # CI/CD pipelines (ci.yml, release.yml)
│   ├── ISSUE_TEMPLATE/      # Bug report and feature request templates
│   ├── CODEOWNERS           # Code review assignments
│   ├── FUNDING.yml          # Sponsorship configuration
│   └── pull_request_template.md
├── 📁 .vscode/              # VS Code development environment
│   ├── extensions.json      # Recommended extensions
│   ├── launch.json          # Debug configurations
│   ├── settings.json        # Editor settings
│   └── tasks.json           # Build and development tasks
├── 📁 api/                  # API definitions and contracts
│   └── backoffice/grpc/     # gRPC service definitions (.proto files)
├── 📁 cmd/                  # Application entry points
│   ├── all-in-one/          # Main application combining all services
│   ├── backoffice-api/      # gRPC API server
│   ├── binance-tick-consumer/ # External data consumer
│   └── ohlc-generator/      # OHLC data generation service
├── 📁 docker/               # Container configurations
│   └── Dockerfile           # Multi-stage production container
├── 📁 docs/                 # Comprehensive documentation
│   ├── ARCHITECTURE.md      # Hexagonal architecture deep dive
│   ├── DEVELOPMENT.md       # Development workflows and guidelines
│   └── hexagonal.png        # Architecture diagram
├── 📁 internal/             # Private application code (core business logic)
│   ├── adapter/            # External system integrations
│   │   ├── binance/        # Binance WebSocket adapter
│   │   ├── grpc/           # gRPC server adapter
│   │   ├── mysql/          # Database adapter
│   │   └── redis/          # Redis adapter
│   ├── app/                # Application services and orchestration
│   ├── entity/             # Domain entities and business objects
│   ├── port/               # Port interfaces (contracts)
│   └── service/            # Domain service implementations
├── 📁 pkg/                  # Public libraries and utilities
│   ├── config/             # Configuration management
│   ├── launcher/           # Application launcher utilities
│   ├── log/                # Structured logging
│   ├── redis/              # Redis client utilities
│   └── sqlgen-db/          # Database utilities
├── 📁 scripts/              # Development and build scripts
│   ├── create-adapter/     # Interactive adapter generator
│   ├── create-service/     # Interactive service generator
│   └── template-init/      # Project initialization script
├── 📁 tests/                # Test suites
│   ├── integration/        # Integration tests with external dependencies
│   └── unit/               # Fast, isolated unit tests
├── 🔧 .air.toml            # Hot reload configuration
├── � .dockerignore        # Docker build context exclusions
├── 📄 .env.example         # Environment variables template
├── 📄 .env.test            # Test environment configuration
├── 🔒 .gitignore           # Git exclusions
├── 📋 CHANGELOG.md         # Version history and release notes
├── 📚 COMMANDS.md          # Complete command reference
├── 🤝 CONTRIBUTING.md      # Contribution guidelines
├── 🚀 docker-compose.yml   # Development environment services
├── 📜 go.mod              # Go module definition
├── 🏗️  Makefile           # Build automation and development tasks
├── ⚡ QUICKSTART.md       # Step-by-step setup guide
└── 🔐 SECURITY.md         # Security policy and reporting
```

## 🛠️ What's Included

### 🏗️ Architecture & Structure
- **Hexagonal Architecture** implementation with clear separation of concerns
- **Domain-Driven Design** patterns and best practices
- **Dependency Injection** through interfaces and ports
- **Clean Architecture** principles throughout

### 🧪 Testing Framework
- **Unit tests** with mocking and isolation
- **Integration tests** with real external dependencies
- **GitHub Actions CI/CD** with automated testing
- **Test coverage** reporting and enforcement
- **Race condition detection** in tests

### 🐳 Containerization
- **Multi-stage Dockerfile** for optimized images
- **Docker Compose** for development environment
- **PostgreSQL** and **Redis** services included
- **Health checks** and proper container lifecycle
- **Non-root user** for security

### 🔧 Development Tools
- **VS Code integration** with tasks, debugging, and settings
- **Hot reload** development server with Air
- **Protocol Buffers** code generation
- **Linting** with golangci-lint
- **Code formatting** with gofmt/gofumpt
- **Dependency management** with Go modules

### 🚀 CI/CD & Automation
- **GitHub Actions** workflows for testing and deployment
- **Release automation** with changelog generation
- **Security scanning** with Gosec
- **Docker image** building and publishing
- **Multi-platform** support (Linux, Windows, macOS)

### 📚 Documentation & Templates
- **Comprehensive documentation** with examples
- **Issue templates** for bug reports and features
- **Pull request templates** with checklists
- **Contributing guidelines** and code of conduct
- **Security policy** and vulnerability reporting

### 🎨 Code Generation
- **Template initialization** script for new projects
- **Adapter generator** for creating new external integrations
- **Service generator** for domain services
- **Interactive CLI** for customization

### Why Hexagonal Architecture?

- ✅ **Technology Independence** - Business logic not tied to specific databases or frameworks
- ✅ **Easy Testing** - Mock external dependencies for fast, reliable tests  
- ✅ **Flexible Integration** - Add HTTP, gRPC, CLI interfaces without changing core logic
- ✅ **Future-Proof** - Easy to evolve and adapt to new requirements
- ✅ **Clear Boundaries** - Well-defined contracts between layers
- ✅ **Dependency Inversion** - External layers depend on internal layers, not vice versa

## 🧪 Testing Strategy

### Three-Layer Testing Approach

1. **🚀 Unit Tests** (`tests/unit/`)
   - Test business logic in isolation
   - Fast execution (< 1ms per test)  
   - Mock all external dependencies
   - 90%+ coverage target

2. **🔗 Integration Tests** (`tests/integration/`)
   - Test with real external systems (PostgreSQL, Redis)
   - Use Docker containers for consistency
   - Test complete data flows
   - Verify adapter implementations

3. **🌐 End-to-End Tests**  
   - Test complete user scenarios
   - gRPC API testing with real services
   - Performance and load testing
   - Production-like environment

### Running Tests

```bash
make test              # All tests with optimal order
make test-unit         # Fast unit tests only
make test-integration  # Integration tests with containers  
make test-coverage     # Generate HTML coverage report
make test-race         # Race condition detection
```

## 🏗️ Development Workflow

### Daily Development Loop

```bash
# 1. Start development environment
make docker-up                 # External services (DB, Redis)
make dev-run                   # Application with hot reload

# 2. Make changes and verify
make test-unit                 # Quick feedback loop
make fmt && make lint          # Code quality checks

# 3. Integration testing
make test-integration          # Test with real dependencies

# 4. Final verification
make test                      # Complete test suite
make build                     # Ensure everything compiles
```

### 🔧 VS Code Integration

This template includes comprehensive VS Code support:

- **🎯 Tasks** (Ctrl+Shift+P → "Tasks: Run Task")
  - Generate gRPC code from proto
  - Run tests (unit/integration/all)
  - Build and run applications
  - Docker operations
  - Template generators

- **🐛 Debug Configurations**
  - Debug all applications individually
  - Attach to running processes
  - Debug current file
  - Debug tests

- **⚙️ Settings & Extensions**
  - Go language server optimized settings
  - Recommended extensions auto-install
  - Format on save and organize imports
  - Integrated terminal configurations

## 🚀 Production Deployment

### Docker Production Build

```bash
# Build optimized production image
docker build -f docker/Dockerfile -t hexago:latest .

# Run with production configuration
docker run -d \
  --name hexago-prod \
  -p 8080:8080 \
  -p 9090:9090 \
  -e APP_ENV=production \
  -e DB_HOST=your-db-host \
  -e REDIS_HOST=your-redis-host \
  hexago:latest
```

### 🔄 CI/CD Pipeline

Automated workflows included:

- **✅ Continuous Integration** (`.github/workflows/ci.yml`)
  - Multi-Go version testing (1.23+)
  - Cross-platform builds (Linux, Windows, macOS)
  - Security scanning with Gosec
  - Dependency vulnerability checks
  - Code quality with golangci-lint

- **🚀 Automated Releases** (`.github/workflows/release.yml`)  
  - Semantic version tagging
  - Docker image builds (multi-arch)
  - GitHub releases with changelogs
  - Container registry publishing

## 🔧 Customization Guide

When using this template for your project:

### 1. **Initialize Your Project**
```bash
go run ./scripts/template-init    # Interactive setup wizard
```

### 2. **Define Your Domain**
- Update entities in `internal/entity/`
- Define business rules and validation
- Create domain services in `internal/service/`

### 3. **Design Your Ports**
- Define input ports (use cases) in `internal/port/`
- Define output ports (external dependencies) 
- Keep interfaces focused and cohesive

### 4. **Build Your Adapters**
- Implement external system integrations
- Database adapters in `internal/adapter/`
- API adapters for external services
- Event publishers and consumers

### 5. **Wire Everything Together**
- Update dependency injection in `internal/app/`
- Configure application services
- Set up proper error handling and logging

## ❓ FAQ

### **Q: How do I add a new external service integration?**
A: Use `make template-create-adapter` to generate boilerplate, then implement the port interface.

### **Q: Can I use this template for HTTP REST APIs instead of gRPC?**  
A: Absolutely! The core architecture is protocol-agnostic. Just replace the gRPC adapter with an HTTP adapter.

### **Q: How do I handle database migrations?**
A: Add migration scripts in a `migrations/` directory and run them in your deployment pipeline.

### **Q: Is this template suitable for microservices?**
A: Yes! Each service can be deployed independently, and the clear boundaries make service extraction easy.

### **Q: How do I add authentication/authorization?**
A: Implement it as a middleware in your adapter layer or as a cross-cutting concern in the application layer.

### **Q: Can I use different databases for different aggregates?**
A: Yes! Create separate repository adapters for each database in `internal/adapter/`.

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details on:

- 🐛 **Reporting bugs** and suggesting features
- 💻 **Code contributions** and development setup  
- 📝 **Documentation** improvements
- 🧪 **Testing** requirements
- 🎯 **Architecture** guidelines

### Contributors

<a href="https://github.com/goregion/hexago/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=goregion/hexago" />
</a>

## 📞 Support & Community

- 📖 **Documentation**: Comprehensive guides in [docs/](docs/) directory
- 🐛 **Bug Reports**: [Create an issue](https://github.com/goregion/hexago/issues/new?template=bug_report.yml)
- 💡 **Feature Requests**: [Create an issue](https://github.com/goregion/hexago/issues/new?template=feature_request.yml)  
- 💬 **Discussions**: [GitHub Discussions](https://github.com/goregion/hexago/discussions)
- 🆘 **Help**: Tag `@goregion` in issues or discussions

## 📈 Roadmap

- [ ] **GraphQL adapter** example implementation
- [ ] **Event sourcing** pattern examples
- [ ] **CQRS** implementation with separate read/write models
- [ ] **Kubernetes** deployment manifests
- [ ] **Observability** with OpenTelemetry integration
- [ ] **Message queue** adapters (RabbitMQ, Apache Kafka)

## 📄 License

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.

**Free to use for commercial and open-source projects!** 🎉

---

<div align="center">

**🚀 Built with ❤️ by the GoRegion team**

⭐ **Star this repo** if you find it helpful! ⭐

[Get Started](QUICKSTART.md) • [View Docs](docs/) • [Report Bug](https://github.com/goregion/hexago/issues) • [Request Feature](https://github.com/goregion/hexago/issues) • [Contribute](CONTRIBUTING.md)

**Hexagonal Architecture** • **Clean Code** • **Production Ready** • **Fully Tested** • **Well Documented**

</div>