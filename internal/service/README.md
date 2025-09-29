# Services (Domain Services)

This directory contains domain services that encapsulate business logic that doesn't naturally belong to any specific entity or value object. These services contain pure business logic without any infrastructure concerns.

## What to create here:

### Domain Services
- **Business logic services** - Complex business operations that span multiple entities
- **Calculation services** - Business calculations and algorithms
- **Policy services** - Business policy implementations
- **Validation services** - Complex business validations
- **Factory services** - Object creation with complex business rules

### When to Create Domain Services:
1. **Multi-entity operations** - Logic that involves multiple entities
2. **Complex calculations** - Business calculations that are too complex for entities
3. **Business policies** - Rules that change frequently or are configurable
4. **External domain concepts** - Business concepts that don't fit in entities
5. **Stateless operations** - Pure functions that operate on domain objects

### Examples:
```go
// Pricing service for complex pricing rules
type PricingService interface {
    CalculateOrderTotal(order *Order, customer *Customer) (Money, error)
    ApplyDiscounts(order *Order, promotions []Promotion) error
    CalculateShipping(order *Order, address Address) (Money, error)
}

// User validation service for complex business rules
type UserValidationService interface {
    ValidateUserRegistration(user *User) error
    CanUserMakeOrder(user *User, order *Order) error
    CheckUserPermissions(user *User, resource Resource) error
}

// Notification policy service
type NotificationPolicyService interface {
    ShouldNotifyUser(user *User, event DomainEvent) bool
    GetNotificationChannels(user *User, notificationType NotificationType) []Channel
}
```

## Structure:
```
service/
├── pricing/        # Pricing and calculation services
├── validation/     # Business validation services
├── policy/         # Business policy services
├── factory/        # Complex object creation services
└── calculation/    # Business calculation services
```

## Key Principles:
- Contain pure business logic without infrastructure concerns
- Should be stateless and deterministic
- Can depend on entities, value objects, and repositories
- Should not depend on external frameworks or adapters
- Express business concepts in the ubiquitous language
- Keep focused on single responsibility
- Can be injected into application services and entities