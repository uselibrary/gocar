# gocar, a cargo for Go

> 一个"类 Rust Cargo"的 Go 项目脚手架与命令行工具，提供简洁的项目初始化和构建体验。

[![License: MIT](https://img.shields.io/badge/License-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)
[![Go](https://img.shields.io/badge/go-1.25+-yellow.svg)](https://golang.org)
[![Platform](https://img.shields.io/badge/platform-Linux%20|%20macOS%20|%20Windows-blue.svg)](https://github.com/uselibrary/gocar)

**[简体中文](README.md)** | **[English](README_en.md)**

## 安装

> `git` 是某些命令的前置依赖，请确保已安装。

### 二进制安装（推荐）
从 [release页面](https://github.com/uselibrary/gocar/releases) 下载适合你操作系统的预编译二进制文件，解压后将其移动到`$PATH`目录中：
```bash
/usr/local/bin/ # Unix-like 系统, 例如 Linux 或 macOS
C:\Program Files\ # Windows 系统，可能需要设置环境变量
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



## 命令详解

### 新建项目

**`gocar new <appName> [--mode simple|project]`**

创建新的 Go 项目:
- `gocar new <appName>` 创建简洁模式项目（默认）
- `gocar new <appName> --mode project` 创建项目模式项目

简洁模型的目录结构：
```
<appName>/
├── go.mod
├── main.go
├── README.md
├── bin/
├── .gitignore
└── .git/
``` 

项目模式的目录结构：
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

> 注意：创建的项目默认不包含 `.gocar.toml`，可通过 `gocar init` 手动生成。

> 简洁模型式适用于小型项目、脚本、CLI 工具等；项目模式适用于大型项目、Web 服务、微服务等，遵循 Go 标准项目布局。

> `<appName>`为项目名称，同时作为目录名和输出的可执行文件名；`--mode`为项目模式，可选 `simple`（默认）或 `project`

> 项目名规则：
> - 必须以字母开头
> - 只能包含字母、数字、下划线 `_` 或连字符 `-`
> - 不能使用保留名称：`test`、`main`、`init`、`internal`、`vendor`


### 编译构建


**`gocar build [--release] [--target <os>/<arch>] [--with-cgo] [--help]`**

构建可执行文件：
- `gocar build` ` 构建 Debug 版本（默认）
- `gocar build --release` 构建 Release 版本（启用CGO_ENABLED=0，ldflags="-s -w" 和 trimpath）
- `gocar build --target <os>/<arch>` 交叉编译到指定平台
- `gocar build --release --target <os>/<arch>` 以 Release 模式交叉编译到指定平台
- `gocar build --with-cgo` 强制启用 CGO（设置 CGO_ENABLED=1）
- `gocar build --help` 显示帮助信息

构建行为：

| 模式 | 命令等价 |
|------|----------|
| debug（默认） | `go build -o bin/<os>/<arch>/<appName> ./main.go` |
| -- release | `CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o bin/<os>/<arch>/<appName> ./main.go` |
| -- target| `GOOS=<os> GOARCH=<arch> go build -o bin/<os>/<arch>/<appName> ./main.go` |
| -- release -- target | `CGO_ENABLED=0 GOOS=<os> GOARCH=<arch> go build -ldflags="-s -w" -trimpath -o bin/<os>/<arch>/<appName> ./main.go` |
| -- with-cgo | `CGO_ENABLED=1 go build -o bin/<os>/<arch>/<appName> ./main.go` |
| -- release -- with-cgo | `CGO_ENABLED=1 go build -ldflags="-s -w" -trimpath -o bin/<os>/<arch>/<appName> ./main.go` |

示例：
```bash
# Debug 构建（默认）
gocar build

# Release 构建（启用CGO_ENABLED=0，ldflags="-s -w" 和 trimpath）
gocar build --release

# 交叉编译到指定系统和架构，例如 Linux AMD64
gocar build --target linux/amd64

# Release 模式交叉编译到 Windows AMD64（启用CGO_ENABLED=0，ldflags="-s -w" 和 trimpath）
gocar build --release --target windows/amd64

# 强制启用 CGO 构建
gocar build --with-cgo

# Release 模式下启用 CGO 构建
gocar build --release --with-cgo

# 显示帮助信息
gocar build --help
```

### 常用命令

**`gocar run [args...]`**

直接运行当前项目（使用 `go run`）。

示例：
```bash
# 运行项目
gocar run

# 传递参数给应用
gocar run --port 8080
```

**`gocar clean`**

清理 `bin/` 目录中的构建产物。

*示例：
```bash
gocar clean
# Cleaned build artifacts for '<appName>'
```

**`gocar help`**

显示帮助信息。

**`gocar version`**

显示版本信息。

### 包操作

**`gocar add <package>...`**

添加、更新依赖：
- `gocar add <package>` 添加指定依赖
- `gocar update <package>` 更新指定依赖
- `gocar update` 更新所有依赖
- `gocar tidy` 整理 `go.mod` 和 `go.sum`
- `gocar add` 等同于 `go get <package>...` 并更新 `go.mod` 和 `go.sum`

依赖行为：
| 命令 | 等价 |
|------|----------|
| gocar add <package>... | go get <package>... |
| gocar update [package]... | go get -u [package]... |
| gocar update  | go get -u ./... & go mod tidy |
| gocar tidy | go mod tidy |


示例：
```bash
# 添加指定依赖
gocar add github.com/gin-gonic/gin

# 更新所有依赖
gocar update

# 更新指定依赖
gocar update github.com/gin-gonic/gin

# 整理依赖
gocar tidy
# Successfully tidied go.mod
```

### 配置文件

**`gocar init`**

在当前项目中生成 `.gocar.toml` 配置文件。配置文件中的设置优先级高于 gocar 的自动检测。

示例：
```bash
# 在已有项目中生成配置文件
gocar init
# Created .gocar.toml in /path/to/project
```

**配置文件结构：**

```toml
# gocar 项目配置文件

# 项目配置
[project]
mode = "project"    # 项目模式: "simple" 或 "project"
name = "myapp"      # 项目名称，留空则使用目录名

# 构建配置
[build]
entry = "cmd/server"                  # 构建入口路径（可修改为 cmd/myapp 等）
output = "bin"                        # 输出目录
ldflags = "-X main.version=1.0.0"     # 额外的 ldflags
# tags = ["jsoniter", "sonic"]        # 构建标签
# extra_env = ["GOPROXY=https://goproxy.cn"]  # 额外环境变量

# 运行配置
[run]
entry = ""                            # 运行入口，留空则使用 build.entry
# args = ["-config", "config.yaml"]   # 默认运行参数

# Debug 构建配置 (gocar build)
[profile.debug]
# ldflags = ""              # Debug 默认无 ldflags
# gcflags = "all=-N -l"     # 禁用优化，方便调试
# trimpath = false          # 保留路径信息
# cgo_enabled = true        # 跟随系统默认
# race = false              # 竞态检测

# Release 构建配置 (gocar build --release)
[profile.release]
ldflags = "-s -w"           # 裁剪符号表和调试信息
# gcflags = ""              # 编译器参数
trimpath = true             # 移除编译路径信息
cgo_enabled = false         # 禁用 CGO 以生成静态二进制
# race = false              # 竞态检测

# 自定义命令
[commands]
vet = "go vet ./..."
fmt = "go fmt ./..."
test = "go test -v ./..."
# lint = "golangci-lint run"
```

**配置项说明：**

| 配置项 | 说明 |
|--------|------|
| `[project].mode` | 指定项目模式 (`simple` 或 `project`)，留空则自动检测 |
| `[project].name` | 自定义项目名称，留空则使用目录名 |
| `[build].entry` | **自定义构建入口路径**，如 `cmd/myapp` 替代默认的 `cmd/server` |
| `[build].ldflags` | 额外的 ldflags，会追加到 profile 的 ldflags 之后 |
| `[build].tags` | 构建标签列表 |
| `[build].extra_env` | 额外的环境变量 |
| `[run].entry` | 运行入口路径，留空则使用 `build.entry` |
| `[run].args` | 默认运行参数 |
| `[profile.debug]` | Debug 构建模式的参数配置 |
| `[profile.release]` | Release 构建模式的参数配置 |
| `[commands]` | 自定义命令映射 |

**Profile 配置项：**

| 配置项 | 说明 | Debug 默认 | Release 默认 |
|--------|------|-------------|---------------|
| `ldflags` | 链接器参数 | `""` | `"-s -w"` |
| `gcflags` | 编译器参数 | `""` | `""` |
| `trimpath` | 移除路径信息 | `false` | `true` |
| `cgo_enabled` | 启用 CGO | `nil` (系统默认) | `false` |
| `race` | 竞态检测 | `false` | `false` |

### 自定义命令

在 `.gocar.toml` 的 `[commands]` 部分定义命令后，可以直接执行：

```bash
# 代码检查
gocar vet

# 代码格式化
gocar fmt

# 运行测试
gocar test

# 传递额外参数
gocar test -run TestXxx
```

命令输出会实时显示到终端。您可以自定义任意命令，例如：

```toml
[commands]
lint = "golangci-lint run"
doc = "godoc -http=:6060"
proto = "protoc --go_out=. --go-grpc_out=. ./proto/*.proto"
dev = "air"  # 热重载
```

---

新建项目的 `main.go` 模板内容如下：
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

---

## License

MIT License
