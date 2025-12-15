# gocar, a cargo for Go

> 一个"类 Rust Cargo"的 Go 项目脚手架与命令行工具，提供简洁的项目初始化和构建体验。

[![License: MIT](https://img.shields.io/badge/License-MIT.svg)](https://opensource.org/licenses/MIT)
[![Go](https://img.shields.io/badge/go-1.25+-blue.svg)](https://golang.org)
[![Platform](https://img.shields.io/badge/platform-Linux%20%7C%20macOS%20%7C%20Windows-lightgrey.svg)](https://github.com/uselibrary/gocar)

**[简体中文](README.md)** | **[English](README_en.md)**

## 安装

> git 是某些命令的前置依赖，请确保已安装。

### 二进制安装（推荐）
从release页面下载适合你操作系统的预编译二进制文件，解压后将其移动到你的`$PATH`目录中：
```bash
/usr/local/bin/ # Unix-like systems, 例如 Linux 或 macOS
C:\Program Files\ # Windows
```
对于Unix-like系统，确保二进制文件具有可执行权限（需要root权限）：
```bash
chown root:root /usr/local/bin/gocar
chmod +x /usr/local/bin/gocar
```

### 或从源码构建：

```bash
git clone https://github.com/uselibrary/gocar.git
cd gocar
CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o gocar main.go
sudo mv gocar /usr/local/bin/
sudo chown root:root /usr/local/bin/gocar
sudo chmod +x /usr/local/bin/gocar
```

## 快速开始

```bash
# 创建新项目（简洁模式）
gocar new myapp

# 进入项目目录
cd myapp

# 构建项目
gocar build

# 运行项目
gocar run

# 清理构建产物
gocar clean
```

## 命令

### `gocar new <name> [--mode simple|project]`

创建新的 Go 项目。

**参数：**
- `<name>` - 项目名称，同时作为目录名和输出的可执行文件名
- `--mode` - 项目模式，可选 `simple`（默认）或 `project`

**项目名规则：**
- 必须以字母开头
- 只能包含字母、数字、下划线 `_` 或连字符 `-`
- 不能使用保留名称：`test`、`main`、`init`、`internal`、`vendor`

**示例：**
```bash
# 创建简洁模式项目（默认）
gocar new myapp

# 创建项目模式项目
gocar new myserver --mode project
```

### `gocar build [--release] [--target <os>/<arch>] [--help]`

构建当前项目。

**参数：**
- `--release` - 使用 release 模式构建（优化体积）
- `--target <os>/<arch>` - 交叉编译到指定平台
- `--help` - 显示 build 命令的帮助信息

**构建行为：**

| 模式 | 命令等价 |
|------|----------|
| Debug（默认） | `go build -o bin/<appName> ./main.go` |
| Release | `CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o bin/<appName> ./main.go` |
| 交叉编译 | `GOOS=<os> GOARCH=<arch> go build ...` |

> 项目模式下入口为 `./cmd/server/main.go`

**常用目标平台：**
- `linux/amd64` - Linux 64位
- `linux/arm64` - Linux ARM 64位
- `darwin/amd64` - macOS Intel
- `darwin/arm64` - macOS Apple Silicon
- `windows/amd64` - Windows 64位
- `windows/386` - Windows 32位

**示例：**
```bash
# Debug 构建
gocar build

# Release 构建（更小的二进制文件）
gocar build --release

# 交叉编译到 Linux AMD64
gocar build --target linux/amd64

# Release 模式交叉编译到 Windows
gocar build --release --target windows/amd64

# 显示帮助信息
gocar build --help
```

### `gocar run [args...]`

直接运行当前项目（使用 `go run`）。

**示例：**
```bash
# 运行项目
gocar run

# 传递参数给应用
gocar run --port 8080
```

### `gocar clean`

清理 `bin/` 目录中的构建产物。

**示例：**
```bash
gocar clean
# Cleaned build artifacts for 'myapp'
```

### `gocar add <package>...`

添加依赖到项目中。

**参数：**
- `<package>` - 要添加的包名，支持一次添加多个包

**示例：**
```bash
# 添加单个依赖
gocar add github.com/gin-gonic/gin

# 添加多个依赖
gocar add github.com/gin-gonic/gin github.com/spf13/cobra
```

### `gocar update [package]...`

更新项目依赖。

**参数：**
- `[package]` - 可选，指定要更新的包名。不指定则更新所有依赖

**示例：**
```bash
# 更新所有依赖
gocar update

# 更新指定依赖
gocar update github.com/gin-gonic/gin
```

### `gocar tidy`

整理 `go.mod` 和 `go.sum` 文件，移除未使用的依赖。

**示例：**
```bash
gocar tidy
# Successfully tidied go.mod
```

### `gocar help`

显示帮助信息。

### `gocar version`

显示版本信息。

---

## 项目模式

### 简洁模式（Simple）

适用于小型项目、脚本、CLI 工具等。

```
myapp/
├── go.mod
├── main.go
├── README.md
├── bin/
├── .gitignore
└── .git/
```

**main.go 模板：**
```go
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

### 项目模式（Project）

适用于大型项目、Web 服务、微服务等，遵循 Go 标准项目布局。

```
myapp/
├── cmd/
│   └── server/
│       └── main.go      # 应用入口
├── internal/            # 私有代码（不可被外部导入）
├── pkg/                 # 公共库代码
├── test/                # 集成测试
├── bin/                 # 构建输出
├── go.mod
├── README.md
├── .gitignore
└── .git/
```

**cmd/server/main.go 模板：**
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

**目录说明：**
- `cmd/` - 应用程序入口
- `internal/` - 私有代码，仅本模块内部使用（Go 编译器强制）
- `pkg/` - 可被外部项目导入的公共库
- `test/` - 集成测试、端到端测试

---

## 特性

- ✅ **自动 Git 初始化** - 创建项目时自动执行 `git init -b main` 并生成 `.gitignore`
- ✅ **智能项目检测** - 自动识别 simple/project 模式
- ✅ **项目名验证** - 确保项目名符合 Go 规范
- ✅ **Release 优化** - 使用 `CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath` 减小二进制体积
- ✅ **清理命令** - 一键清理构建产物
- ✅ **跨平台支持**

---

## .gitignore 模板

自动生成的 `.gitignore` 包含：

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

---

## License

MIT License
