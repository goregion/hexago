# Tests

This directory contains all test files that don't belong to specific packages. It includes integration tests, end-to-end tests, and test utilities that span multiple packages.

## What to create here:

### Integration Tests
- **API integration tests** - Test full HTTP API workflows
- **Database integration tests** - Test database operations with real DB
- **Service integration tests** - Test service interactions
- **External service integration tests** - Test third-party API integrations

### End-to-End Tests
- **Full workflow tests** - Complete user journey tests
- **Performance tests** - Load and stress testing
- **Contract tests** - API contract validation
- **Smoke tests** - Basic functionality validation after deployment

### Test Utilities
- **Test fixtures** - Shared test data and setup
- **Test helpers** - Common testing utilities
- **Mock factories** - Reusable mock objects
- **Test containers** - Docker container setup for tests

### Test Structure:
```
tests/
├── integration/    # Integration tests
│   ├── api/        # API integration tests
│   ├── database/   # Database integration tests
│   └── service/    # Service integration tests
├── e2e/            # End-to-end tests
│   ├── workflows/  # Complete user workflows
│   └── performance/ # Load and performance tests
├── fixtures/       # Test data and fixtures
│   ├── data/       # JSON/YAML test data
│   └── sql/        # Database test data
└── utils/          # Test utilities and helpers
    ├── containers/ # Test container setup
    ├── mocks/      # Shared mock objects
    └── helpers/    # Test helper functions
```

### Integration Test Example:
```go
// tests/integration/api/user_test.go
package api

import (
    "testing"
    "net/http/httptest"
    
    "github.com/goregion/hexago/internal/adapter/http"
    "github.com/goregion/hexago/tests/utils"
)

func TestCreateUser_Integration(t *testing.T) {
    // Setup test environment
    app := utils.SetupTestApp(t)
    server := httptest.NewServer(http.NewHandler(app))
    defer server.Close()
    
    // Test user creation workflow
    response := utils.PostJSON(t, server.URL+"/users", map[string]interface{}{
        "email": "test@example.com",
        "name":  "Test User",
    })
    
    utils.AssertStatusCode(t, response, 201)
    utils.AssertUserCreated(t, app.DB, "test@example.com")
}
```

### E2E Test Example:
```go
// tests/e2e/workflows/user_registration_test.go
package workflows

import (
    "testing"
    
    "github.com/goregion/hexago/tests/utils"
)

func TestUserRegistrationWorkflow(t *testing.T) {
    client := utils.NewTestClient(t)
    
    // Step 1: Register user
    user := client.RegisterUser("john@example.com", "password123")
    
    // Step 2: Verify email (mock)
    client.VerifyEmail(user.ID, user.VerificationToken)
    
    // Step 3: Login
    token := client.Login("john@example.com", "password123")
    
    // Step 4: Access protected resource
    profile := client.GetProfile(token)
    
    assert.Equal(t, "john@example.com", profile.Email)
    assert.True(t, profile.EmailVerified)
}
```

### Test Utilities Example:
```go
// tests/utils/helpers.go
package utils

import (
    "testing"
    "database/sql"
    
    "github.com/testcontainers/testcontainers-go"
    "github.com/goregion/hexago/pkg/launcher"
)

func SetupTestApp(t *testing.T) *launcher.App {
    // Setup test database container
    db := SetupTestDatabase(t)
    
    // Initialize app with test configuration
    app := launcher.NewApp(
        launcher.WithDatabase(db),
        launcher.WithTestMode(true),
    )
    
    return app
}

func SetupTestDatabase(t *testing.T) *sql.DB {
    // Start PostgreSQL container for testing
    container := testcontainers.StartPostgreSQLContainer(t)
    
    // Run migrations
    RunMigrations(t, container.ConnectionString())
    
    return container.DB()
}
```

## Key Principles:
- Tests should be independent and repeatable
- Use real dependencies for integration tests (database, external services)
- Clean up resources after each test
- Use test containers for isolated testing environments
- Separate fast unit tests from slower integration/e2e tests
- Maintain test data fixtures for consistent testing
- Mock external dependencies appropriately

## Running Tests:
```bash
# Run all tests
make test

# Run only integration tests
make test-integration

# Run only e2e tests  
make test-e2e

# Run tests with coverage
make test-coverage
```