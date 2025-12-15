package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

const version = "0.1.0"

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(0)
	}

	command := os.Args[1]

	switch command {
	case "new":
		handleNew(os.Args[2:])
	case "build":
		handleBuild(os.Args[2:])
	case "run":
		handleRun(os.Args[2:])
	case "clean":
		handleClean()
	case "help", "-h", "--help":
		printHelp()
	case "version", "-v", "--version":
		fmt.Printf("gocar %s\n", version)
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printHelp()
		os.Exit(1)
	}
}

func printHelp() {
	help := `gocar - A cargo-like tool for Go projects

USAGE:
    gocar <COMMAND> [OPTIONS]

COMMANDS:
    new <name> [--mode simple|project]     Create a new Go project
    build [--release]                      Build the project
    run [args...]                          Run the project
    clean                                  Clean build artifacts
    help                                   Print this help message
    version                                Print version info

EXAMPLES:
    gocar new myapp                        Create a simple project
    gocar new myapp --mode project         Create a project-mode project
    gocar build                            Build in debug mode
    gocar build --release                  Build in release mode
    gocar run                              Build and run the project
`
	fmt.Print(help)
}

// ==================== NEW COMMAND ====================

func handleNew(args []string) {
	if len(args) < 1 {
		fmt.Println("Error: Missing project name")
		fmt.Println("Usage: gocar new <name> [--mode simple|project]")
		os.Exit(1)
	}

	appName := args[0]

	// Validate project name
	if err := validateProjectName(appName); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	mode := "simple" // default mode

	// Parse --mode flag
	for i := 1; i < len(args); i++ {
		if args[i] == "--mode" && i+1 < len(args) {
			mode = args[i+1]
			if mode != "simple" && mode != "project" {
				fmt.Printf("Error: Invalid mode '%s'. Use 'simple' or 'project'\n", mode)
				os.Exit(1)
			}
			break
		}
	}

	// Check if directory already exists
	if _, err := os.Stat(appName); !os.IsNotExist(err) {
		fmt.Printf("Error: Directory '%s' already exists\n", appName)
		os.Exit(1)
	}

	fmt.Printf("Creating new %s project: %s\n", mode, appName)

	var err error
	if mode == "simple" {
		err = createSimpleProject(appName)
	} else {
		err = createProjectMode(appName)
	}

	if err != nil {
		fmt.Printf("Error creating project: %v\n", err)
		os.Exit(1)
	}

	// Initialize git
	if err := initGit(appName); err != nil {
		fmt.Printf("Warning: Failed to initialize git: %v\n", err)
	}

	fmt.Printf("\nSuccessfully created project '%s'\n", appName)
	fmt.Printf("\nTo get started:\n")
	fmt.Printf("    cd %s\n", appName)
	fmt.Printf("    gocar build\n")
	fmt.Printf("    gocar run\n")
}

func createSimpleProject(appName string) error {
	// Create directories
	dirs := []string{
		appName,
		filepath.Join(appName, "bin"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Create go.mod
	if err := runCommandSilent(appName, "go", "mod", "init", appName); err != nil {
		return fmt.Errorf("failed to initialize go.mod: %w", err)
	}

	// Create main.go
	mainGo := `package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Hello, gocar! A golang package manager.")
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
}
`
	if err := writeFile(filepath.Join(appName, "main.go"), mainGo); err != nil {
		return err
	}

	// Create README.md
	readme := fmt.Sprintf(`# %s

A Go project created with gocar.

## Build

`+"```bash"+`
# Debug build
gocar build

# Release build
gocar build --release
`+"```"+`

## Run

`+"```bash"+`
gocar run
`+"```"+`

## Output

- Debug build: `+"`./bin/%s`"+`
- Release build: `+"`./bin/%s`"+` (with release flags: -ldflags="-s -w" -trimpath)
`, appName, appName, appName)

	if err := writeFile(filepath.Join(appName, "README.md"), readme); err != nil {
		return err
	}

	// Create .gitignore
	gitignore := fmt.Sprintf(`# Binaries
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
`, appName)

	if err := writeFile(filepath.Join(appName, ".gitignore"), gitignore); err != nil {
		return err
	}

	return nil
}

func createProjectMode(appName string) error {
	// Create directories
	dirs := []string{
		appName,
		filepath.Join(appName, "cmd", "server"),
		filepath.Join(appName, "internal"),
		filepath.Join(appName, "pkg"),
		filepath.Join(appName, "test"),
		filepath.Join(appName, "bin"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Create go.mod
	if err := runCommandSilent(appName, "go", "mod", "init", appName); err != nil {
		return fmt.Errorf("failed to initialize go.mod: %w", err)
	}

	// Create cmd/server/main.go
	mainGo := `package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Hello, gocar! A golang package manager.")
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
}
`
	if err := writeFile(filepath.Join(appName, "cmd", "server", "main.go"), mainGo); err != nil {
		return err
	}

	// Create .gitkeep files for empty directories
	emptyDirs := []string{
		filepath.Join(appName, "internal", ".gitkeep"),
		filepath.Join(appName, "pkg", ".gitkeep"),
		filepath.Join(appName, "test", ".gitkeep"),
	}
	for _, f := range emptyDirs {
		if err := writeFile(f, ""); err != nil {
			return err
		}
	}

	// Create README.md
	readme := fmt.Sprintf(`# %s

A Go project created with gocar (project mode).

## Project Structure

`+"```"+`
%s/
├── cmd/
│   └── server/
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
# Debug build
gocar build

# Release build
gocar build --release
`+"```"+`

## Run

`+"```bash"+`
gocar run
`+"```"+`

## Output

- Debug build: `+"`./bin/%s`"+`
- Release build: `+"`./bin/%s`"+` (with release flags: CGO_ENABLED=0 -ldflags="-s -w" -trimpath)

## Directories

- **cmd/**: Main applications for this project
- **internal/**: Private application and library code (not importable by other projects)
- **pkg/**: Library code that can be used by external applications
- **test/**: Integration tests, black-box tests
`, appName, appName, appName, appName)

	if err := writeFile(filepath.Join(appName, "README.md"), readme); err != nil {
		return err
	}

	// Create .gitignore
	gitignore := fmt.Sprintf(`# Binaries
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
`, appName)

	if err := writeFile(filepath.Join(appName, ".gitignore"), gitignore); err != nil {
		return err
	}

	return nil
}

func initGit(appName string) error {
	// git init with main as default branch
	if err := runCommandSilent(appName, "git", "init", "-b", "main"); err != nil {
		return err
	}

	// git add .
	if err := runCommandSilent(appName, "git", "add", "."); err != nil {
		return err
	}

	return nil
}

// ==================== BUILD COMMAND ====================

func handleBuild(args []string) {
	release := false
	for _, arg := range args {
		if arg == "--release" {
			release = true
			break
		}
	}

	// Get project info
	projectRoot, appName, projectMode, err := detectProject()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	outputPath := filepath.Join("bin", appName)
	if runtime.GOOS == "windows" {
		outputPath += ".exe"
	}

	var buildArgs []string
	var env []string

	if release {
		fmt.Println("Building in release mode...")
		env = append(os.Environ(), "CGO_ENABLED=0")
		buildArgs = []string{"build", "-ldflags=-s -w", "-trimpath", "-o", outputPath}
	} else {
		fmt.Println("Building in debug mode...")
		env = os.Environ()
		buildArgs = []string{"build", "-o", outputPath}
	}

	// Determine source path based on project mode
	if projectMode == "project" {
		buildArgs = append(buildArgs, "./cmd/server/main.go")
	} else {
		buildArgs = append(buildArgs, "./main.go")
	}

	// Ensure bin directory exists
	binDir := filepath.Join(projectRoot, "bin")
	if err := os.MkdirAll(binDir, 0755); err != nil {
		fmt.Printf("Error creating bin directory: %v\n", err)
		os.Exit(1)
	}

	cmd := exec.Command("go", buildArgs...)
	cmd.Dir = projectRoot
	cmd.Env = env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Build failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Build successful: %s\n", filepath.Join(projectRoot, outputPath))
}

// ==================== RUN COMMAND ====================

func handleRun(args []string) {
	// Get project info
	projectRoot, appName, projectMode, err := detectProject()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Determine source path based on project mode
	var sourcePath string
	if projectMode == "project" {
		sourcePath = "./cmd/server/main.go"
	} else {
		sourcePath = "./main.go"
	}

	fmt.Printf("Running %s...\n\n", appName)

	runArgs := []string{"run", sourcePath}
	runArgs = append(runArgs, args...)

	cmd := exec.Command("go", runArgs...)
	cmd.Dir = projectRoot
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		// Don't print error for normal exit
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		fmt.Printf("Run failed: %v\n", err)
		os.Exit(1)
	}
}

// ==================== CLEAN COMMAND ====================

func handleClean() {
	projectRoot, appName, _, err := detectProject()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	binDir := filepath.Join(projectRoot, "bin")

	// Remove bin directory contents
	entries, err := os.ReadDir(binDir)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Nothing to clean.")
			return
		}
		fmt.Printf("Error reading bin directory: %v\n", err)
		os.Exit(1)
	}

	if len(entries) == 0 {
		fmt.Println("Nothing to clean.")
		return
	}

	for _, entry := range entries {
		path := filepath.Join(binDir, entry.Name())
		if err := os.RemoveAll(path); err != nil {
			fmt.Printf("Error removing %s: %v\n", path, err)
		}
	}

	fmt.Printf("Cleaned build artifacts for '%s'\n", appName)
}

// ==================== HELPER FUNCTIONS ====================

func validateProjectName(name string) error {
	// Check if empty
	if name == "" {
		return fmt.Errorf("project name cannot be empty")
	}

	// Check if starts with a dash or dot
	if strings.HasPrefix(name, "-") || strings.HasPrefix(name, ".") {
		return fmt.Errorf("project name cannot start with '-' or '.'")
	}

	// Check for valid characters (alphanumeric, dash, underscore)
	validName := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*$`)
	if !validName.MatchString(name) {
		return fmt.Errorf("project name must start with a letter and contain only letters, numbers, dashes, or underscores")
	}

	// Check for reserved names
	reserved := []string{"test", "main", "init", "internal", "vendor"}
	for _, r := range reserved {
		if strings.ToLower(name) == r {
			return fmt.Errorf("'%s' is a reserved name in Go", name)
		}
	}

	return nil
}

func detectProject() (projectRoot, appName, projectMode string, err error) {
	// Find project root by looking for go.mod
	cwd, err := os.Getwd()
	if err != nil {
		return "", "", "", fmt.Errorf("failed to get current directory: %w", err)
	}

	projectRoot = cwd
	for {
		if _, err := os.Stat(filepath.Join(projectRoot, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(projectRoot)
		if parent == projectRoot {
			return "", "", "", fmt.Errorf("not in a Go module (go.mod not found)")
		}
		projectRoot = parent
	}

	// Get app name from directory name
	appName = filepath.Base(projectRoot)

	// Detect project mode
	if _, err := os.Stat(filepath.Join(projectRoot, "cmd", "server", "main.go")); err == nil {
		projectMode = "project"
	} else if _, err := os.Stat(filepath.Join(projectRoot, "main.go")); err == nil {
		projectMode = "simple"
	} else {
		return "", "", "", fmt.Errorf("cannot detect project mode: no main.go or cmd/server/main.go found")
	}

	return projectRoot, appName, projectMode, nil
}

func writeFile(path, content string) error {
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", path, err)
	}
	return nil
}

func runCommandSilent(dir string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	return cmd.Run()
}
