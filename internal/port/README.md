# Ports

This directory contains all port definitions (interfaces) that define the contracts between the application core and external adapters. Ports represent the boundaries of the hexagonal architecture.

## What to create here:

### Primary Ports (Driving Ports)
These interfaces are implemented by the application core and called by primary adapters:

- **Use case interfaces** - Define application use cases and commands
- **Service interfaces** - Application service contracts  
- **API interfaces** - Define what the application can do
- **Handler interfaces** - Event and command handler contracts

### Secondary Ports (Driven Ports)
These interfaces are implemented by secondary adapters and called by the application core:

- **Repository interfaces** - Data persistence contracts
- **External service interfaces** - Third-party API contracts
- **Notification interfaces** - Email, SMS, push notification contracts
- **Cache interfaces** - Caching service contracts
- **Queue interfaces** - Message queue contracts
- **File storage interfaces** - File system operation contracts

### Examples:
```go
// Primary Port (Driving)
type UserUseCase interface {
    CreateUser(ctx context.Context, cmd CreateUserCommand) (*User, error)
    GetUser(ctx context.Context, id UserID) (*User, error)
    UpdateUser(ctx context.Context, cmd UpdateUserCommand) error
}

// Secondary Port (Driven) 
type UserRepository interface {
    Save(ctx context.Context, user *User) error
    FindByID(ctx context.Context, id UserID) (*User, error)
    FindByEmail(ctx context.Context, email Email) (*User, error)
    Delete(ctx context.Context, id UserID) error
}

type EmailService interface {
    SendWelcomeEmail(ctx context.Context, user *User) error
    SendPasswordResetEmail(ctx context.Context, user *User, token string) error
}
```

## Structure:
```
port/
├── primary/        # Primary ports (driving adapters call these)
├── secondary/      # Secondary ports (application calls these)  
├── repository/     # Repository interface definitions
├── service/        # External service interfaces
└── common/         # Shared port types and contracts
```

## Key Principles:
- Ports are just interfaces - no implementations here
- Define contracts between layers
- Should be stable and change infrequently  
- Use domain language in interface names and methods
- Keep interfaces focused and cohesive (Interface Segregation Principle)
- Depend only on entities and domain types