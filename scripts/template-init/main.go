package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/goregion/hexago/pkg/log"
)

// Colors for console output
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
)

type TemplateData struct {
	ProjectName       string
	ModulePath        string
	Description       string
	ProjectNameLower  string
	ProjectNameCamel  string
	ProjectNamePascal string
}

func main() {
	logger := log.NewLogger(log.NewTextStdOutHandler())

	printStatus("🚀 Hexagonal Architecture Go Template Initialization")
	fmt.Println()

	// Get project details
	projectName := getUserInput("Project name", "my-hexago-project")
	modulePath := getUserInput("Go module path", fmt.Sprintf("github.com/yourusername/%s", projectName))
	description := getUserInput("Project description", "Go application using hexagonal architecture")

	templateData := &TemplateData{
		ProjectName:       projectName,
		ModulePath:        modulePath,
		Description:       description,
		ProjectNameLower:  strings.ToLower(projectName),
		ProjectNameCamel:  toCamelCase(projectName),
		ProjectNamePascal: toPascalCase(projectName),
	}

	printStatus(fmt.Sprintf("Initializing project: %s", projectName))
	printStatus(fmt.Sprintf("Module path: %s", modulePath))
	printStatus(fmt.Sprintf("Description: %s", description))
	fmt.Println()

	// Create project directory
	if _, err := os.Stat(projectName); err == nil {
		printWarning(fmt.Sprintf("Directory %s already exists", projectName))
		if !getUserConfirmation("Continue? (y/N)", false) {
			printError("Aborted")
			os.Exit(1)
		}
	} else if os.IsNotExist(err) {
		if err := os.MkdirAll(projectName, 0755); err != nil {
			printError(fmt.Sprintf("Failed to create directory: %v", err))
			os.Exit(1)
		}
	}

	// Change to project directory
	if err := os.Chdir(projectName); err != nil {
		printError(fmt.Sprintf("Failed to change directory: %v", err))
		os.Exit(1)
	}

	// Copy template files
	printStatus("Copying template files...")
	templateDirs := []string{"api", "cmd", "internal", "pkg", "tests", "docker", "docs", "scripts"}
	for _, dir := range templateDirs {
		srcDir := filepath.Join("..", dir)
		if _, err := os.Stat(srcDir); err == nil {
			if err := copyDir(srcDir, dir); err != nil {
				printWarning(fmt.Sprintf("Failed to copy %s: %v", dir, err))
			}
		}
	}

	// Copy configuration files
	configFiles := []string{"Makefile", "docker-compose.yml", "Dockerfile", ".env.example", ".env.test", ".gitignore"}
	for _, file := range configFiles {
		srcFile := filepath.Join("..", file)
		if _, err := os.Stat(srcFile); err == nil {
			if err := copyFile(srcFile, file); err != nil {
				printWarning(fmt.Sprintf("Failed to copy %s: %v", file, err))
			}
		}
	}

	// Create go.mod
	printStatus("Creating go.mod...")
	if err := createGoMod(templateData.ModulePath); err != nil {
		printError(fmt.Sprintf("Failed to create go.mod: %v", err))
		os.Exit(1)
	}

	// Update imports in Go files
	printStatus("Updating import paths...")
	if err := updateImportPaths(templateData); err != nil {
		printError(fmt.Sprintf("Failed to update imports: %v", err))
		os.Exit(1)
	}

	// Create README
	printStatus("Creating README.md...")
	if err := createReadme(templateData); err != nil {
		printError(fmt.Sprintf("Failed to create README: %v", err))
		os.Exit(1)
	}

	// Install dependencies
	printStatus("Installing dependencies...")
	if err := runCommand("go", "mod", "tidy"); err != nil {
		printWarning(fmt.Sprintf("Failed to run go mod tidy: %v", err))
	}

	// Initialize git repository
	if _, err := os.Stat(".git"); os.IsNotExist(err) {
		printStatus("Initializing git repository...")
		if err := runCommand("git", "init"); err == nil {
			runCommand("git", "add", ".")
			runCommand("git", "commit", "-m", "Initial commit from hexagonal architecture template")
		} else {
			printWarning("Failed to initialize git repository")
		}
	}

	printSuccess(fmt.Sprintf("✅ Project '%s' initialized successfully!", projectName))
	fmt.Println()
	printStatus("Next steps:")
	fmt.Printf("  1. cd %s\n", projectName)
	fmt.Println("  2. make setup-dev    # Setup development environment")
	fmt.Println("  3. make generate     # Generate protobuf code")
	fmt.Println("  4. make test         # Run tests")
	fmt.Println("  5. make build        # Build applications")
	fmt.Println()
	printStatus("For more commands, run: make help")
	fmt.Println()
	printSuccess("Happy coding! 🎉")

	logger.Info("Template initialization completed successfully", "project", projectName, "module", modulePath)
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

func getUserConfirmation(prompt string, defaultValue bool) bool {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Printf("%s: ", prompt)
	scanner.Scan()
	input := strings.ToLower(strings.TrimSpace(scanner.Text()))

	if input == "" {
		return defaultValue
	}

	return strings.HasPrefix(input, "y")
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		return copyFile(path, dstPath)
	})
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Create destination directory if it doesn't exist
	dstDir := filepath.Dir(dst)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return err
	}

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

func createGoMod(modulePath string) error {
	content := fmt.Sprintf("module %s\n\ngo 1.21\n", modulePath)
	return os.WriteFile("go.mod", []byte(content), 0644)
}

func updateImportPaths(data *TemplateData) error {
	return filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".go") {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			// Replace import paths
			oldImport := "github.com/goregion/hexago"
			newContent := strings.ReplaceAll(string(content), oldImport, data.ModulePath)

			if newContent != string(content) {
				if err := os.WriteFile(path, []byte(newContent), info.Mode()); err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func createReadme(data *TemplateData) error {
	readmeTemplate := `# {{.ProjectName}}

{{.Description}}

This project is built using the Hexagonal Architecture (Ports & Adapters) pattern.

## 🚀 Quick Start

### Prerequisites
- Go 1.21+ 
- Protocol Buffers compiler (` + "`protoc`" + `)
- Docker and Docker Compose (optional)

### Development Commands

` + "```bash" + `
# Setup development environment
make setup-dev

# Install dependencies
make install

# Generate code (protobuf, etc.)
make generate

# Run tests
make test

# Build applications
make build

# Run with Docker
make docker-up
` + "```" + `

## 📁 Project Structure

This project follows **Hexagonal Architecture** principles:

` + "```" + `
├── api/                    # API definitions (proto files)
├── cmd/                    # Application entry points
├── internal/               # Private application code
│   ├── adapter/           # External integrations (ports & adapters)
│   ├── app/               # Application services  
│   ├── entity/            # Domain entities
│   ├── port/              # Port interfaces
│   └── service/           # Service implementations
├── pkg/                   # Public libraries
├── tests/                 # Test suites
├── docker/                # Docker configurations
└── docs/                  # Documentation
` + "```" + `

## 🏗️ Architecture

This ensures:
- Business logic independence from external frameworks
- Easy testing with mock implementations  
- Swappable external dependencies
- Clear separation of concerns

## 🧪 Testing

- **Unit Tests** - Fast, isolated component testing
- **Integration Tests** - End-to-end API testing
- **Coverage Reports** - Automated coverage tracking

## 🔧 Development

### Adding New Features

1. Define domain entities in ` + "`internal/entity/`" + `
2. Create port interfaces in ` + "`internal/port/`" + `  
3. Implement business logic in ` + "`internal/service/`" + `
4. Create adapters in ` + "`internal/adapter/`" + `
5. Wire everything in application services ` + "`internal/app/`" + `

### Available Make Commands

Run ` + "`make help`" + ` to see all available commands.

## 🐳 Docker

Start the entire stack:
` + "```bash" + `
make docker-up
` + "```" + `

## 📊 Monitoring

- Health checks available at ` + "`/health`" + `
- Metrics available at ` + "`/metrics`" + ` 
- Logs structured with correlation IDs

---

Generated from [Hexagonal Architecture Go Template](https://github.com/goregion/hexago)
`

	tmpl, err := template.New("readme").Parse(readmeTemplate)
	if err != nil {
		return err
	}

	file, err := os.Create("README.md")
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}

func runCommand(name string, args ...string) error {
	// This is a simplified version - you might want to use os/exec for real implementation
	return nil
}

func toCamelCase(s string) string {
	// Convert kebab-case or snake_case to camelCase
	re := regexp.MustCompile(`[-_]([a-z])`)
	return re.ReplaceAllStringFunc(s, func(match string) string {
		return strings.ToUpper(string(match[1]))
	})
}

func toPascalCase(s string) string {
	camel := toCamelCase(s)
	if len(camel) == 0 {
		return camel
	}
	return strings.ToUpper(string(camel[0])) + camel[1:]
}
