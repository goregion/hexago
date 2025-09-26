# Hexagonal Architecture Guide

## 🏛️ Overview

Hexagonal Architecture, also known as Ports and Adapters pattern, was introduced by Alistair Cockburn. It aims to create loosely coupled application components that can be easily connected to their software environment by means of ports and adapters.

## 🎯 Core Principles

### 1. **Business Logic Isolation**
- Keep business rules and domain logic independent of external concerns
- No direct dependencies on databases, web frameworks, or external services in the core
- Business logic should not know about HTTP, databases, or any I/O

### 2. **Dependency Inversion**
- High-level modules should not depend on low-level modules
- Both should depend on abstractions (interfaces)
- Abstractions should not depend on details

### 3. **Testability**
- Core business logic can be tested without external dependencies
- Mock implementations can be easily substituted during testing
- Isolated unit tests run fast and are reliable

## 🏗️ Architecture Layers

```
┌─────────────────────────────────────────────────────┐
│                  External World                     │
│  (Web, CLI, Message Queue, Database, File System)  │
└─────────────┬───────────────────────┬───────────────┘
              │                       │
         ┌────▼────┐               ┌───▼────┐
         │ Driving │               │ Driven │
         │Adapters │               │Adapters│
         │(Primary)│               │(Second)│
         └────┬────┘               └───▲────┘
              │                       │
         ┌────▼───────────────────────┴────┐
         │            Ports                │
         │  (Interfaces/Contracts)         │
         └────┬───────────────────────▲────┘
              │                       │
         ┌────▼───────────────────────┴────┐
         │        Application Core         │
         │  (Business Logic & Entities)    │
         └─────────────────────────────────┘
```

### Layer Descriptions

#### **Application Core** (`internal/entity/`, `internal/service/`)
- **Entities**: Domain objects with business rules
- **Services**: Business logic implementations
- **No dependencies** on external frameworks or libraries

#### **Ports** (`internal/port/`)
- **Driving Ports**: Interfaces for incoming operations (use cases)
- **Driven Ports**: Interfaces for outgoing operations (repositories, external services)
- Define **contracts** between core and external world

#### **Adapters** (`internal/adapter/`)
- **Driving Adapters**: Implement driving ports (HTTP handlers, CLI, gRPC)
- **Driven Adapters**: Implement driven ports (databases, external APIs, message queues)
- Handle **translation** between external formats and internal models

#### **Application Services** (`internal/app/`)
- **Orchestrate** business workflows
- **Coordinate** between multiple domain services
- Handle **cross-cutting concerns** (transactions, logging, metrics)

## 📁 Directory Structure

```
internal/
├── entity/           # Domain entities and value objects
│   ├── user.go       # User domain entity
│   └── order.go      # Order domain entity
├── port/             # Interface definitions
│   ├── user.go       # User-related interfaces
│   └── payment.go    # Payment-related interfaces  
├── service/          # Business logic implementations
│   ├── user/         # User domain services
│   └── order/        # Order domain services
├── adapter/          # External integrations
│   ├── http/         # HTTP API adapter
│   ├── postgres/     # PostgreSQL adapter
│   ├── redis/        # Redis cache adapter
│   └── stripe/       # Stripe payment adapter
└── app/              # Application orchestration
    ├── user-service/ # User management app
    └── order-service/# Order processing app
```

## 🔌 Ports and Adapters

### Driving Ports (Primary/Inbound)
Use cases that the application offers to the outside world.

```go
// internal/port/user.go
type UserService interface {
    CreateUser(ctx context.Context, req CreateUserRequest) (*entity.User, error)
    GetUser(ctx context.Context, id string) (*entity.User, error)
}
```

### Driven Ports (Secondary/Outbound)
Operations that the application needs from external systems.

```go
// internal/port/user.go  
type UserRepository interface {
    Save(ctx context.Context, user *entity.User) error
    FindByID(ctx context.Context, id string) (*entity.User, error)
}

type EmailSender interface {
    SendWelcomeEmail(ctx context.Context, user *entity.User) error
}
```

### Driving Adapters (Primary)
Implement driving ports and handle incoming requests.

```go
// internal/adapter/http/user_handler.go
type UserHandler struct {
    userService port.UserService
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    // Parse HTTP request
    // Call userService.CreateUser()
    // Return HTTP response
}
```

### Driven Adapters (Secondary)
Implement driven ports and handle outgoing operations.

```go
// internal/adapter/postgres/user_repository.go
type UserRepository struct {
    db *sql.DB
}

func (r *UserRepository) Save(ctx context.Context, user *entity.User) error {
    // Save user to PostgreSQL database
}
```

## 🎭 Implementation Patterns

### 1. **Entity Pattern**
Domain entities encapsulate business rules and data.

```go
// internal/entity/user.go
type User struct {
    ID       string
    Email    string
    Name     string
    Status   UserStatus
}

func (u *User) Activate() error {
    if u.Status == StatusActive {
        return ErrUserAlreadyActive
    }
    u.Status = StatusActive
    return nil
}

func (u *User) IsValid() error {
    if u.Email == "" {
        return ErrInvalidEmail
    }
    return nil
}
```

### 2. **Service Pattern**
Services implement business workflows.

```go
// internal/service/user/user_service.go
type UserService struct {
    repository port.UserRepository
    emailSender port.EmailSender
}

func (s *UserService) CreateUser(ctx context.Context, req CreateUserRequest) (*entity.User, error) {
    // Validate request
    user := &entity.User{
        ID:    generateID(),
        Email: req.Email,
        Name:  req.Name,
        Status: entity.StatusPending,
    }
    
    if err := user.IsValid(); err != nil {
        return nil, err
    }
    
    // Save to repository
    if err := s.repository.Save(ctx, user); err != nil {
        return nil, err
    }
    
    // Send welcome email
    if err := s.emailSender.SendWelcomeEmail(ctx, user); err != nil {
        // Log error but don't fail the operation
        log.Error("Failed to send welcome email", err)
    }
    
    return user, nil
}
```

### 3. **Adapter Pattern**
Adapters translate between external systems and internal models.

```go
// internal/adapter/postgres/user_repository.go
type UserRepository struct {
    db *sql.DB
}

func (r *UserRepository) Save(ctx context.Context, user *entity.User) error {
    query := `INSERT INTO users (id, email, name, status) VALUES ($1, $2, $3, $4)`
    _, err := r.db.ExecContext(ctx, query, user.ID, user.Email, user.Name, user.Status)
    return err
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
    query := `SELECT id, email, name, status FROM users WHERE id = $1`
    
    var user entity.User
    err := r.db.QueryRowContext(ctx, query, id).Scan(
        &user.ID, &user.Email, &user.Name, &user.Status,
    )
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, ErrUserNotFound
        }
        return nil, err
    }
    
    return &user, nil
}
```

## 🧪 Testing Strategy

### Unit Tests
Test business logic in isolation.

```go
func TestUserService_CreateUser(t *testing.T) {
    // Arrange
    mockRepo := &mockUserRepository{}
    mockEmailSender := &mockEmailSender{}
    service := user.NewUserService(mockRepo, mockEmailSender)
    
    req := CreateUserRequest{
        Email: "test@example.com",
        Name:  "Test User",
    }
    
    // Act
    result, err := service.CreateUser(context.Background(), req)
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, req.Email, result.Email)
    
    // Verify interactions
    assert.True(t, mockRepo.SaveCalled)
    assert.True(t, mockEmailSender.SendCalled)
}
```

### Integration Tests
Test with real adapters.

```go
func TestUserRepository_Integration(t *testing.T) {
    // Setup test database
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)
    
    repo := postgres.NewUserRepository(db)
    
    user := &entity.User{
        ID:    "user-123",
        Email: "test@example.com",
        Name:  "Test User",
    }
    
    // Test Save
    err := repo.Save(context.Background(), user)
    assert.NoError(t, err)
    
    // Test FindByID
    found, err := repo.FindByID(context.Background(), user.ID)
    assert.NoError(t, err)
    assert.Equal(t, user.Email, found.Email)
}
```

## 📊 Benefits

### ✅ **Advantages**

1. **Testability**: Easy to test business logic in isolation
2. **Flexibility**: Easy to swap adapters (e.g., change from MySQL to PostgreSQL)
3. **Maintainability**: Clear separation of concerns
4. **Technology Independence**: Core logic not tied to specific frameworks
5. **Evolution**: Easy to add new interfaces (HTTP, gRPC, CLI)

### ⚠️ **Considerations**

1. **Complexity**: More layers and abstractions
2. **Over-engineering**: May be overkill for simple applications  
3. **Learning Curve**: Requires understanding of DDD concepts
4. **Boilerplate**: More interfaces and structures to maintain

## 🎯 Best Practices

### 1. **Start Simple**
- Begin with clear domain boundaries
- Add complexity only when needed
- Focus on business value first

### 2. **Define Clear Contracts**
- Make port interfaces explicit and focused
- Use meaningful names for methods and types
- Document expected behavior and error conditions

### 3. **Keep Dependencies Flowing Inward**
- Core should never import adapter packages
- Use dependency injection to wire components
- Interfaces belong to the consumer, not the implementer

### 4. **Handle Errors Gracefully**
- Define domain-specific error types
- Don't leak implementation details in errors
- Use appropriate logging levels

### 5. **Test at the Right Level**
- Unit test business logic with mocks
- Integration test adapters with real dependencies
- End-to-end test critical user journeys

## 🔧 Common Patterns

### Repository Pattern
```go
type UserRepository interface {
    Save(ctx context.Context, user *entity.User) error
    FindByID(ctx context.Context, id string) (*entity.User, error)
    FindByEmail(ctx context.Context, email string) (*entity.User, error)
    Delete(ctx context.Context, id string) error
}
```

### Event Publisher Pattern
```go
type EventPublisher interface {
    Publish(ctx context.Context, event Event) error
}

type Event interface {
    Type() string
    AggregateID() string
}
```

### Command/Query Separation
```go
// Commands (write operations)
type UserCommandService interface {
    CreateUser(ctx context.Context, cmd CreateUserCommand) error
    UpdateUser(ctx context.Context, cmd UpdateUserCommand) error
    DeleteUser(ctx context.Context, cmd DeleteUserCommand) error
}

// Queries (read operations)  
type UserQueryService interface {
    GetUser(ctx context.Context, query GetUserQuery) (*UserView, error)
    ListUsers(ctx context.Context, query ListUsersQuery) ([]*UserView, error)
}
```

## 📚 Further Reading

- [Hexagonal Architecture by Alistair Cockburn](https://alistair.cockburn.us/hexagonal-architecture/)
- [Clean Architecture by Robert Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Domain-Driven Design by Eric Evans](https://www.domainlanguage.com/ddd/)
- [Growing Object-Oriented Software, Guided by Tests](http://www.growing-object-oriented-software.com/)