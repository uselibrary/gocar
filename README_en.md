# gocar, a cargo for Go

> A Rust Cargo-like Go project scaffold and CLI tool that provides a simple project
initialization and build workflow.

[![License: MIT](https://img.shields.io/badge/License-MIT.svg)](https://opensource.org/licenses/MIT)
[![Go](https://img.shields.io/badge/go-1.25+-blue.svg)](https://golang.org)
[![Platform](https://img.shields.io/badge/platform-Linux%20%7C%20macOS%20%7C%20Windows-lightgrey.svg)](https://github.com/uselibrary/gocar)

**[简体中文](README.md)** | **[English](README_en.md)**

## Installation

> git is required for some commands, please ensure it is installed.

### Binary Installation (Recommended)
Download the precompiled binary for your operating system from the release page,
extract it, and move it to your `$PATH` directory:

```bash
/usr/local/bin/ # Unix-like systems, e.g., Linux or macOS
C:\Program Files\ # Windows
```
For Unix-like systems, ensure the binary has executable permissions (requires root):
```bash
chown root:root /usr/local/bin/gocar
chmod +x /usr/local/bin/gocar
```

### Or Build from Source:

```bash
git clone https://github.com/uselibrary/gocar.git
cd gocar
CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o gocar main.go
sudo mv gocar /usr/local/bin/
sudo chown root:root /usr/local/bin/gocar
sudo chmod +x /usr/local/bin/gocar
```

## Quick Start

```bash
# Create a new project (simple mode)
gocar new myapp

# Enter the project directory
cd myapp

# Build the project
gocar build

# Run the project
gocar run

# Clean build artifacts
gocar clean
```

## Commands

### `gocar new <name> [--mode simple|project]`

Create a new Go project.

Arguments:
- `<name>` - project name; used as directory name and output binary name
- `--mode` - project mode, either `simple` (default) or `project`

Project name rules:
- Must start with a letter
- May contain letters, digits, underscores `_`, or hyphens `-`
- Must not use reserved names: `test`, `main`, `init`, `internal`, `vendor`

Examples:

```bash
# Create a simple mode project (default)
gocar new myapp

# Create a project mode project
gocar new myserver --mode project
```

### `gocar build [--release] [--target <os>/<arch>] [--help]`

Build the current project.

**Options:**
- `--release` - build in release mode (optimized binary size)
- `--target <os>/<arch>` - cross-compile for target platform
- `--help` - show build command help

**Build behavior:**

| Mode | Equivalent command |
|------|--------------------|
| Debug (default) | `go build -o bin/<appName> ./main.go` |
| Release | `CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o bin/<appName> ./main.go` |
| Cross-compile | `GOOS=<os> GOARCH=<arch> go build ...` |

> In project mode the entry point is `./cmd/server/main.go`

**Common target platforms:**
- `linux/amd64` - Linux 64-bit
- `linux/arm64` - Linux ARM 64-bit
- `darwin/amd64` - macOS Intel
- `darwin/arm64` - macOS Apple Silicon
- `windows/amd64` - Windows 64-bit
- `windows/386` - Windows 32-bit

**Examples:**

```bash
# Debug build
gocar build

# Release build (smaller binary)
gocar build --release

# Cross-compile for Linux AMD64
gocar build --target linux/amd64

# Release mode cross-compile for Windows
gocar build --release --target windows/amd64

# Show help information
gocar build --help
```

### `gocar run [args...]`

Run the current project directly using `go run`.

Examples:

```bash
# Run the project
gocar run

# Pass arguments to the application
gocar run --port 8080
```

### `gocar clean`

Clean the `bin/` directory build artifacts.

Example:

```bash
gocar clean
# Cleaned build artifacts for 'myapp'
```

### `gocar add <package>...`

Add dependencies to the project.

**Arguments:**
- `<package>` - package name(s) to add, supports multiple packages

**Examples:**
```bash
# Add a single dependency
gocar add github.com/gin-gonic/gin

# Add multiple dependencies
gocar add github.com/gin-gonic/gin github.com/spf13/cobra
```

### `gocar update [package]...`

Update project dependencies.

**Arguments:**
- `[package]` - optional, specific package(s) to update. If omitted, updates all dependencies

**Examples:**
```bash
# Update all dependencies
gocar update

# Update specific dependencies
gocar update github.com/gin-gonic/gin
```

### `gocar tidy`

Tidy up `go.mod` and `go.sum` files, removing unused dependencies.

**Example:**
```bash
gocar tidy
# Successfully tidied go.mod
```

### `gocar help`

Show help information.

### `gocar version`

Show version information.

## Project Modes

### Simple Mode

Suitable for small projects, scripts, or CLI tools.

```
myapp/
├── go.mod
├── main.go
├── README.md
├── bin/
├── .gitignore
└── .git/
```

main.go template:

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Hello, gocar! A golang project scaffolding tool for <appName>")
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
}
```

### Project Mode

Suitable for larger projects, web services, or microservices and follows a
standard Go project layout.

```
myapp/
├── cmd/
│   └── server/
│       └── main.go      # application entry
├── internal/            # private code (not importable by other modules)
├── pkg/                 # public libraries for external import
├── test/                # integration tests
├── bin/                 # build output
├── go.mod
├── README.md
├── .gitignore
└── .git/
```

`cmd/server/main.go` template:

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Hello, gocar! A golang project scaffolding tool for <appName>")
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
}
```

Directory notes:
- `cmd/` - application entry points
- `internal/` - private code (enforced by Go)
- `pkg/` - public libraries intended for external use
- `test/` - integration and end-to-end tests

## Features

- ✅ Automatic Git initialization (`git init -b main`) and `.gitignore` generation
- ✅ Smart project mode detection (simple vs project)
- ✅ Project name validation following Go conventions
- ✅ Release optimizations using `CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath`
- ✅ Clean command to remove build artifacts
- ✅ Cross-platform support

## .gitignore Template

Automatically generated `.gitignore` includes:

```gitignore
# Binaries
<appName>
bin/
*.exe
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
```

## License

MIT License


