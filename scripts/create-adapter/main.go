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

type AdapterData struct {
	AdapterName   string
	AdapterType   string
	PortName      string
	PackageName   string
	AdapterStruct string
	PortInterface string
	ModulePath    string
}

func main() {
	printStatus("🔌 Creating New Adapter")
	fmt.Println()

	// Get adapter details
	adapterName := getUserInput("Adapter name (e.g., postgres, kafka, http)", "")
	if adapterName == "" {
		printError("Adapter name is required")
		os.Exit(1)
	}

	adapterType := getUserInput("Adapter type", "driven")
	portName := getUserInput("Port interface name (e.g., UserRepository, EventPublisher)", "")
	if portName == "" {
		printError("Port interface name is required")
		os.Exit(1)
	}

	// Get module path from go.mod
	modulePath := getModulePath()
	if modulePath == "" {
		printError("Could not determine module path. Make sure you're in a Go project directory")
		os.Exit(1)
	}

	data := &AdapterData{
		AdapterName:   adapterName,
		AdapterType:   adapterType,
		PortName:      portName,
		PackageName:   strings.ReplaceAll(adapterName, "-", "_"),
		AdapterStruct: toPascalCase(adapterName) + "Adapter",
		PortInterface: portName,
		ModulePath:    modulePath,
	}

	adapterDir := filepath.Join("internal", "adapter", adapterName)

	printStatus(fmt.Sprintf("Creating adapter: %s", data.AdapterName))
	printStatus(fmt.Sprintf("Adapter struct: %s", data.AdapterStruct))
	printStatus(fmt.Sprintf("Port interface: %s", data.PortInterface))
	printStatus(fmt.Sprintf("Directory: %s", adapterDir))
	fmt.Println()

	// Create directory
	if err := os.MkdirAll(adapterDir, 0755); err != nil {
		printError(fmt.Sprintf("Failed to create directory: %v", err))
		os.Exit(1)
	}

	// Create adapter implementation
	if err := createAdapterFile(data, adapterDir); err != nil {
		printError(fmt.Sprintf("Failed to create adapter file: %v", err))
		os.Exit(1)
	}

	// Create config file
	if err := createConfigFile(data, adapterDir); err != nil {
		printError(fmt.Sprintf("Failed to create config file: %v", err))
		os.Exit(1)
	}

	// Create test file
	if err := createTestFile(data, adapterDir); err != nil {
		printError(fmt.Sprintf("Failed to create test file: %v", err))
		os.Exit(1)
	}

	// Create port interface if it doesn't exist
	portFile := filepath.Join("internal", "port", adapterName+".go")
	if _, err := os.Stat(portFile); os.IsNotExist(err) {
		printStatus(fmt.Sprintf("Creating port interface: %s", portFile))
		if err := createPortFile(data, portFile); err != nil {
			printError(fmt.Sprintf("Failed to create port file: %v", err))
			os.Exit(1)
		}
	}

	printSuccess(fmt.Sprintf("✅ Adapter '%s' created successfully!", data.AdapterName))
	fmt.Println()
	printStatus("Files created:")
	fmt.Printf("  - %s/adapter.go\n", adapterDir)
	fmt.Printf("  - %s/config.go\n", adapterDir)
	fmt.Printf("  - %s/adapter_test.go\n", adapterDir)
	if _, err := os.Stat(portFile); err == nil {
		fmt.Printf("  - %s\n", portFile)
	}
	fmt.Println()
	printStatus("Next steps:")
	fmt.Printf("  1. Define methods in internal/port/%s.go\n", adapterName)
	fmt.Printf("  2. Implement methods in %s/adapter.go\n", adapterDir)
	fmt.Printf("  3. Add configuration fields in %s/config.go\n", adapterDir)
	fmt.Printf("  4. Write tests in %s/adapter_test.go\n", adapterDir)
	fmt.Println("  5. Wire up the adapter in your application service")
	fmt.Println()
	printWarning("Don't forget to:")
	fmt.Println("  - Add necessary dependencies to go.mod")
	fmt.Println("  - Update your dependency injection/wiring")
	fmt.Println("  - Add integration tests if needed")
}

func createAdapterFile(data *AdapterData, dir string) error {
	adapterTemplate := `package {{.PackageName}}

import (
	"context"

	"{{.ModulePath}}/internal/port"
)

// {{.AdapterStruct}} implements the {{.PortInterface}} port
type {{.AdapterStruct}} struct {
	// Add your dependencies here
	// e.g., db *sql.DB, client *http.Client, etc.
}

// New{{.AdapterStruct}} creates a new instance of {{.AdapterStruct}}
func New{{.AdapterStruct}}(/* dependencies */) *{{.AdapterStruct}} {
	return &{{.AdapterStruct}}{
		// Initialize dependencies
	}
}

// Ensure {{.AdapterStruct}} implements {{.PortInterface}}
var _ port.{{.PortInterface}} = (*{{.AdapterStruct}})(nil)

// TODO: Implement the methods from {{.PortInterface}} interface
// Example:
// func (a *{{.AdapterStruct}}) MethodName(ctx context.Context, param string) error {
//     // Implementation here
//     return nil
// }
`

	tmpl, err := template.New("adapter").Parse(adapterTemplate)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(dir, "adapter.go"))
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}

func createConfigFile(data *AdapterData, dir string) error {
	configTemplate := `package {{.PackageName}}

// Config holds configuration for {{.AdapterStruct}}
type Config struct {
	// Add configuration fields here
	// e.g., Host string ` + "`env:\"HOST\" required:\"true\"`" + `
	// e.g., Port int    ` + "`env:\"PORT\" default:\"8080\"`" + `
	// e.g., Timeout time.Duration ` + "`env:\"TIMEOUT\" default:\"30s\"`" + `
}

// DefaultConfig returns default configuration
func DefaultConfig() Config {
	return Config{
		// Set default values
	}
}
`

	tmpl, err := template.New("config").Parse(configTemplate)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(dir, "config.go"))
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}

func createTestFile(data *AdapterData, dir string) error {
	testTemplate := `package {{.PackageName}}

import (
	"context"
	"testing"
)

func TestNew{{.AdapterStruct}}(t *testing.T) {
	adapter := New{{.AdapterStruct}}()
	if adapter == nil {
		t.Fatal("Expected adapter to be created")
	}
}

// TODO: Add more tests for your adapter methods
// Example:
// func Test{{.AdapterStruct}}_MethodName(t *testing.T) {
//     adapter := New{{.AdapterStruct}}()
//     
//     err := adapter.MethodName(context.Background(), "test")
//     if err != nil {
//         t.Errorf("Expected no error, got: %v", err)
//     }
// }
`

	tmpl, err := template.New("test").Parse(testTemplate)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(dir, "adapter_test.go"))
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}

func createPortFile(data *AdapterData, portFile string) error {
	portTemplate := `package port

import "context"

// {{.PortInterface}} defines the contract for {{.AdapterName}} operations
type {{.PortInterface}} interface {
	// TODO: Define your interface methods here
	// Example:
	// ProcessData(ctx context.Context, data string) error
	// GetStatus(ctx context.Context) (string, error)
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

	if scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())
		if input == "" && defaultValue != "" {
			return defaultValue
		}
		return input
	}

	// If scan fails, return default or empty
	if defaultValue != "" {
		return defaultValue
	}
	return ""
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
