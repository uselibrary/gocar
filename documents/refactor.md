
# Gocar 项目重构方案

## 新的项目结构

```
gocar/
├── cmd/
│   └── gocar/
│       └── main.go                 # 程序入口，最小化逻辑
├── internal/
│   ├── cli/
│   │   ├── cli.go                  # CLI 应用初始化和路由
│   │   ├── new.go                  # new 命令实现
│   │   ├── build.go                # build 命令实现
│   │   ├── run.go                  # run 命令实现
│   │   ├── clean.go                # clean 命令实现
│   │   ├── dep.go                  # add/update/tidy 命令实现
│   │   └── version.go              # version 命令实现
│   ├── project/
│   │   ├── project.go              # 项目检测和信息
│   │   ├── creator.go              # 项目创建逻辑
│   │   ├── validator.go            # 项目名称验证
│   │   └── template.go             # 项目模板内容
│   ├── build/
│   │   ├── builder.go              # 构建逻辑封装
│   │   ├── config.go               # 构建配置
│   │   └── target.go               # 目标平台处理
│   └── util/
│       ├── exec.go                 # 命令执行工具
│       ├── file.go                 # 文件操作工具
│       └── git.go                  # Git 操作工具
├── pkg/                            # (可选) 公共库
├── .gitignore
├── go.mod
├── go.sum
├── LICENSE
├── README.md
└── Makefile                        # 构建脚本
```

## 重构要点

### 1. 项目结构优化
- **cmd/gocar/main.go**: 只负责启动 CLI 应用
- **internal/cli**: 命令行接口层，每个命令一个文件
- **internal/project**: 项目相关的核心逻辑
- **internal/build**: 构建相关的核心逻辑
- **internal/util**: 通用工具函数

### 2. 关注点分离
- 命令处理（CLI 层）与业务逻辑（domain 层）分离
- 文件操作、进程执行等工具函数独立封装
- 模板内容与逻辑代码分离

### 3. 代码改进
- 使用接口提高可测试性
- 错误处理更加统一和明确
- 添加适当的注释和文档
- 使用常量管理魔法字符串

### 4. 建议的改进
- 添加单元测试（testing）
- 使用 cobra 库来管理 CLI（可选，当前实现也可以）
- 添加配置文件支持（.gocar.yaml）
- 添加日志系统
- 改进错误类型和错误处理

## 核心文件示例

### cmd/gocar/main.go
```go
package main

import (
    "os"
    "gocar/internal/cli"
)

func main() {
    app := cli.NewApp()
    if err := app.Run(os.Args); err != nil {
        os.Exit(1)
    }
}
```

### internal/cli/cli.go
```go
package cli

import (
    "fmt"
    "os"
)

const Version = "0.1.3"

type App struct {
    commands map[string]Command
}

type Command interface {
    Run(args []string) error
    Help() string
}

func NewApp() *App {
    app := &App{
        commands: make(map[string]Command),
    }
    
    // 注册命令
    app.commands["new"] = &NewCommand{}
    app.commands["build"] = &BuildCommand{}
    app.commands["run"] = &RunCommand{}
    app.commands["clean"] = &CleanCommand{}
    app.commands["add"] = &AddCommand{}
    app.commands["update"] = &UpdateCommand{}
    app.commands["tidy"] = &TidyCommand{}
    
    return app
}

func (a *App) Run(args []string) error {
    if len(args) < 2 {
        printHelp()
        return nil
    }
    
    cmdName := args[1]
    
    // 处理特殊命令
    switch cmdName {
    case "help", "-h", "--help":
        printHelp()
        return nil
    case "version", "-v", "--version":
        fmt.Printf("gocar %s\n", Version)
        return nil
    }
    
    // 执行命令
    cmd, ok := a.commands[cmdName]
    if !ok {
        fmt.Printf("Unknown command: %s\n", cmdName)
        printHelp()
        return fmt.Errorf("unknown command: %s", cmdName)
    }
    
    return cmd.Run(args[2:])
}
```

### internal/project/project.go
```go
package project

import (
    "fmt"
    "os"
    "path/filepath"
)

type Info struct {
    Root string
    Name string
    Mode string // "simple" or "project"
}

type Detector struct{}

func NewDetector() *Detector {
    return &Detector{}
}

func (d *Detector) Detect() (*Info, error) {
    // 查找项目根目录
    cwd, err := os.Getwd()
    if err != nil {
        return nil, fmt.Errorf("failed to get current directory: %w", err)
    }
    
    root := cwd
    for {
        if _, err := os.Stat(filepath.Join(root, "go.mod")); err == nil {
            break
        }
        parent := filepath.Dir(root)
        if parent == root {
            return nil, fmt.Errorf("not in a Go module (go.mod not found)")
        }
        root = parent
    }
    
    // 检测项目模式
    mode := d.detectMode(root)
    if mode == "" {
        return nil, fmt.Errorf("cannot detect project mode")
    }
    
    return &Info{
        Root: root,
        Name: filepath.Base(root),
        Mode: mode,
    }, nil
}

func (d *Detector) detectMode(root string) string {
    cmdServerDir := filepath.Join(root, "cmd", "server")
    if stat, err := os.Stat(cmdServerDir); err == nil && stat.IsDir() {
        return "project"
    }
    
    if _, err := os.Stat(filepath.Join(root, "main.go")); err == nil {
        return "simple"
    }
    
    return ""
}
```

### internal/build/builder.go
```go
package build

import (
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "runtime"
)

type Config struct {
    Release    bool
    TargetOS   string
    TargetArch string
    WithCGO    bool
}

type Builder struct {
    config     *Config
    projectRoot string
    appName     string
    projectMode string
}

func NewBuilder(projectRoot, appName, projectMode string, config *Config) *Builder {
    if config.TargetOS == "" {
        config.TargetOS = runtime.GOOS
    }
    if config.TargetArch == "" {
        config.TargetArch = runtime.GOARCH
    }
    
    return &Builder{
        config:      config,
        projectRoot: projectRoot,
        appName:     appName,
        projectMode: projectMode,
    }
}

func (b *Builder) Build() error {
    outputPath := b.getOutputPath()
    
    // 确保输出目录存在
    if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
        return fmt.Errorf("failed to create output directory: %w", err)
    }
    
    // 构建命令
    cmd := b.buildCommand(outputPath)
    
    // 执行构建
    output, err := cmd.CombinedOutput()
    if len(output) > 0 {
        fmt.Print(string(output))
    }
    
    if err != nil {
        return fmt.Errorf("build failed: %w", err)
    }
    
    fmt.Printf("Build successful: %s\n", outputPath)
    return nil
}

func (b *Builder) getOutputPath() string {
    buildMode := "debug"
    if b.config.Release {
        buildMode = "release"
    }
    
    targetDir := fmt.Sprintf("%s-%s", b.config.TargetOS, b.config.TargetArch)
    outputDir := filepath.Join(b.projectRoot, "bin", buildMode, targetDir)
    outputPath := filepath.Join(outputDir, b.appName)
    
    if b.config.TargetOS == "windows" {
        outputPath += ".exe"
    }
    
    return outputPath
}

func (b *Builder) buildCommand(outputPath string) *exec.Cmd {
    args := []string{"build"}
    
    if b.config.Release {
        args = append(args, "-ldflags=-s -w", "-trimpath")
    }
    
    args = append(args, "-o", outputPath)
    
    // 添加源码路径
    if b.projectMode == "project" {
        args = append(args, "./cmd/<appName>")
    } else {
        args = append(args, ".")
    }
    
    cmd := exec.Command("go", args...)
    cmd.Dir = b.projectRoot
    cmd.Env = b.buildEnv()
    
    return cmd
}

func (b *Builder) buildEnv() []string {
    env := os.Environ()
    
    env = append(env, fmt.Sprintf("GOOS=%s", b.config.TargetOS))
    env = append(env, fmt.Sprintf("GOARCH=%s", b.config.TargetArch))
    
    if b.config.WithCGO {
        env = append(env, "CGO_ENABLED=1")
    } else if b.config.Release {
        env = append(env, "CGO_ENABLED=0")
    }
    
    return env
}
```

## 测试结构

```
gocar/
├── internal/
│   ├── cli/
│   │   └── cli_test.go
│   ├── project/
│   │   ├── project_test.go
│   │   ├── creator_test.go
│   │   └── validator_test.go
│   ├── build/
│   │   └── builder_test.go
│   └── util/
│       ├── exec_test.go
│       └── file_test.go
└── testdata/                       # 测试数据
    └── projects/
        ├── simple/
        └── project/
```

## 构建和安装

### Makefile
```makefile
.PHONY: build install test clean

build:
	go build -o bin/gocar ./cmd/gocar

install:
	go install ./cmd/gocar

test:
	go test -v ./...

test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

clean:
	rm -rf bin/
	rm -f coverage.out

lint:
	golangci-lint run

.DEFAULT_GOAL := build
```

## 优势总结

1. **更好的代码组织**: 按功能模块清晰分层
2. **易于测试**: 逻辑与 CLI 分离，可以单独测试
3. **易于维护**: 单一职责原则，每个文件职责明确
4. **易于扩展**: 添加新命令或新功能更容易
5. **符合 Go 标准**: 遵循 Go 社区的最佳实践
6. **更好的错误处理**: 统一的错误处理和返回机制