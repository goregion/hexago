# Testing Guide

## 🧪 Testing Philosophy

Our testing strategy follows the **Test Pyramid** principle:

```
                    /\
                   /  \
               E2E /    \ (Few, Slow, High Value)
                  /______\
                 /        \
        Integration/        \ (Some, Medium Speed)
               /____________\
              /              \
            Unit/              \ (Many, Fast, Focused)
           /__________________\
```

## 🏗️ Test Structure

### Directory Layout

```
tests/
├── unit/                    # Fast, isolated unit tests
│   ├── mocks.go            # Test doubles and mocks
│   ├── adapter/            # Adapter unit tests
│   ├── service/            # Service unit tests
│   └── entity/             # Entity unit tests
├── integration/            # Integration tests with real dependencies
│   ├── helpers.go          # Test helpers and utilities
│   ├── basic_test.go       # Basic integration tests
│   ├── concurrent_test.go  # Concurrency tests
│   └── error_test.go       # Error handling tests
└── e2e/                    # End-to-end tests (future)
    └── api_test.go         # Full API workflow tests
```

## 🎯 Unit Testing

### Principles

1. **Fast** - Tests should run in milliseconds
2. **Isolated** - No external dependencies (databases, networks, files)
3. **Deterministic** - Same input always produces same output
4. **Independent** - Tests don't depend on each other

### Testing Entities

Entities contain business logic and should be thoroughly tested:

```go
// internal/entity/user_test.go
package entity

import (
    "testing"
    "time"
)

func TestUser_Activate(t *testing.T) {
    tests := []struct {
        name        string
        user        User
        wantErr     bool
        wantStatus  UserStatus
    }{
        {
            name: "activate pending user",
            user: User{
                ID:     "user-123",
                Status: StatusPending,
            },
            wantErr:    false,
            wantStatus: StatusActive,
        },
        {
            name: "activate already active user",
            user: User{
                ID:     "user-123", 
                Status: StatusActive,
            },
            wantErr:    true,
            wantStatus: StatusActive,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.user.Activate()
            
            if tt.wantErr && err == nil {
                t.Error("expected error but got nil")
            }
            if !tt.wantErr && err != nil {
                t.Errorf("expected no error but got: %v", err)
            }
            
            if tt.user.Status != tt.wantStatus {
                t.Errorf("expected status %v, got %v", tt.wantStatus, tt.user.Status)
            }
        })
    }
}

func TestUser_IsValid(t *testing.T) {
    tests := []struct {
        name    string
        user    User
        wantErr bool
    }{
        {
            name: "valid user",
            user: User{
                ID:    "user-123",
                Email: "test@example.com",
                Name:  "Test User",
            },
            wantErr: false,
        },
        {
            name: "invalid email",
            user: User{
                ID:   "user-123",
                Name: "Test User",
                // Missing email
            },
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.user.IsValid()
            
            if tt.wantErr && err == nil {
                t.Error("expected error but got nil")
            }
            if !tt.wantErr && err != nil {
                t.Errorf("expected no error but got: %v", err)
            }
        })
    }
}
```

### Testing Services

Services contain business logic and orchestration:

```go
// tests/unit/service/user_service_test.go
package service

import (
    "context"
    "testing"
    
    "github.com/goregion/hexago/internal/entity"
    "github.com/goregion/hexago/internal/service/user"
)

// Mock implementations
type mockUserRepository struct {
    users   map[string]*entity.User
    saveErr error
    findErr error
}

func (m *mockUserRepository) Save(ctx context.Context, user *entity.User) error {
    if m.saveErr != nil {
        return m.saveErr
    }
    if m.users == nil {
        m.users = make(map[string]*entity.User)
    }
    m.users[user.ID] = user
    return nil
}

func (m *mockUserRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
    if m.findErr != nil {
        return nil, m.findErr
    }
    user, exists := m.users[id]
    if !exists {
        return nil, entity.ErrUserNotFound
    }
    return user, nil
}

type mockEventPublisher struct {
    publishedEvents []interface{}
    publishErr      error
}

func (m *mockEventPublisher) Publish(ctx context.Context, event interface{}) error {
    if m.publishErr != nil {
        return m.publishErr
    }
    m.publishedEvents = append(m.publishedEvents, event)
    return nil
}

func TestUserService_CreateUser(t *testing.T) {
    tests := []struct {
        name           string
        request        user.CreateUserRequest
        setupMocks     func(*mockUserRepository, *mockEventPublisher)
        wantErr        bool
        wantUser       bool
        wantEventCount int
    }{
        {
            name: "successful user creation",
            request: user.CreateUserRequest{
                Email: "test@example.com",
                Name:  "Test User",
            },
            setupMocks: func(repo *mockUserRepository, pub *mockEventPublisher) {
                // No special setup needed for success case
            },
            wantErr:        false,
            wantUser:       true,
            wantEventCount: 1,
        },
        {
            name: "invalid email",
            request: user.CreateUserRequest{
                Email: "invalid-email",
                Name:  "Test User",
            },
            setupMocks: func(repo *mockUserRepository, pub *mockEventPublisher) {
                // No setup needed
            },
            wantErr:        true,
            wantUser:       false,
            wantEventCount: 0,
        },
        {
            name: "repository save error",
            request: user.CreateUserRequest{
                Email: "test@example.com",
                Name:  "Test User",
            },
            setupMocks: func(repo *mockUserRepository, pub *mockEventPublisher) {
                repo.saveErr = errors.New("database error")
            },
            wantErr:        true,
            wantUser:       false,
            wantEventCount: 0,
        },
        {
            name: "event publish error (should not fail operation)",
            request: user.CreateUserRequest{
                Email: "test@example.com",
                Name:  "Test User",
            },
            setupMocks: func(repo *mockUserRepository, pub *mockEventPublisher) {
                pub.publishErr = errors.New("publish error")
            },
            wantErr:        false, // Service should not fail if event publishing fails
            wantUser:       true,
            wantEventCount: 0, // Event not published due to error
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Arrange
            mockRepo := &mockUserRepository{}
            mockPublisher := &mockEventPublisher{}
            tt.setupMocks(mockRepo, mockPublisher)
            
            service := user.NewUserService(mockRepo, mockPublisher)
            
            // Act
            result, err := service.CreateUser(context.Background(), tt.request)
            
            // Assert
            if tt.wantErr && err == nil {
                t.Error("expected error but got nil")
            }
            if !tt.wantErr && err != nil {
                t.Errorf("expected no error but got: %v", err)
            }
            
            if tt.wantUser && result == nil {
                t.Error("expected user to be returned")
            }
            if !tt.wantUser && result != nil {
                t.Error("expected no user to be returned")
            }
            
            if len(mockPublisher.publishedEvents) != tt.wantEventCount {
                t.Errorf("expected %d events, got %d", tt.wantEventCount, len(mockPublisher.publishedEvents))
            }
            
            // Verify user was saved correctly
            if tt.wantUser {
                savedUser := mockRepo.users[result.ID]
                if savedUser == nil {
                    t.Error("user was not saved to repository")
                }
                if savedUser.Email != tt.request.Email {
                    t.Errorf("expected email %s, got %s", tt.request.Email, savedUser.Email)
                }
            }
        })
    }
}

func TestUserService_GetUser(t *testing.T) {
    // Arrange
    existingUser := &entity.User{
        ID:    "user-123",
        Email: "test@example.com",
        Name:  "Test User",
    }
    
    mockRepo := &mockUserRepository{
        users: map[string]*entity.User{
            existingUser.ID: existingUser,
        },
    }
    mockPublisher := &mockEventPublisher{}
    
    service := user.NewUserService(mockRepo, mockPublisher)
    
    tests := []struct {
        name    string
        userID  string
        wantErr bool
        wantUser *entity.User
    }{
        {
            name:     "existing user",
            userID:   "user-123",
            wantErr:  false,
            wantUser: existingUser,
        },
        {
            name:     "non-existing user",
            userID:   "user-999",
            wantErr:  true,
            wantUser: nil,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Act
            result, err := service.GetUser(context.Background(), tt.userID)
            
            // Assert
            if tt.wantErr && err == nil {
                t.Error("expected error but got nil")
            }
            if !tt.wantErr && err != nil {
                t.Errorf("expected no error but got: %v", err)
            }
            
            if tt.wantUser != nil && result == nil {
                t.Error("expected user but got nil")
            }
            if tt.wantUser == nil && result != nil {
                t.Error("expected nil but got user")
            }
            
            if tt.wantUser != nil && result != nil {
                if result.ID != tt.wantUser.ID {
                    t.Errorf("expected user ID %s, got %s", tt.wantUser.ID, result.ID)
                }
            }
        })
    }
}
```

### Testing Adapters (Unit Level)

Test adapters in isolation using mocks for external dependencies:

```go
// tests/unit/adapter/postgres_user_repository_test.go
package postgres

import (
    "context"
    "database/sql"
    "testing"
    
    "github.com/DATA-DOG/go-sqlmock"
    "github.com/goregion/hexago/internal/entity"
    "github.com/goregion/hexago/internal/adapter/postgres"
)

func TestUserRepository_Save(t *testing.T) {
    // Create mock database
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("failed to create mock db: %v", err)
    }
    defer db.Close()
    
    repo := postgres.NewUserRepository(db)
    
    tests := []struct {
        name      string
        user      *entity.User
        setupMock func(sqlmock.Sqlmock)
        wantErr   bool
    }{
        {
            name: "successful save",
            user: &entity.User{
                ID:    "user-123",
                Email: "test@example.com",
                Name:  "Test User",
                Status: entity.StatusActive,
            },
            setupMock: func(mock sqlmock.Sqlmock) {
                mock.ExpectExec("INSERT INTO users").
                    WithArgs("user-123", "test@example.com", "Test User", entity.StatusActive).
                    WillReturnResult(sqlmock.NewResult(1, 1))
            },
            wantErr: false,
        },
        {
            name: "database error",
            user: &entity.User{
                ID:    "user-123",
                Email: "test@example.com",
                Name:  "Test User",
                Status: entity.StatusActive,
            },
            setupMock: func(mock sqlmock.Sqlmock) {
                mock.ExpectExec("INSERT INTO users").
                    WithArgs("user-123", "test@example.com", "Test User", entity.StatusActive).
                    WillReturnError(sql.ErrConnDone)
            },
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Setup mock expectations
            tt.setupMock(mock)
            
            // Act
            err := repo.Save(context.Background(), tt.user)
            
            // Assert
            if tt.wantErr && err == nil {
                t.Error("expected error but got nil")
            }
            if !tt.wantErr && err != nil {
                t.Errorf("expected no error but got: %v", err)
            }
            
            // Verify all expectations were met
            if err := mock.ExpectationsWereMet(); err != nil {
                t.Errorf("unmet expectations: %v", err)
            }
        })
    }
}
```

## 🔗 Integration Testing

Integration tests verify that components work together correctly with real dependencies.

### Database Integration Tests

```go
// tests/integration/user_repository_test.go
package integration

import (
    "context"
    "database/sql"
    "testing"
    
    _ "github.com/go-sql-driver/mysql"
    "github.com/goregion/hexago/internal/entity"
    "github.com/goregion/hexago/internal/adapter/postgres"
)

func setupTestDB(t *testing.T) *sql.DB {
    // Use test database or test container
    db, err := sql.Open("mysql", "test:test@tcp(localhost:3306)/test_db")
    if err != nil {
        t.Fatalf("failed to connect to test database: %v", err)
    }
    
    // Run migrations or setup schema
    setupSchema(t, db)
    
    return db
}

func cleanupTestDB(t *testing.T, db *sql.DB) {
    // Clean up test data
    _, err := db.Exec("DELETE FROM users")
    if err != nil {
        t.Logf("failed to cleanup test data: %v", err)
    }
    db.Close()
}

func setupSchema(t *testing.T, db *sql.DB) {
    schema := `
    CREATE TABLE IF NOT EXISTS users (
        id VARCHAR(255) PRIMARY KEY,
        email VARCHAR(255) UNIQUE NOT NULL,
        name VARCHAR(255) NOT NULL,
        status VARCHAR(50) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );`
    
    _, err := db.Exec(schema)
    if err != nil {
        t.Fatalf("failed to setup schema: %v", err)
    }
}

func TestUserRepository_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test in short mode")
    }
    
    // Setup
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)
    
    repo := postgres.NewUserRepository(db)
    
    user := &entity.User{
        ID:    "user-123",
        Email: "test@example.com",
        Name:  "Test User",
        Status: entity.StatusActive,
    }
    
    t.Run("save and find user", func(t *testing.T) {
        // Save user
        err := repo.Save(context.Background(), user)
        if err != nil {
            t.Fatalf("failed to save user: %v", err)
        }
        
        // Find user by ID
        found, err := repo.FindByID(context.Background(), user.ID)
        if err != nil {
            t.Fatalf("failed to find user: %v", err)
        }
        
        // Verify user data
        if found.ID != user.ID {
            t.Errorf("expected ID %s, got %s", user.ID, found.ID)
        }
        if found.Email != user.Email {
            t.Errorf("expected email %s, got %s", user.Email, found.Email)
        }
        if found.Name != user.Name {
            t.Errorf("expected name %s, got %s", user.Name, found.Name)
        }
    })
    
    t.Run("find non-existing user", func(t *testing.T) {
        _, err := repo.FindByID(context.Background(), "non-existing")
        if err == nil {
            t.Error("expected error for non-existing user")
        }
        if !errors.Is(err, entity.ErrUserNotFound) {
            t.Errorf("expected ErrUserNotFound, got %v", err)
        }
    })
}
```

### gRPC Integration Tests

```go
// tests/integration/grpc_api_test.go
package integration

import (
    "context"
    "net"
    "testing"
    
    "google.golang.org/grpc"
    "google.golang.org/grpc/test/bufconn"
    
    pb "github.com/goregion/hexago/internal/adapter/grpc-api/gen"
    "github.com/goregion/hexago/internal/adapter/grpc-api/impl"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func bufDialer(context.Context, string) (net.Conn, error) {
    return lis.Dial()
}

func setupGRPCServer(t *testing.T) pb.OHLCServiceClient {
    lis = bufconn.Listen(bufSize)
    
    s := grpc.NewServer()
    
    // Setup your gRPC server implementation
    server := impl.NewOHLCServer(/* dependencies */)
    pb.RegisterOHLCServiceServer(s, server)
    
    go func() {
        if err := s.Serve(lis); err != nil {
            t.Logf("Server exited with error: %v", err)
        }
    }()
    
    // Create client connection
    conn, err := grpc.DialContext(context.Background(), "bufnet", 
        grpc.WithContextDialer(bufDialer), 
        grpc.WithInsecure())
    if err != nil {
        t.Fatalf("Failed to dial bufnet: %v", err)
    }
    
    return pb.NewOHLCServiceClient(conn)
}

func TestOHLCService_GetOHLC(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test in short mode")
    }
    
    client := setupGRPCServer(t)
    
    tests := []struct {
        name    string
        request *pb.GetOHLCRequest
        wantErr bool
    }{
        {
            name: "valid request",
            request: &pb.GetOHLCRequest{
                Symbol: "BTCUSDT",
                From:   1234567890,
                To:     1234567900,
            },
            wantErr: false,
        },
        {
            name: "invalid symbol",
            request: &pb.GetOHLCRequest{
                Symbol: "",
                From:   1234567890,
                To:     1234567900,
            },
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            resp, err := client.GetOHLC(context.Background(), tt.request)
            
            if tt.wantErr && err == nil {
                t.Error("expected error but got nil")
            }
            if !tt.wantErr && err != nil {
                t.Errorf("expected no error but got: %v", err)
            }
            
            if !tt.wantErr && resp == nil {
                t.Error("expected response but got nil")
            }
        })
    }
}
```

## 🎭 Test Doubles

### Types of Test Doubles

1. **Dummy** - Objects passed around but never used
2. **Fake** - Working implementations with shortcuts (in-memory database)
3. **Stub** - Provide canned answers to calls
4. **Mock** - Pre-programmed with expectations

### Mock Implementation Guidelines

```go
// tests/unit/mocks.go
package unit

import (
    "context"
    
    "github.com/goregion/hexago/internal/entity"
)

// MockUserRepository implements port.UserRepository for testing
type MockUserRepository struct {
    // Control behavior
    SaveFunc    func(ctx context.Context, user *entity.User) error
    FindByIDFunc func(ctx context.Context, id string) (*entity.User, error)
    
    // Track calls
    SaveCalls    []SaveCall
    FindByIDCalls []FindByIDCall
}

type SaveCall struct {
    Ctx  context.Context
    User *entity.User
}

type FindByIDCall struct {
    Ctx context.Context
    ID  string
}

func (m *MockUserRepository) Save(ctx context.Context, user *entity.User) error {
    m.SaveCalls = append(m.SaveCalls, SaveCall{Ctx: ctx, User: user})
    
    if m.SaveFunc != nil {
        return m.SaveFunc(ctx, user)
    }
    
    return nil // Default behavior
}

func (m *MockUserRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
    m.FindByIDCalls = append(m.FindByIDCalls, FindByIDCall{Ctx: ctx, ID: id})
    
    if m.FindByIDFunc != nil {
        return m.FindByIDFunc(ctx, id)
    }
    
    return nil, entity.ErrUserNotFound // Default behavior
}

// Helper methods for assertions
func (m *MockUserRepository) WasSaveCalled() bool {
    return len(m.SaveCalls) > 0
}

func (m *MockUserRepository) SaveCallCount() int {
    return len(m.SaveCalls)
}

func (m *MockUserRepository) LastSaveCall() *SaveCall {
    if len(m.SaveCalls) == 0 {
        return nil
    }
    return &m.SaveCalls[len(m.SaveCalls)-1]
}
```

## 🏃‍♂️ Running Tests

### Commands

```bash
# Run all tests
make test

# Run only unit tests
make test-unit

# Run only integration tests  
make test-integration

# Run tests with coverage
make test-coverage

# Run tests with race detection
make test-race

# Run specific test
go test -run TestUserService_CreateUser ./tests/unit/service

# Run tests in verbose mode
go test -v ./tests/...

# Run tests in short mode (skip integration tests)
go test -short ./tests/...
```

### Test Flags

```bash
# Parallel execution
go test -parallel 4 ./tests/...

# Timeout
go test -timeout 30s ./tests/...

# CPU profiling
go test -cpuprofile cpu.prof ./tests/...

# Memory profiling  
go test -memprofile mem.prof ./tests/...

# Benchmark
go test -bench . ./tests/...
```

## 📊 Coverage Analysis

### Generate Coverage Report

```bash
# Generate coverage profile
go test -coverprofile=coverage.out ./tests/...

# View coverage in terminal
go tool cover -func=coverage.out

# Generate HTML report
go tool cover -html=coverage.out -o coverage.html
```

### Coverage Targets

- **Overall**: > 80%
- **Business Logic** (services): > 90%
- **Domain Entities**: > 95%
- **Adapters**: > 70%

## 🐛 Testing Anti-Patterns

### ❌ Avoid These

#### Testing Implementation Details
```go
// ❌ Bad - testing internal implementation
func TestUserService_CreateUser_CallsRepository(t *testing.T) {
    mockRepo := &MockUserRepository{}
    service := NewUserService(mockRepo)
    
    service.CreateUser(ctx, request)
    
    if !mockRepo.WasSaveCalled() {
        t.Error("expected Save to be called") // Testing implementation
    }
}
```

#### Fragile Tests
```go
// ❌ Bad - depends on specific time
func TestUser_CreatedAt(t *testing.T) {
    user := CreateUser()
    
    if user.CreatedAt != time.Now() { // Will fail due to timing
        t.Error("wrong creation time")
    }
}
```

#### Giant Test Methods
```go
// ❌ Bad - testing too many things
func TestUserService_Everything(t *testing.T) {
    // Tests create, update, delete, find all in one method
    // Hard to understand what failed
}
```

### ✅ Better Approaches

#### Test Behavior, Not Implementation
```go
// ✅ Good - testing behavior
func TestUserService_CreateUser_ReturnsUser(t *testing.T) {
    service := NewUserService(mockRepo)
    
    user, err := service.CreateUser(ctx, request)
    
    assert.NoError(t, err)
    assert.Equal(t, request.Email, user.Email) // Testing behavior
}
```

#### Use Test Helpers
```go
// ✅ Good - helper for time-sensitive tests
func TestUser_CreatedAt(t *testing.T) {
    now := time.Now()
    user := CreateUser()
    
    if user.CreatedAt.Before(now) || user.CreatedAt.After(now.Add(time.Second)) {
        t.Error("creation time not within expected range")
    }
}
```

#### One Concept Per Test
```go
// ✅ Good - focused tests
func TestUserService_CreateUser_ValidInput(t *testing.T) {
    // Test successful creation
}

func TestUserService_CreateUser_InvalidEmail(t *testing.T) {
    // Test invalid email handling
}

func TestUserService_CreateUser_RepositoryError(t *testing.T) {
    // Test repository error handling
}
```

## 🔧 Test Configuration

### Environment Variables

```bash
# .env.test
DATABASE_URL=postgres://test:test@localhost:5432/test_db
REDIS_URL=redis://localhost:6379/1
LOG_LEVEL=error
```

### Test Helpers

```go
// tests/helpers/database.go
package helpers

import (
    "database/sql"
    "testing"
)

func SetupTestDB(t *testing.T) *sql.DB {
    t.Helper() // Mark as helper function
    
    db, err := sql.Open("postgres", os.Getenv("TEST_DATABASE_URL"))
    if err != nil {
        t.Fatalf("failed to connect to test database: %v", err)
    }
    
    // Run cleanup when test finishes
    t.Cleanup(func() {
        cleanupDB(t, db)
        db.Close()
    })
    
    return db
}

func cleanupDB(t *testing.T, db *sql.DB) {
    t.Helper()
    
    tables := []string{"users", "orders", "payments"}
    for _, table := range tables {
        _, err := db.Exec(fmt.Sprintf("DELETE FROM %s", table))
        if err != nil {
            t.Logf("failed to cleanup table %s: %v", table, err)
        }
    }
}
```

## 📈 Continuous Integration

### GitHub Actions

```yaml
# .github/workflows/test.yml
name: Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    
    services:
      postgres:
        image: postgres:13
        env:
          POSTGRES_PASSWORD: postgres
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
    
    steps:
    - uses: actions/checkout@v2
    
    - name: Set up Go
      uses: actions/setup-go@v2
      with:
        go-version: 1.21
    
    - name: Install dependencies
      run: make install
    
    - name: Run unit tests
      run: make test-unit
    
    - name: Run integration tests
      run: make test-integration
      env:
        DATABASE_URL: postgres://postgres:postgres@localhost:5432/postgres
    
    - name: Upload coverage
      uses: codecov/codecov-action@v1
      with:
        file: ./coverage.out
```

## 🎯 Best Practices Summary

1. **Write tests first** - TDD helps with better design
2. **Test behavior** - Not implementation details
3. **Keep tests simple** - Easy to understand and maintain
4. **Use descriptive names** - Test names should explain what they test
5. **Arrange, Act, Assert** - Clear test structure
6. **Independent tests** - No dependencies between tests
7. **Mock external dependencies** - Keep unit tests isolated
8. **Test edge cases** - Boundary conditions and error paths
9. **Maintain test code** - Same quality standards as production code
10. **Run tests frequently** - Fast feedback loop

Remember: **Tests are documentation** - they show how your code is supposed to work!