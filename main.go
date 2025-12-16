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

const version = "0.1.3"

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
	case "add":
		handleAdd(os.Args[2:])
	case "update":
		handleUpdate(os.Args[2:])
	case "tidy":
		handleTidy()
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
    add <package>...                       Add dependencies to go.mod
    update [package]...                    Update dependencies
    tidy                                   Tidy up go.mod and go.sum
    help                                   Print this help message
    version                                Print version info

EXAMPLES:
    gocar new myapp                        Create a simple project
    gocar new myapp --mode project         Create a project-mode project
    gocar build                            Build in debug mode
    gocar build --release                  Build in release mode
    gocar run                              Build and run the project
    gocar add github.com/gin-gonic/gin     Add a dependency
    gocar update                           Update all dependencies
    gocar tidy                             Clean up go.mod
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

	// Check for help
	if args[0] == "help" || args[0] == "--help" || args[0] == "-h" {
		fmt.Println("gocar new - Create a new Go project")
		fmt.Println()
		fmt.Println("USAGE:")
		fmt.Println("    gocar new <name> [--mode simple|project]")
		fmt.Println()
		fmt.Println("OPTIONS:")
		fmt.Println("    --mode <mode>    Project mode: 'simple' (default) or 'project'")
		fmt.Println()
		fmt.Println("EXAMPLES:")
		fmt.Println("    gocar new myapp              Create a simple project")
		fmt.Println("    gocar new myapp --mode project    Create a project-mode project")
		os.Exit(0)
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
		switch args[i] {
		case "--mode":
			if i+1 < len(args) {
				mode = args[i+1]
				if mode != "simple" && mode != "project" {
					fmt.Printf("Error: Invalid mode '%s'. Use 'simple' or 'project'\n", mode)
					os.Exit(1)
				}
				i++ // skip next arg
			} else {
				fmt.Println("Error: --mode requires a value")
				os.Exit(1)
			}
		default:
			if strings.HasPrefix(args[i], "-") {
				fmt.Printf("Error: Unknown option '%s'\n", args[i])
				fmt.Println("Run 'gocar new --help' for usage.")
				os.Exit(1)
			}
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
	mainGo := fmt.Sprintf(`package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Hello, gocar! A golang project scaffolding tool for %s.")
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
}
`, appName)
	if err := writeFile(filepath.Join(appName, "main.go"), mainGo); err != nil {
		return err
	}

	// Create README.md
	readme := fmt.Sprintf(`# %s

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
`, appName, appName, appName, appName, appName)

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
	mainGo := fmt.Sprintf(`package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Hello, gocar! A golang project scaffolding tool for %s.")
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
}
`, appName)
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
`, appName, appName, appName, appName, appName, appName)

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

func printBuildHelp() {
	help := `gocar build - Build the project

USAGE:
    gocar build [OPTIONS]

OPTIONS:
    --release              Build in release mode (optimized binary)
    --target <os>/<arch>   Cross-compile for target platform
    --with-cgo             Force enable CGO (sets CGO_ENABLED=1)
    --help                 Show this help message

EXAMPLES:
    gocar build                                  Build for current platform (debug)
    gocar build --release                        Build for current platform (release)
    gocar build --target linux/amd64             Cross-compile for Linux AMD64
    gocar build --target windows/amd64           Cross-compile for Windows AMD64
    gocar build --release --target linux/arm64   Cross-compile for Linux ARM (release)
    gocar build --with-cgo                       Build with CGO enabled
    gocar build --release --with-cgo             Build in release mode with CGO enabled

COMMON TARGETS:
    linux/amd64     Linux AMD 64-bit
    linux/arm64     Linux ARM 64-bit
    linux/arm       Linux ARM 32-bit
    darwin/amd64    macOS Intel
    darwin/arm64    macOS Apple Silicon
    windows/amd64   Windows 64-bit
    windows/386     Windows 32-bit
`
	fmt.Print(help)
}

func handleBuild(args []string) {
	release := false
	target := ""
	withCgo := false

	// Parse arguments
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "help", "--help", "-h":
			printBuildHelp()
			os.Exit(0)
		case "--release":
			release = true
		case "--with-cgo":
			withCgo = true
		case "--target":
			if i+1 < len(args) {
				target = args[i+1]
				i++ // skip next arg
			} else {
				fmt.Println("Error: --target requires a value")
				fmt.Println("Usage: gocar build --target <os>/<arch>")
				fmt.Println("Example: gocar build --target linux/amd64")
				os.Exit(1)
			}
		default:
			fmt.Printf("Error: Unknown option '%s'\n", arg)
			fmt.Println("Run 'gocar build --help' for usage.")
			os.Exit(1)
		}
	}

	// Get project info
	projectRoot, appName, projectMode, err := detectProject()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Parse target if specified
	targetOS := runtime.GOOS
	targetArch := runtime.GOARCH
	if target != "" {
		parts := strings.Split(target, "/")
		if len(parts) != 2 {
			fmt.Printf("Error: Invalid target format '%s'\n", target)
			fmt.Println("Expected format: <os>/<arch>")
			fmt.Println("Example: linux/amd64, windows/amd64, darwin/arm64")
			os.Exit(1)
		}
		targetOS = parts[0]
		targetArch = parts[1]
	}

	// Determine output path following Cargo's structure
	// bin/debug/<os>-<arch>/appname or bin/release/<os>-<arch>/appname
	buildMode := "debug"
	if release {
		buildMode = "release"
	}
	targetDir := fmt.Sprintf("%s-%s", targetOS, targetArch)
	outputDir := filepath.Join("bin", buildMode, targetDir)
	outputPath := filepath.Join(outputDir, appName)
	if targetOS == "windows" {
		outputPath += ".exe"
	}

	var buildArgs []string
	env := os.Environ()

	// Set cross-compilation environment variables
	if target != "" {
		env = append(env, fmt.Sprintf("GOOS=%s", targetOS))
		env = append(env, fmt.Sprintf("GOARCH=%s", targetArch))
	}

	// Set CGO_ENABLED based on flags
	if withCgo {
		env = append(env, "CGO_ENABLED=1")
	} else if release {
		// Only disable CGO in release mode if --with-cgo is not specified
		env = append(env, "CGO_ENABLED=0")
	}

	if release {
		if target != "" {
			fmt.Printf("Building in release mode for %s/%s", targetOS, targetArch)
		} else {
			fmt.Print("Building in release mode")
		}
		if withCgo {
			fmt.Print(" with CGO enabled")
		}
		fmt.Println("...")
		buildArgs = []string{"build", "-ldflags=-s -w", "-trimpath", "-o", outputPath}
	} else {
		if target != "" {
			fmt.Printf("Building in debug mode for %s/%s", targetOS, targetArch)
		} else {
			fmt.Print("Building in debug mode")
		}
		if withCgo {
			fmt.Print(" with CGO enabled")
		}
		fmt.Println("...")
		buildArgs = []string{"build", "-o", outputPath}
	}

	// Determine source path based on project mode
	if projectMode == "project" {
		buildArgs = append(buildArgs, "./cmd/server")
	} else {
		buildArgs = append(buildArgs, ".")
	}

	// Ensure output directory exists
	fullOutputDir := filepath.Join(projectRoot, outputDir)
	if err := os.MkdirAll(fullOutputDir, 0755); err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	cmd := exec.Command("go", buildArgs...)
	cmd.Dir = projectRoot
	cmd.Env = env

	// Capture both stdout and stderr to ensure we display all output
	output, err := cmd.CombinedOutput()

	if err != nil {
		// Display all output from the build command
		if len(output) > 0 {
			fmt.Print(string(output))
		}
		fmt.Printf("\nBuild failed: %v\n", err)
		os.Exit(1)
	}

	// Display any output even on success (e.g., warnings)
	if len(output) > 0 {
		fmt.Print(string(output))
	}

	fmt.Printf("Build successful: %s\n", outputPath)
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
		sourcePath = "./cmd/server"
	} else {
		sourcePath = "."
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

	// Detect project mode: prioritize checking directory structure
	// Check for project mode first (cmd/server directory exists)
	cmdServerDir := filepath.Join(projectRoot, "cmd", "server")
	if stat, err := os.Stat(cmdServerDir); err == nil && stat.IsDir() {
		projectMode = "project"
	} else if _, err := os.Stat(filepath.Join(projectRoot, "main.go")); err == nil {
		// Simple mode: main.go in root
		projectMode = "simple"
	} else {
		return "", "", "", fmt.Errorf("cannot detect project mode: no main.go found and cmd/server directory doesn't exist")
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

func runCommand(dir string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir

	// Capture output to display in case of error
	output, err := cmd.CombinedOutput()

	// Always display output (for progress messages, warnings, etc.)
	if len(output) > 0 {
		fmt.Print(string(output))
	}

	return err
}

// ==================== ADD COMMAND ====================

func handleAdd(args []string) {
	if len(args) < 1 {
		fmt.Println("Error: Missing package name")
		fmt.Println("Usage: gocar add <package>...")
		fmt.Println("Example: gocar add github.com/gin-gonic/gin")
		os.Exit(1)
	}

	// Check for help
	if args[0] == "help" || args[0] == "--help" || args[0] == "-h" {
		fmt.Println("gocar add - Add dependencies to the project")
		fmt.Println()
		fmt.Println("USAGE:")
		fmt.Println("    gocar add <package>...")
		fmt.Println()
		fmt.Println("EXAMPLES:")
		fmt.Println("    gocar add github.com/gin-gonic/gin")
		fmt.Println("    gocar add github.com/gin-gonic/gin github.com/spf13/cobra")
		os.Exit(0)
	}

	// Check if we're in a Go module
	projectRoot, appName, _, err := detectProject()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Adding dependencies to '%s'...\n", appName)

	// Add each package
	for _, pkg := range args {
		fmt.Printf("  Adding %s...\n", pkg)
		getArgs := []string{"get", pkg}
		if err := runCommand(projectRoot, "go", getArgs...); err != nil {
			fmt.Printf("Error adding %s: %v\n", pkg, err)
			os.Exit(1)
		}
	}

	// Run go mod tidy to clean up
	fmt.Println("Tidying go.mod...")
	if err := runCommand(projectRoot, "go", "mod", "tidy"); err != nil {
		fmt.Printf("Warning: Failed to tidy go.mod: %v\n", err)
	}

	fmt.Println("Successfully added dependencies")
}

// ==================== UPDATE COMMAND ====================

func handleUpdate(args []string) {
	// Check for help
	if len(args) > 0 && (args[0] == "help" || args[0] == "--help" || args[0] == "-h") {
		fmt.Println("gocar update - Update dependencies")
		fmt.Println()
		fmt.Println("USAGE:")
		fmt.Println("    gocar update [package]...")
		fmt.Println()
		fmt.Println("EXAMPLES:")
		fmt.Println("    gocar update                           Update all dependencies")
		fmt.Println("    gocar update github.com/gin-gonic/gin  Update specific package")
		os.Exit(0)
	}

	// Check if we're in a Go module
	projectRoot, appName, _, err := detectProject()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if len(args) == 0 {
		// Update all dependencies
		fmt.Printf("Updating all dependencies for '%s'...\n", appName)
		if err := runCommand(projectRoot, "go", "get", "-u", "./..."); err != nil {
			fmt.Printf("Error updating dependencies: %v\n", err)
			os.Exit(1)
		}
	} else {
		// Update specific packages
		fmt.Printf("Updating specified dependencies for '%s'...\n", appName)
		for _, pkg := range args {
			fmt.Printf("  Updating %s...\n", pkg)
			if err := runCommand(projectRoot, "go", "get", "-u", pkg); err != nil {
				fmt.Printf("Error updating %s: %v\n", pkg, err)
				os.Exit(1)
			}
		}
	}

	// Run go mod tidy to clean up
	fmt.Println("Tidying go.mod...")
	if err := runCommand(projectRoot, "go", "mod", "tidy"); err != nil {
		fmt.Printf("Warning: Failed to tidy go.mod: %v\n", err)
	}

	fmt.Println("Successfully updated dependencies")
}

// ==================== TIDY COMMAND ====================

func handleTidy() {
	// Check if we're in a Go module
	projectRoot, appName, _, err := detectProject()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Tidying go.mod for '%s'...\n", appName)

	if err := runCommand(projectRoot, "go", "mod", "tidy"); err != nil {
		fmt.Printf("Error tidying go.mod: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Successfully tidied go.mod")
}
