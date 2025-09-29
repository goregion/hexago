# Application Layer

This directory contains the application layer of the hexagonal architecture. It orchestrates the business logic by coordinating between ports, services, and use cases.

## What to create here:

### Use Cases / Application Services
- **Command handlers** - Handle business commands and operations
- **Query handlers** - Handle data retrieval and filtering
- **Application services** - Coordinate complex business workflows
- **DTOs (Data Transfer Objects)** - Input/output models for use cases
- **Application interfaces** - Contracts for application services

### Application Logic
- **Workflow orchestration** - Multi-step business processes
- **Transaction management** - Cross-service transactions
- **Validation rules** - Input validation and business rule validation
- **Error handling** - Application-level error processing
- **Event coordination** - Domain event handling and publishing

### Use Case Examples:
- `CreateUserUseCase` - Handle user registration workflow
- `ProcessOrderUseCase` - Handle order processing pipeline  
- `GenerateReportUseCase` - Handle report generation
- `SendNotificationUseCase` - Handle notification workflows

## Structure:
```
app/
├── usecase/        # Use case implementations
├── dto/            # Data transfer objects
├── command/        # Command patterns and handlers
├── query/          # Query patterns and handlers
└── service/        # Application services
```

## Key Principles:
- Application layer should be technology-agnostic
- Contains application-specific business rules
- Orchestrates calls to domain services and repositories
- Should not contain framework-specific code
- Depends only on entities and ports (inward dependencies)
- Implements use cases that represent user intentions