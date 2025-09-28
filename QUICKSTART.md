# Quick Start Guide

Get up and running with Hexago in minutes!

## 🚀 1. Create Your Project

### Option A: Use GitHub Template (Recommended)
1. Click **"Use this template"** on the [GitHub repository](https://github.com/goregion/hexago)
2. Create your new repository
3. Clone it locally

### Option B: Clone and Setup
```bash
git clone https://github.com/goregion/hexago.git my-awesome-project
cd my-awesome-project
```

## 🔧 2. Initialize Template

Run the interactive setup:
```bash
go run ./scripts/template-init
```

This will:
- ✅ Ask for your project details
- ✅ Update module paths and imports  
- ✅ Customize configuration files
- ✅ Set up your project structure

## 🛠️ 3. Setup Development Environment

```bash
# Install development tools and dependencies
make setup-dev

# Generate protobuf code
make generate

# Verify everything works
make test
```

## 🐳 4. Start Development Environment

```bash
# Start external services (PostgreSQL, Redis)
make docker-up

# Run the application in development mode  
make dev-run
```

Your application will be available at:
- 🌐 **HTTP API**: http://localhost:8080
- 🔌 **gRPC API**: localhost:9090
- 📊 **Health Check**: http://localhost:8080/health

## 🧪 5. Verify Installation

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Check code quality
make lint
```

## 🎯 6. Start Building

### Create Your First Adapter
```bash
make template-create-adapter
# Follow prompts to generate boilerplate
```

### Create Your First Service
```bash
make template-create-service
# Follow prompts to generate domain service
```

### Modify Business Logic
1. Define your domain entities in `internal/entity/`
2. Create service interfaces in `internal/port/`
3. Implement services in `internal/service/`
4. Build adapters in `internal/adapter/`

## 📖 Next Steps

- 📚 Read the [Architecture Guide](docs/ARCHITECTURE.md)
- 🔧 Check available [Commands](COMMANDS.md)
- 🧪 Learn about [Testing Strategy](docs/DEVELOPMENT.md)
- 🐳 Explore [Docker Development](docker-compose.yml)

## 🆘 Need Help?

- 📖 Check the [Documentation](docs/)
- 🐛 [Report Issues](https://github.com/goregion/hexago/issues)
- 💬 [Start Discussion](https://github.com/goregion/hexago/discussions)
- 🤝 [Contributing Guide](CONTRIBUTING.md)

---

**🎉 Congratulations!** You now have a production-ready Go application following hexagonal architecture principles.

Happy coding! 🚀