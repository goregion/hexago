# Entities (Domain Layer)

This directory contains the core business entities and domain logic. This is the heart of the hexagonal architecture and should contain the most important business rules and logic.

## What to create here:

### Domain Entities
- **Business entities** - Core objects that represent business concepts (User, Order, Product, etc.)
- **Value objects** - Immutable objects that represent values (Email, Money, Address, etc.)
- **Aggregates** - Collections of entities that are treated as a single unit
- **Domain events** - Events that represent something important that happened in the domain

### Domain Logic
- **Business rules** - Core business validation and logic
- **Domain services** - Services that contain domain logic that doesn't naturally fit in entities
- **Specifications** - Business rule specifications for complex validations
- **Domain exceptions** - Business-specific error types

### Examples:
```go
type User struct {
    ID       UserID
    Email    Email
    Profile  UserProfile
    Status   UserStatus
    // ... business methods
}

type Order struct {
    ID          OrderID
    CustomerID  CustomerID
    Items       []OrderItem
    Status      OrderStatus
    Total       Money
    // ... business methods
}

type Email struct {
    value string
    // ... validation methods
}
```

## Key Principles:
- Entities should contain business logic and rules
- Should be independent of external frameworks and libraries
- Rich domain models with behavior, not just data containers
- Entities should protect their invariants
- Should not depend on any external layers
- Contains the ubiquitous language of the domain