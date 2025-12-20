# gocar, a cargo for Go

> A “Rust Cargo–like” scaffolding and CLI tool for Go projects, providing a clean experience for project initialization and building.

[![License: MIT](https://img.shields.io/badge/License-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)
[![Go](https://img.shields.io/badge/go-1.25+-yellow.svg)](https://golang.org)
[![Platform](https://img.shields.io/badge/platform-Linux%20|%20macOS%20|%20Windows-blue.svg)](https://github.com/uselibrary/gocar)

**[简体中文](README.md)** | **[English](README_en.md)**

## Installation

> `git` is a prerequisite for some commands. Please make sure it is installed.

### Install via binary (recommended)

Download the prebuilt binary for your OS from the [Releases page](https://github.com/uselibrary/gocar/releases). Extract it and move it into a directory in your `$PATH`:

```
/usr/local/bin/ # Unix-like systems, e.g. Linux or macOS
C:\Program Files\ # Windows, may need to set environment variables
```

On Unix-like systems, ensure the binary has executable permissions (root required):

```
chown root:root /usr/local/bin/gocar
chmod +x /usr/local/bin/gocar
```

### Or build from source

```
git clone https://github.com/uselibrary/gocar.git
cd gocar
CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o gocar main.go
sudo mv gocar /usr/local/bin/
sudo chown root:root /usr/local/bin/gocar
sudo chmod +x /usr/local/bin/gocar
```

## Quick Start

```
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

## Command Reference

### Create a new project

**`gocar new <appName> [--mode simple|project]`**

Create a new Go project:

- `gocar new <appName>` creates a simple-mode project (default)
- `gocar new <appName> --mode project` creates a project-mode project

Directory structure for **simple mode**:

```
<appName>/
├── go.mod
├── main.go
├── README.md
├── bin/
├── .gitignore
└── .git/
```

Directory structure for **project mode**:

```
<appName>/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
├── pkg/
├── test/
├── bin/
├── go.mod
├── README.md
├── .gitignore
└── .git/
```

> Note: Created projects do not include `.gocar.toml` by default. Use `gocar init` to generate it manually.

> Simple mode is suitable for small projects, scripts, CLI tools, etc. Project mode is suitable for larger projects, web services, microservices, etc., following the standard Go project layout.

> `<appName>` is the project name, used as the directory name and the output executable name. `--mode` selects the project mode: `simple` (default) or `project`.

> Project name rules:
>
> - Must start with a letter
> - May contain only letters, digits, underscore `_`, or hyphen `-`
> - Reserved names are not allowed: `test`, `main`, `init`, `internal`, `vendor`

### Build / Compile

**`gocar build [--release] [--target <os>/<arch>] [--with-cgo] [--help]`**

Build the executable:

- `gocar build` builds a Debug binary (default)
- `gocar build --release` builds a Release binary (enables `CGO_ENABLED=0`, `ldflags="-s -w"`, and `-trimpath`)
- `gocar build --target <os>/<arch>` cross-compiles for the specified platform
- `gocar build --release --target <os>/<arch>` cross-compiles in Release mode for the specified platform
- `gocar build --with-cgo` forces CGO to be enabled (sets `CGO_ENABLED=1`)
- `gocar build --help` shows help information

Build behavior:

| Mode                    | Equivalent command                                           |
| ----------------------- | ------------------------------------------------------------ |
| debug (default)         | `go build -o bin/<os>/<arch>/<appName> ./main.go`            |
| --release               | `CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o bin/<os>/<arch>/<appName> ./main.go` |
| --target                | `GOOS=<os> GOARCH=<arch> go build -o bin/<os>/<arch>/<appName> ./main.go` |
| --release --target      | `CGO_ENABLED=0 GOOS=<os> GOARCH=<arch> go build -ldflags="-s -w" -trimpath -o bin/<os>/<arch>/<appName> ./main.go` |
| --with-cgo              | `CGO_ENABLED=1 go build -o bin/<os>/<arch>/<appName> ./main.go` |
| --release --with-cgo    | `CGO_ENABLED=1 go build -ldflags="-s -w" -trimpath -o bin/<os>/<arch>/<appName> ./main.go` |

Examples:

```
# Debug build (default)
gocar build

# Release build (enables CGO_ENABLED=0, ldflags="-s -w" and trimpath)
gocar build --release

# Cross-compile for a target OS/arch, e.g. Linux AMD64
gocar build --target linux/amd64

# Release cross-compile for Windows AMD64 (enables CGO_ENABLED=0, ldflags="-s -w" and trimpath)
gocar build --release --target windows/amd64

# Force enable CGO
gocar build --with-cgo

# Release build with CGO enabled
gocar build --release --with-cgo

# Show help
gocar build --help
```

### Common commands

**`gocar run [args...]`**

Run the current project directly (uses `go run`).

Examples:

```
# Run the project
gocar run

# Pass arguments to the app
gocar run --port 8080
```

**`gocar clean`**

Remove build artifacts in the `bin/` directory.

Example:

```
gocar clean
# Cleaned build artifacts for '<appName>'
```

**`gocar help`**

Show help information.

**`gocar version`**

Show version information.

### Package management

**`gocar add <package>...`**

Add or update dependencies:

- `gocar add <package>` add the specified dependency
- `gocar update <package>` update the specified dependency
- `gocar update` update all dependencies
- `gocar tidy` tidy `go.mod` and `go.sum`
- `gocar add` is equivalent to `go get <package>...` and updates `go.mod` and `go.sum`

Dependency behavior:

| Command                     | Equivalent                      |
| --------------------------- | ------------------------------- |
| `gocar add <package>...`    | `go get <package>...`           |
| `gocar update [package]...` | `go get -u [package]...`        |
| `gocar update`              | `go get -u ./... & go mod tidy` |
| `gocar tidy`                | `go mod tidy`                   |

Examples:

```
# Add a dependency
gocar add github.com/gin-gonic/gin

# Update all dependencies
gocar update

# Update a specific dependency
gocar update github.com/gin-gonic/gin

# Tidy dependencies
gocar tidy
# Successfully tidied go.mod
```

### Configuration File

**`gocar init`**

Generate a `.gocar.toml` configuration file in the current project. Settings in the config file take priority over gocar's auto-detection.

Example:
```bash
# Generate config file in existing project
gocar init
# Created .gocar.toml in /path/to/project
```

**Configuration file structure:**

```toml
# gocar project configuration file

# Project configuration
[project]
mode = "project"    # Project mode: "simple" or "project"
name = "myapp"      # Project name, uses directory name if empty

# Build configuration
[build]
entry = "cmd/server"                  # Build entry path (can be changed to cmd/myapp, etc.)
output = "bin"                        # Output directory
ldflags = "-X main.version=1.0.0"     # Additional ldflags
# tags = ["jsoniter", "sonic"]        # Build tags
# extra_env = ["GOPROXY=https://goproxy.cn"]  # Additional environment variables

# Run configuration
[run]
entry = ""                            # Run entry, uses build.entry if empty
# args = ["-config", "config.yaml"]   # Default run arguments

# Debug build configuration (gocar build)
[profile.debug]
# ldflags = ""              # Debug has no ldflags by default
# gcflags = "all=-N -l"     # Disable optimization for debugging
# trimpath = false          # Keep path information
# cgo_enabled = true        # Follow system default
# race = false              # Race detection

# Release build configuration (gocar build --release)
[profile.release]
ldflags = "-s -w"           # Strip symbol table and debug info
# gcflags = ""              # Compiler flags
trimpath = true             # Remove build path information
cgo_enabled = false         # Disable CGO for static binary
# race = false              # Race detection

# Custom commands
[commands]
vet = "go vet ./..."
fmt = "go fmt ./..."
test = "go test -v ./..."
# lint = "golangci-lint run"
```

**Configuration options:**

| Option | Description |
|--------|-------------|
| `[project].mode` | Specify project mode (`simple` or `project`), auto-detected if empty |
| `[project].name` | Custom project name, uses directory name if empty |
| `[build].entry` | **Custom build entry path**, e.g., `cmd/myapp` instead of default `cmd/server` |
| `[build].ldflags` | Additional ldflags, appended to profile ldflags |
| `[build].tags` | Build tags list |
| `[build].extra_env` | Additional environment variables |
| `[run].entry` | Run entry path, uses `build.entry` if empty |
| `[run].args` | Default run arguments |
| `[profile.debug]` | Debug build mode parameters |
| `[profile.release]` | Release build mode parameters |
| `[commands]` | Custom command mappings |

**Profile options:**

| Option | Description | Debug Default | Release Default |
|--------|-------------|---------------|------------------|
| `ldflags` | Linker flags | `""` | `"-s -w"` |
| `gcflags` | Compiler flags | `""` | `""` |
| `trimpath` | Remove path info | `false` | `true` |
| `cgo_enabled` | Enable CGO | `nil` (system) | `false` |
| `race` | Race detection | `false` | `false` |

### Custom Commands

After defining commands in the `[commands]` section of `.gocar.toml`, you can execute them directly:

```bash
# Code checking
gocar vet

# Code formatting
gocar fmt

# Run tests
gocar test

# Pass additional arguments
gocar test -run TestXxx
```

Command output is displayed in real-time to the terminal. You can define any custom commands, for example:

```toml
[commands]
lint = "golangci-lint run"
doc = "godoc -http=:6060"
proto = "protoc --go_out=. --go-grpc_out=. ./proto/*.proto"
dev = "air"  # Hot reload
```

------

The `main.go` template content for a new project is:

```
package main

import (
    "fmt"
    "time"
)

func main() {
    fmt.Println("Hello, gocar! A golang project scaffolding tool for <appName>.")
    fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
}
```

------

## License

MIT License