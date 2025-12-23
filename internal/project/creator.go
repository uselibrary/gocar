package project

import (
	"fmt"
	"os"
	"path/filepath"

	"gocar/internal/util"
)

// Creator 项目创建器
type Creator struct {
	Name string
	Mode string
}

// NewCreator 创建项目创建器
func NewCreator(name, mode string) *Creator {
	return &Creator{
		Name: name,
		Mode: mode,
	}
}

// Create 创建项目
func (c *Creator) Create() error {
	// Check if directory already exists
	if _, err := os.Stat(c.Name); !os.IsNotExist(err) {
		return fmt.Errorf("directory '%s' already exists", c.Name)
	}

	var err error
	if c.Mode == "simple" {
		err = c.createSimpleProject()
	} else {
		err = c.createProjectMode()
	}

	if err != nil {
		return err
	}

	// Initialize git
	if err := util.InitGit(c.Name); err != nil {
		fmt.Printf("Warning: Failed to initialize git: %v\n", err)
	}

	return nil
}

// createSimpleProject 创建简单项目
func (c *Creator) createSimpleProject() error {
	// Create directories
	dirs := []string{
		c.Name,
		filepath.Join(c.Name, "bin"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Create go.mod
	if err := util.RunCommandSilent(c.Name, "go", "mod", "init", c.Name); err != nil {
		return fmt.Errorf("failed to initialize go.mod: %w", err)
	}

	// Create main.go
	if err := util.WriteFile(filepath.Join(c.Name, "main.go"), c.simpleMainTemplate()); err != nil {
		return err
	}

	// Create README.md
	if err := util.WriteFile(filepath.Join(c.Name, "README.md"), c.simpleReadmeTemplate()); err != nil {
		return err
	}

	// Create .gitignore
	if err := util.WriteFile(filepath.Join(c.Name, ".gitignore"), c.gitignoreTemplate()); err != nil {
		return err
	}

	return nil
}

// createProjectMode 创建项目模式
func (c *Creator) createProjectMode() error {
	// Create directories
	dirs := []string{
		c.Name,
		filepath.Join(c.Name, "cmd", c.Name),
		filepath.Join(c.Name, "internal"),
		filepath.Join(c.Name, "pkg"),
		filepath.Join(c.Name, "test"),
		filepath.Join(c.Name, "bin"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Create go.mod
	if err := util.RunCommandSilent(c.Name, "go", "mod", "init", c.Name); err != nil {
		return fmt.Errorf("failed to initialize go.mod: %w", err)
	}

	// Create cmd/<appName>/main.go
	if err := util.WriteFile(filepath.Join(c.Name, "cmd", c.Name, "main.go"), c.projectMainTemplate()); err != nil {
		return err
	}

	// Create .gitkeep files for empty directories
	emptyDirs := []string{
		filepath.Join(c.Name, "internal", ".gitkeep"),
		filepath.Join(c.Name, "pkg", ".gitkeep"),
		filepath.Join(c.Name, "test", ".gitkeep"),
	}
	for _, f := range emptyDirs {
		if err := util.WriteFile(f, ""); err != nil {
			return err
		}
	}

	// Create README.md
	if err := util.WriteFile(filepath.Join(c.Name, "README.md"), c.projectReadmeTemplate()); err != nil {
		return err
	}

	// Create .gitignore
	if err := util.WriteFile(filepath.Join(c.Name, ".gitignore"), c.gitignoreTemplate()); err != nil {
		return err
	}

	return nil
}

// simpleMainTemplate 生成简单项目的 main.go 内容
func (c *Creator) simpleMainTemplate() string {
	return fmt.Sprintf(`package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Hello, gocar! A golang project scaffolding tool for %s.")
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
}
`, c.Name)
}

// projectMainTemplate 生成项目模式的 main.go 内容
func (c *Creator) projectMainTemplate() string {
	return fmt.Sprintf(`package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Hello, gocar! A golang project scaffolding tool for %s.")
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
}
`, c.Name)
}

// simpleReadmeTemplate 生成简单项目的 README.md 内容
func (c *Creator) simpleReadmeTemplate() string {
	return fmt.Sprintf(`# %s

A Go project created with gocar.

## Build

`+"```bash"+`
# Debug build (current platform)
gocar build

# Release build (current platform)
gocar build --release

# Cross-compile for Linux on AMD64
gocar build --target linux/amd64
`+"```"+`

## Run

`+"```bash"+`
gocar run
`+"```"+`

## Output Structure

`+"```"+`
bin/
├── debug/
│   └── <os>-<arch>/
│       └── %s
└── release/
    └── <os>-<arch>/
        └── %s
`+"```"+`

Build artifacts are organized by:
- **Build mode**: debug or release
- **Target platform**: OS and architecture (e.g., linux-amd64, darwin-arm64)

Examples:
- Debug build for current platform: `+"`./bin/debug/linux-amd64/%s`"+`
- Release build for Windows: `+"`./bin/release/windows-amd64/%s.exe`"+`
`, c.Name, c.Name, c.Name, c.Name, c.Name)
}

// projectReadmeTemplate 生成项目模式的 README.md 内容
func (c *Creator) projectReadmeTemplate() string {
	return fmt.Sprintf(`# %s

A Go project created with gocar (project mode).

## Project Structure

`+"```"+`
%s/
├── cmd/
│   └── %s/
│       └── main.go      # Application entry point
├── internal/            # Private application code
├── pkg/                 # Public library code
├── test/                # Integration tests
├── bin/                 # Build output
├── go.mod
└── README.md
`+"```"+`

## Build

`+"```bash"+`
# Debug build (current platform)
gocar build

# Release build (current platform)
gocar build --release

# Cross-compile for Linux
gocar build --target linux/amd64
`+"```"+`

## Run

`+"```bash"+`
gocar run
`+"```"+`

## Output Structure

`+"```"+`
bin/
├── debug/
│   └── <os>-<arch>/
│       └── %s
└── release/
    └── <os>-<arch>/
        └── %s
`+"```"+`

Build artifacts are organized by:
- **Build mode**: debug or release
- **Target platform**: OS and architecture (e.g., linux-amd64, darwin-arm64)

Examples:
- Debug build for current platform: `+"`./bin/debug/linux-amd64/%s`"+`
- Release build for Windows: `+"`./bin/release/windows-amd64/%s.exe`"+`

## Directories

- **cmd/**: Main applications for this project
- **internal/**: Private application and library code (not importable by other projects)
- **pkg/**: Library code that can be used by external applications
- **test/**: Integration tests, black-box tests
`, c.Name, c.Name, c.Name, c.Name, c.Name, c.Name, c.Name)
}

// gitignoreTemplate 生成 .gitignore 内容
func (c *Creator) gitignoreTemplate() string {
	return fmt.Sprintf(`# Binaries
%s
bin/
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary
*.test

# Output of go coverage
*.out

# Dependency directories
vendor/

# IDE
.idea/
.vscode/
*.swp
*.swo

# OS files
.DS_Store
Thumbs.db
`, c.Name)
}
