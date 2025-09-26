# Development Guide

## 🚀 Getting Started

### Prerequisites

- **Go 1.21+** - [Install Go](https://golang.org/dl/)
- **Protocol Buffers** - [Install protoc](https://grpc.io/docs/protoc-installation/)
- **Docker & Docker Compose** - [Install Docker](https://docs.docker.com/get-docker/)
- **Make** - For running automated tasks

### Initial Setup

1. **Clone the repository**
```bash
git clone <your-repo-url>
cd hexago
```

2. **Setup development environment**
```bash
make setup-dev
```

3. **Install dependencies**
```bash
make install
```

4. **Generate code**
```bash
make generate
```

5. **Run tests**
```bash
make test
```

## 🏗️ Project Structure

### Core Architecture Files

```
internal/
├── entity/           # 🏛️ Domain entities (business objects)
├── port/             # 🔌 Interface contracts
├── service/          # ⚙️ Business logic services  
├── adapter/          # 🔄 External integrations
└── app/              # 🎭 Application orchestration
```

### Supporting Files

```
api/                  # 📡 API definitions (protobuf)
cmd/                  # 🚀 Application entry points
pkg/                  # 📦 Shared libraries  
tests/                # 🧪 Test suites
docker/               # 🐳 Docker configurations
scripts/              # 🔧 Automation scripts
docs/                 # 📚 Documentation
```

## 🔧 Development Workflow

### Adding New Features

1. **Define Domain Entity** (if needed)
```bash
# Create or update entity
internal/entity/new_entity.go
```

2. **Create Port Interfaces**
```bash
# Define contracts
internal/port/new_entity.go
```

3. **Implement Business Logic**
```bash
# Create service
internal/service/new_entity/service.go
```

4. **Create Adapters**
```bash
# Use helper script
make template-create-adapter

# Or manually create
internal/adapter/new_adapter/
```

5. **Wire Everything Up**
```bash
# Update application service
internal/app/your-app/service.go
```

6. **Add Tests**
```bash
# Unit tests
internal/service/new_entity/service_test.go

# Integration tests  
tests/integration/new_feature_test.go
```

### Using Generation Scripts

#### Create New Adapter
```bash
./scripts/create-adapter.sh
# Follow prompts to generate adapter boilerplate
```

#### Create New Service
```bash
./scripts/create-service.sh  
# Follow prompts to generate service boilerplate
```

#### Initialize New Project
```bash
./scripts/template-init.sh
# Creates new project from this template
```

## 🧪 Testing Strategy

### Test Pyramid

```
                    /\
                   /  \
               E2E /    \ (Few)
                  /______\
                 /        \
        Integration/        \ (Some)
               /____________\
              /              \
            Unit/              \ (Many)
           /__________________\
```

### Unit Tests
- Test business logic in isolation
- Use mocks for dependencies
- Fast execution (< 100ms per test)

```bash
# Run unit tests
make test-unit

# Run with coverage
make test-coverage
```

### Integration Tests  
- Test with real adapters
- Use test databases/services
- Slower but more comprehensive

```bash
# Run integration tests
make test-integration
```

### End-to-End Tests
- Test complete user journeys
- Use real external services (or realistic fakes)
- Slowest but most confidence

```bash
# Run all tests
make test
```

## 📝 Code Style Guide

### Naming Conventions

#### Packages
- Use lowercase, single word names
- Use domain names: `user`, `order`, `payment`

#### Types
- Use PascalCase: `UserService`, `OrderRepository`
- Interfaces should be verbs or agent nouns: `UserCreator`, `EmailSender`

#### Variables
- Use camelCase: `userID`, `emailAddress`
- Be descriptive: `userRepository`, not `ur`

#### Constants
- Use PascalCase for exported: `DefaultTimeout`
- Use camelCase for unexported: `defaultRetries`

### Interface Design

#### Small and Focused
```go
// ✅ Good - focused interface
type UserCreator interface {
    CreateUser(ctx context.Context, user *entity.User) error
}

// ❌ Bad - too broad
type UserManager interface {
    CreateUser(ctx context.Context, user *entity.User) error
    UpdateUser(ctx context.Context, user *entity.User) error
    DeleteUser(ctx context.Context, id string) error
    FindUser(ctx context.Context, id string) (*entity.User, error)
    ListUsers(ctx context.Context) ([]*entity.User, error)
    SendEmail(ctx context.Context, user *entity.User) error
}
```

#### Define at Point of Use
```go
// ✅ Good - interface in same package as consumer
package service

type UserRepository interface {
    Save(ctx context.Context, user *entity.User) error
}

type UserService struct {
    repository UserRepository
}
```

### Error Handling

#### Create Domain Errors
```go
// internal/entity/errors.go
var (
    ErrUserNotFound     = errors.New("user not found")
    ErrUserAlreadyExists = errors.New("user already exists")
    ErrInvalidEmail     = errors.New("invalid email address")
)
```

#### Wrap Context
```go
func (s *UserService) CreateUser(ctx context.Context, req CreateUserRequest) error {
    user := &entity.User{Email: req.Email}
    
    if err := s.repository.Save(ctx, user); err != nil {
        return fmt.Errorf("failed to save user %s: %w", user.Email, err)
    }
    
    return nil
}
```

### Context Usage

#### Always Pass Context
```go
// ✅ Good
func (s *UserService) GetUser(ctx context.Context, id string) (*entity.User, error)

// ❌ Bad  
func (s *UserService) GetUser(id string) (*entity.User, error)
```

#### Check Context Cancellation
```go
func (s *UserService) ProcessUsers(ctx context.Context, users []*entity.User) error {
    for _, user := range users {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
            if err := s.processUser(ctx, user); err != nil {
                return err
            }
        }
    }
    return nil
}
```

## 🔄 Common Patterns

### Repository Pattern
```go
type UserRepository interface {
    Save(ctx context.Context, user *entity.User) error
    FindByID(ctx context.Context, id string) (*entity.User, error)
    FindByEmail(ctx context.Context, email string) (*entity.User, error)
}

// Implementation
type postgresUserRepository struct {
    db *sql.DB
}

func (r *postgresUserRepository) Save(ctx context.Context, user *entity.User) error {
    query := `INSERT INTO users (id, email, name) VALUES ($1, $2, $3)`
    _, err := r.db.ExecContext(ctx, query, user.ID, user.Email, user.Name)
    return err
}
```

### Service Pattern
```go
type UserService struct {
    repository UserRepository
    publisher  EventPublisher
}

func (s *UserService) CreateUser(ctx context.Context, req CreateUserRequest) (*entity.User, error) {
    // Validate
    if err := validateCreateUserRequest(req); err != nil {
        return nil, err
    }
    
    // Create entity
    user := &entity.User{
        ID:    generateID(),
        Email: req.Email,
        Name:  req.Name,
    }
    
    // Save
    if err := s.repository.Save(ctx, user); err != nil {
        return nil, fmt.Errorf("failed to save user: %w", err)
    }
    
    // Publish event
    event := UserCreatedEvent{UserID: user.ID}
    if err := s.publisher.Publish(ctx, event); err != nil {
        // Log but don't fail
        log.Error("failed to publish user created event", err)
    }
    
    return user, nil
}
```

### Application Launcher Pattern

Проект использует специальный паттерн для запуска приложений с помощью `pkg/launcher`:

```go
// cmd/binance-tick-consumer/main.go
func main() {
    // configure logger
    logger := log.NewLogger(
        log.NewTextStdOutHandler(),
    )

    launcher.NewAppLauncher().
        // inject context that is canceled on SIGINT, SIGTERM
        WithGrexitContext().
        // inject logger into context
        WithLoggerContext(logger).
        // Run application, wait for it to finish
        WaitApplication(app_binance_tick_consumer.Launch).
        // log error if any
        LogIfError(logger)
}
```

### Application Service Structure

Каждое приложение имеет функцию `Launch` в своем app package:

```go
// internal/app/binance-tick-consumer/service.go
func Launch(ctx context.Context) error {
    logger, logStopService := log.MustGetLoggerFromContext(ctx).
        StartService("binance-tick-consumer")
    defer logStopService()

    // Load config
    var cfg = must.Return(
        config.ParseEnv[serviceConfig](),
    )
    logger.Info("service config loaded", "config", cfg)

    // Initialize clients
    redisClient, redisClose := must.Return2(
        redis.NewClient(ctx, cfg.RedisURL),
    )
    defer redisClose()

    // Initialize adapters
    var tickPublisher = adapter_redis.NewTickPublisher(redisClient)

    // Initialize services
    var tickProcessor = service_tick.NewTickProcessor(tickPublisher)

    // Initialize consumers  
    var binanceListener = adapter_binance.NewLPTickConsumer(
        cfg.Symbols, 
        tickProcessor,
        func(err error) {
            logger.Error("failed to handle tick event", "error", err)
        },
    )

    // Start the service
    if err := binanceListener.Launch(ctx); err != nil {
        return errors.Wrap(err, "failed to run binance tick consumer")
    }

    return nil
}
```

### Multi-Service Application

Для запуска нескольких сервисов одновременно используется `WaitApplications`:

```go
// cmd/all-in-one/main.go
func main() {
    logger := log.NewLogger(
        log.NewTextStdOutHandler(),
    )

    launcher.NewAppLauncher().
        WithGrexitContext().
        WithLoggerContext(logger).
        // Run multiple applications concurrently
        WaitApplications(
            app_binance_tick_consumer.Launch,
            app_ohlc_generator.Launch,
            app_backoffice_api.Launch,
        ).
        LogIfError(logger)
}
```

### Configuration Pattern

Каждое приложение определяет свою структуру конфигурации:

```go
// serviceConfig holds the configuration for the service
type serviceConfig struct {
    RedisURL string   `env:"REDIS_URL" required:"true"`
    Symbols  []string `env:"SYMBOLS" required:"true"`
}

// Load configuration from environment
var cfg = must.Return(
    config.ParseEnv[serviceConfig](),
)
```

## 🚀 Deployment

### Docker

#### Build Images
```bash
make docker-build
```

#### Run Services
```bash
make docker-up
```

#### View Logs
```bash
make docker-logs
```

### Environment Variables

Create `.env` file:
```env
# Database
DATABASE_URL=postgres://user:pass@localhost/dbname
DATABASE_MAX_CONNECTIONS=25

# Redis  
REDIS_URL=redis://localhost:6379
REDIS_MAX_CONNECTIONS=10

# gRPC
GRPC_PORT=8080
GRPC_TIMEOUT=30s

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
```

## 📊 Monitoring

### Health Checks
```go
type HealthChecker interface {
    Check(ctx context.Context) error
}

type HealthService struct {
    checkers map[string]HealthChecker
}

func (h *HealthService) CheckAll(ctx context.Context) map[string]error {
    results := make(map[string]error)
    for name, checker := range h.checkers {
        results[name] = checker.Check(ctx)
    }
    return results
}
```

### Metrics
```go
type Metrics interface {
    IncCounter(name string, labels map[string]string)
    ObserveHistogram(name string, value float64, labels map[string]string)
    SetGauge(name string, value float64, labels map[string]string)
}

// Usage in service
func (s *UserService) CreateUser(ctx context.Context, req CreateUserRequest) (*entity.User, error) {
    start := time.Now()
    defer func() {
        s.metrics.ObserveHistogram("user_create_duration", time.Since(start).Seconds(), nil)
    }()
    
    // ... implementation
    
    s.metrics.IncCounter("users_created", nil)
    return user, nil
}
```

### Structured Logging
```go
type Logger interface {
    Info(msg string, fields ...Field)
    Error(msg string, err error, fields ...Field)
    Debug(msg string, fields ...Field)
}

// Usage
func (s *UserService) CreateUser(ctx context.Context, req CreateUserRequest) (*entity.User, error) {
    s.logger.Info("creating user", 
        String("email", req.Email),
        String("trace_id", getTraceID(ctx)),
    )
    
    user, err := s.createUserInternal(ctx, req)
    if err != nil {
        s.logger.Error("failed to create user", err,
            String("email", req.Email),
            String("trace_id", getTraceID(ctx)),
        )
        return nil, err
    }
    
    s.logger.Info("user created successfully",
        String("user_id", user.ID),
        String("email", user.Email),
        String("trace_id", getTraceID(ctx)),
    )
    
    return user, nil
}
```

## 🔧 Troubleshooting

### Common Issues

#### Import Cycles
```
# Error: import cycle not allowed
package internal/service/user imports internal/adapter/postgres
package internal/adapter/postgres imports internal/service/user
```

**Solution**: Use interfaces and dependency injection
```go
// ✅ Good: Service depends on interface
type UserService struct {
    repository UserRepository // interface
}

// ✅ Good: Adapter implements interface
type PostgresUserRepository struct {
    db *sql.DB
}
```

#### Missing Context
```go
// ❌ Bad
func (s *UserService) GetUser(id string) (*entity.User, error)

// ✅ Good
func (s *UserService) GetUser(ctx context.Context, id string) (*entity.User, error)
```

#### Testing Issues
```bash
# Run specific test
go test -run TestUserService_CreateUser ./internal/service/user

# Run with verbose output
go test -v ./...

# Run with race detection
go test -race ./...
```

### Debugging

#### Enable Debug Logging
```env
LOG_LEVEL=debug
```

#### Use Trace IDs
```go
ctx = context.WithValue(ctx, "trace_id", generateTraceID())
```

#### Check Health Endpoints
```bash
curl http://localhost:8080/health
curl http://localhost:8080/metrics
```

## 📚 Resources

- [Go Documentation](https://golang.org/doc/)
- [Hexagonal Architecture](./ARCHITECTURE.md)
- [Testing Guide](./TESTING.md)
- [API Documentation](../api/README.md)

## 🤝 Contributing

1. Fork the repository
2. Create feature branch: `git checkout -b feature/new-feature`
3. Run tests: `make test`
4. Commit changes: `git commit -am 'Add new feature'`
5. Push branch: `git push origin feature/new-feature`
6. Create pull request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.