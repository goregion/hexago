# Adapters

This directory contains all external adapters that implement the ports defined in the `port` package. Adapters are the outermost layer of the hexagonal architecture and handle communication with external systems.

## What to create here:

### Primary Adapters (Driving)
- **HTTP handlers/controllers** - REST API endpoints that receive HTTP requests
- **gRPC servers** - gRPC service implementations
- **CLI commands** - Command-line interface handlers
- **Message queue consumers** - Handlers for incoming messages from queues like RabbitMQ, Kafka
- **WebSocket handlers** - Real-time communication handlers
- **GraphQL resolvers** - GraphQL query and mutation resolvers

### Secondary Adapters (Driven)
- **Database repositories** - Implementations of repository interfaces (PostgreSQL, MySQL, MongoDB, etc.)
- **External API clients** - HTTP clients for third-party services
- **Message queue publishers** - Publishers for sending messages to queues
- **File system handlers** - File operations (local, S3, etc.)
- **Cache implementations** - Redis, Memcached implementations
- **Email service clients** - SMTP, SendGrid, Mailgun clients
- **Payment gateway clients** - Stripe, PayPal, etc.

## Structure:
```
adapter/
├── http/           # HTTP handlers and middleware
├── grpc/           # gRPC server implementations
├── cli/            # Command-line interface
├── repository/     # Database implementations
├── cache/          # Cache implementations
├── external/       # External API clients
└── queue/          # Message queue implementations
```

## Key Principles:
- Adapters should only contain infrastructure concerns
- Business logic should be delegated to the application layer
- Adapters should implement interfaces defined in the `port` package
- Keep adapters thin and focused on data transformation