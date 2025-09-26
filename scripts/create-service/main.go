package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// Colors for console output
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
)

type ServiceData struct {
	ServiceName      string
	ServiceType      string
	ModulePath       string
	PackageName      string
	ServiceStruct    string
	ServiceInterface string
	AppServiceStruct string
}

func main() {
	printStatus("⚙️ Creating New Service")
	fmt.Println()

	// Get service details
	serviceName := getUserInput("Service name (e.g., user, notification, payment)", "")
	if serviceName == "" {
		printError("Service name is required")
		os.Exit(1)
	}

	serviceType := getUserInput("Service type", "domain")

	// Get module path from go.mod
	modulePath := getModulePath()
	if modulePath == "" {
		printError("Could not determine module path. Make sure you're in a Go project directory")
		os.Exit(1)
	}

	data := &ServiceData{
		ServiceName:      serviceName,
		ServiceType:      serviceType,
		ModulePath:       modulePath,
		PackageName:      strings.ReplaceAll(serviceName, "-", "_"),
		ServiceStruct:    toPascalCase(serviceName) + "Service",
		ServiceInterface: toPascalCase(serviceName) + "Service",
		AppServiceStruct: toPascalCase(serviceName) + "ApplicationService",
	}

	serviceDir := filepath.Join("internal", "service", serviceName)
	appServiceDir := filepath.Join("internal", "app", serviceName)

	printStatus(fmt.Sprintf("Creating service: %s", data.ServiceName))
	printStatus(fmt.Sprintf("Domain service: %s", data.ServiceStruct))
	printStatus(fmt.Sprintf("App service: %s", data.AppServiceStruct))
	printStatus(fmt.Sprintf("Service dir: %s", serviceDir))
	printStatus(fmt.Sprintf("App service dir: %s", appServiceDir))
	fmt.Println()

	// Create directories
	if err := os.MkdirAll(serviceDir, 0755); err != nil {
		printError(fmt.Sprintf("Failed to create service directory: %v", err))
		os.Exit(1)
	}

	if err := os.MkdirAll(appServiceDir, 0755); err != nil {
		printError(fmt.Sprintf("Failed to create app service directory: %v", err))
		os.Exit(1)
	}

	// Create domain service
	if err := createDomainServiceFile(data, serviceDir); err != nil {
		printError(fmt.Sprintf("Failed to create domain service file: %v", err))
		os.Exit(1)
	}

	// Create application service
	if err := createAppServiceFile(data, appServiceDir); err != nil {
		printError(fmt.Sprintf("Failed to create application service file: %v", err))
		os.Exit(1)
	}

	// Create test files
	if err := createDomainServiceTestFile(data, serviceDir); err != nil {
		printError(fmt.Sprintf("Failed to create domain service test file: %v", err))
		os.Exit(1)
	}

	if err := createAppServiceTestFile(data, appServiceDir); err != nil {
		printError(fmt.Sprintf("Failed to create application service test file: %v", err))
		os.Exit(1)
	}

	// Create service interface in port if needed
	portFile := filepath.Join("internal", "port", serviceName+"-service.go")
	if _, err := os.Stat(portFile); os.IsNotExist(err) {
		printStatus(fmt.Sprintf("Creating service interface: %s", portFile))
		if err := createServicePortFile(data, portFile); err != nil {
			printError(fmt.Sprintf("Failed to create service port file: %v", err))
			os.Exit(1)
		}
	}

	printSuccess(fmt.Sprintf("✅ Service '%s' created successfully!", data.ServiceName))
	fmt.Println()
	printStatus("Files created:")
	fmt.Printf("  - %s/%s.go (domain service)\n", serviceDir, serviceName)
	fmt.Printf("  - %s/%s_test.go (domain service test)\n", serviceDir, serviceName)
	fmt.Printf("  - %s/service.go (application service)\n", appServiceDir)
	fmt.Printf("  - %s/service_test.go (application service test)\n", appServiceDir)
	if _, err := os.Stat(portFile); err == nil {
		fmt.Printf("  - %s\n", portFile)
	}
	fmt.Println()
	printStatus("Next steps:")
	fmt.Printf("  1. Define business logic in %s/%s.go\n", serviceDir, serviceName)
	fmt.Printf("  2. Implement use cases in %s/service.go\n", appServiceDir)
	fmt.Printf("  3. Add necessary ports/adapters dependencies\n")
	fmt.Printf("  4. Write comprehensive tests\n")
	fmt.Println("  5. Wire up the services in your main application")
	fmt.Println()
	printWarning("Don't forget to:")
	fmt.Println("  - Add business entities if needed")
	fmt.Println("  - Define repository interfaces in ports")
	fmt.Println("  - Update your dependency injection/wiring")
	fmt.Println("  - Add integration tests")
}

func createDomainServiceFile(data *ServiceData, dir string) error {
	serviceTemplate := `package {{.PackageName}}

import (
	"context"
)

// {{.ServiceStruct}} represents domain service for {{.ServiceName}}
type {{.ServiceStruct}} struct {
	// Add your domain service dependencies here
	// Usually repositories, other domain services, etc.
}

// New{{.ServiceStruct}} creates a new instance of {{.ServiceStruct}}
func New{{.ServiceStruct}}(/* dependencies */) *{{.ServiceStruct}} {
	return &{{.ServiceStruct}}{
		// Initialize dependencies
	}
}

// TODO: Implement your domain business logic methods here
// Example:
// func (s *{{.ServiceStruct}}) ProcessBusinessLogic(ctx context.Context, data string) error {
//     // Domain business rules and logic
//     return nil
// }

// func (s *{{.ServiceStruct}}) ValidateBusinessRules(ctx context.Context, entity *Entity) error {
//     // Business validation logic
//     return nil
// }
`

	tmpl, err := template.New("service").Parse(serviceTemplate)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(dir, data.ServiceName+".go"))
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}

func createAppServiceFile(data *ServiceData, dir string) error {
	appServiceTemplate := `package {{.ServiceName}}

import (
	"context"

	"{{.ModulePath}}/internal/service/{{.ServiceName}}"
)

// {{.AppServiceStruct}} orchestrates use cases for {{.ServiceName}}
type {{.AppServiceStruct}} struct {
	// Add your ports (interfaces) here
	// e.g., userRepo port.UserRepository
	// e.g., eventPublisher port.EventPublisher
	// e.g., logger port.Logger
	
	// Domain service
	domainService *{{.PackageName}}.{{.ServiceStruct}}
}

// New{{.AppServiceStruct}} creates a new instance of {{.AppServiceStruct}}
func New{{.AppServiceStruct}}(
	// Add your dependencies as parameters
	domainService *{{.PackageName}}.{{.ServiceStruct}},
) *{{.AppServiceStruct}} {
	return &{{.AppServiceStruct}}{
		// Initialize ports
		domainService: domainService,
	}
}

// TODO: Implement your use case methods here
// Example:
// func (s *{{.AppServiceStruct}}) CreateUser(ctx context.Context, req CreateUserRequest) (*User, error) {
//     // 1. Validate request
//     if err := s.validateCreateUserRequest(req); err != nil {
//         return nil, err
//     }
//     
//     // 2. Use domain service for business logic
//     user, err := s.domainService.CreateUser(ctx, req.Name, req.Email)
//     if err != nil {
//         return nil, err
//     }
//     
//     // 3. Persist using repository
//     if err := s.userRepo.Save(ctx, user); err != nil {
//         return nil, err
//     }
//     
//     // 4. Publish events if needed
//     s.eventPublisher.Publish(ctx, NewUserCreatedEvent(user))
//     
//     return user, nil
// }

// Request/Response DTOs
// TODO: Define your request/response structures
// type CreateUserRequest struct {
//     Name  string ` + "`json:\"name\" validate:\"required\"`" + `
//     Email string ` + "`json:\"email\" validate:\"required,email\"`" + `
// }
`

	tmpl, err := template.New("appservice").Parse(appServiceTemplate)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(dir, "service.go"))
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}

func createDomainServiceTestFile(data *ServiceData, dir string) error {
	testTemplate := `package {{.PackageName}}

import (
	"context"
	"testing"
)

func TestNew{{.ServiceStruct}}(t *testing.T) {
	service := New{{.ServiceStruct}}()
	if service == nil {
		t.Fatal("Expected service to be created")
	}
}

// TODO: Add more tests for your domain service methods
// Example:
// func Test{{.ServiceStruct}}_BusinessLogic(t *testing.T) {
//     service := New{{.ServiceStruct}}()
//     
//     err := service.ProcessBusinessLogic(context.Background(), "test")
//     if err != nil {
//         t.Errorf("Expected no error, got: %v", err)
//     }
// }
`

	tmpl, err := template.New("test").Parse(testTemplate)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(dir, data.ServiceName+"_test.go"))
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}

func createAppServiceTestFile(data *ServiceData, dir string) error {
	appTestTemplate := `package {{.ServiceName}}

import (
	"context"
	"testing"

	"{{.ModulePath}}/internal/service/{{.ServiceName}}"
)

func TestNew{{.AppServiceStruct}}(t *testing.T) {
	domainService := {{.PackageName}}.New{{.ServiceStruct}}()
	appService := New{{.AppServiceStruct}}(domainService)
	
	if appService == nil {
		t.Fatal("Expected application service to be created")
	}
}

// TODO: Add more tests for your application service use cases
// Example:
// func Test{{.AppServiceStruct}}_UseCase(t *testing.T) {
//     // Arrange
//     mockRepo := &MockUserRepository{}
//     domainService := {{.PackageName}}.New{{.ServiceStruct}}(mockRepo)
//     appService := New{{.AppServiceStruct}}(domainService)
//     
//     // Act
//     result, err := appService.SomeUseCase(context.Background(), request)
//     
//     // Assert
//     if err != nil {
//         t.Errorf("Expected no error, got: %v", err)
//     }
//     // Add more assertions...
// }
`

	tmpl, err := template.New("apptest").Parse(appTestTemplate)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(dir, "service_test.go"))
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}

func createServicePortFile(data *ServiceData, portFile string) error {
	portTemplate := `package port

import "context"

// {{.ServiceInterface}} defines the contract for {{.ServiceName}} use cases
type {{.ServiceInterface}} interface {
	// TODO: Define your use case methods here
	// Example:
	// CreateUser(ctx context.Context, req CreateUserRequest) (*User, error)
	// GetUser(ctx context.Context, userID string) (*User, error)
	// UpdateUser(ctx context.Context, userID string, req UpdateUserRequest) (*User, error)
	// DeleteUser(ctx context.Context, userID string) error
}
`

	// Create port directory if it doesn't exist
	portDir := filepath.Dir(portFile)
	if err := os.MkdirAll(portDir, 0755); err != nil {
		return err
	}

	tmpl, err := template.New("port").Parse(portTemplate)
	if err != nil {
		return err
	}

	file, err := os.Create(portFile)
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}

func printStatus(msg string) {
	fmt.Printf("%s[INFO]%s %s\n", ColorBlue, ColorReset, msg)
}

func printSuccess(msg string) {
	fmt.Printf("%s[SUCCESS]%s %s\n", ColorGreen, ColorReset, msg)
}

func printWarning(msg string) {
	fmt.Printf("%s[WARNING]%s %s\n", ColorYellow, ColorReset, msg)
}

func printError(msg string) {
	fmt.Printf("%s[ERROR]%s %s\n", ColorRed, ColorReset, msg)
}

func getUserInput(prompt, defaultValue string) string {
	scanner := bufio.NewScanner(os.Stdin)

	if defaultValue != "" {
		fmt.Printf("%s [%s]: ", prompt, defaultValue)
	} else {
		fmt.Printf("%s: ", prompt)
	}

	scanner.Scan()
	input := strings.TrimSpace(scanner.Text())

	if input == "" && defaultValue != "" {
		return defaultValue
	}

	return input
}

func getModulePath() string {
	content, err := os.ReadFile("go.mod")
	if err != nil {
		return ""
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module"))
		}
	}

	return ""
}

func toPascalCase(s string) string {
	parts := strings.FieldsFunc(s, func(c rune) bool {
		return c == '-' || c == '_' || c == ' '
	})

	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(string(part[0])) + strings.ToLower(part[1:])
		}
	}

	return strings.Join(parts, "")
}
