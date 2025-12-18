package config

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

// ConfigFileName 配置文件名
const ConfigFileName = ".gocar.toml"

// GocarConfig gocar 配置结构
type GocarConfig struct {
	Project  ProjectConfig     `toml:"project"`
	Build    BuildConfig       `toml:"build"`
	Run      RunConfig         `toml:"run"`
	Commands map[string]string `toml:"commands"`
}

// ProjectConfig 项目配置
type ProjectConfig struct {
	Mode string `toml:"mode"` // simple | project
	Name string `toml:"name"` // 项目名称，为空时使用目录名
}

// BuildConfig 构建配置
type BuildConfig struct {
	Entry    string   `toml:"entry"`     // 构建入口路径
	Output   string   `toml:"output"`    // 输出目录
	Ldflags  string   `toml:"ldflags"`   // 额外的 ldflags
	Tags     []string `toml:"tags"`      // 构建标签
	ExtraEnv []string `toml:"extra_env"` // 额外的环境变量
}

// RunConfig 运行配置
type RunConfig struct {
	Entry string   `toml:"entry"` // 运行入口路径
	Args  []string `toml:"args"`  // 默认运行参数
}

// DefaultConfig 返回默认配置
func DefaultConfig() *GocarConfig {
	return &GocarConfig{
		Project: ProjectConfig{
			Mode: "",
			Name: "",
		},
		Build: BuildConfig{
			Entry:    "",
			Output:   "bin",
			Ldflags:  "",
			Tags:     []string{},
			ExtraEnv: []string{},
		},
		Run: RunConfig{
			Entry: "",
			Args:  []string{},
		},
		Commands: map[string]string{
			"vet":  "go vet ./...",
			"fmt":  "go fmt ./...",
			"test": "go test -v ./...",
		},
	}
}

// DefaultConfigTemplate 返回默认配置文件模板
func DefaultConfigTemplate(projectName, projectMode string) string {
	entry := "."
	if projectMode == "project" {
		entry = "cmd/server"
	}

	return fmt.Sprintf(`# gocar 项目配置文件
# 文档: https://github.com/uselibrary/gocar

# 项目配置
[project]
# 项目模式: "simple" (单文件) 或 "project" (标准目录结构)
# 留空则自动检测
mode = "%s"

# 项目名称，留空则使用目录名
name = "%s"

# 构建配置
[build]
# 构建入口路径 (相对于项目根目录)
# simple 模式默认为 ".", project 模式默认为 "cmd/server"
entry = "%s"

# 输出目录
output = "bin"

# 额外的 ldflags，会追加到默认 ldflags 之后
# 例如: "-X main.version=1.0.0"
ldflags = ""

# 构建标签
# tags = ["jsoniter", "sonic"]

# 额外的环境变量
# extra_env = ["GOPROXY=https://goproxy.cn"]

# 运行配置
[run]
# 运行入口路径，留空则使用 build.entry
entry = ""

# 默认运行参数
# args = ["-config", "config.yaml"]

# 自定义命令
# 格式: 命令名 = "要执行的 shell 命令"
# 使用: gocar <命令名>
# 命令会在项目根目录下执行
[commands]
# 代码检查
vet = "go vet ./..."

# 代码格式化
fmt = "go fmt ./..."

# 运行测试
test = "go test -v ./..."

# lint = "golangci-lint run"
# doc = "godoc -http=:6060"
# proto = "protoc --go_out=. --go-grpc_out=. ./proto/*.proto"
`, projectMode, projectName, entry)
}

// Load 从指定目录加载配置
func Load(projectRoot string) (*GocarConfig, error) {
	configPath := filepath.Join(projectRoot, ConfigFileName)

	// 如果配置文件不存在，返回默认配置
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return DefaultConfig(), nil
	}

	config := DefaultConfig()
	if _, err := toml.DecodeFile(configPath, config); err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", ConfigFileName, err)
	}

	return config, nil
}

// Exists 检查配置文件是否存在
func Exists(projectRoot string) bool {
	configPath := filepath.Join(projectRoot, ConfigFileName)
	_, err := os.Stat(configPath)
	return err == nil
}

// Save 保存配置到文件
func Save(projectRoot, projectName, projectMode string) error {
	configPath := filepath.Join(projectRoot, ConfigFileName)
	content := DefaultConfigTemplate(projectName, projectMode)
	return os.WriteFile(configPath, []byte(content), 0644)
}

// GetBuildEntry 获取构建入口路径
func (c *GocarConfig) GetBuildEntry(defaultMode string) string {
	if c.Build.Entry != "" {
		return c.Build.Entry
	}

	mode := c.Project.Mode
	if mode == "" {
		mode = defaultMode
	}

	if mode == "project" {
		return "cmd/server"
	}
	return "."
}

// GetRunEntry 获取运行入口路径
func (c *GocarConfig) GetRunEntry(defaultMode string) string {
	if c.Run.Entry != "" {
		return c.Run.Entry
	}
	return c.GetBuildEntry(defaultMode)
}

// GetProjectMode 获取项目模式
func (c *GocarConfig) GetProjectMode() string {
	return c.Project.Mode
}

// GetProjectName 获取项目名称
func (c *GocarConfig) GetProjectName(defaultName string) string {
	if c.Project.Name != "" {
		return c.Project.Name
	}
	return defaultName
}

// GetCommand 获取自定义命令
func (c *GocarConfig) GetCommand(name string) (string, bool) {
	cmd, ok := c.Commands[name]
	return cmd, ok
}

// RunCustomCommand 执行自定义命令
func (c *GocarConfig) RunCustomCommand(projectRoot, name string, extraArgs []string) error {
	cmdStr, ok := c.Commands[name]
	if !ok {
		return fmt.Errorf("command '%s' not defined in %s", name, ConfigFileName)
	}

	// 如果有额外参数，追加到命令后面
	if len(extraArgs) > 0 {
		cmdStr = cmdStr + " " + strings.Join(extraArgs, " ")
	}

	// 使用 shell 执行命令
	cmd := exec.Command("sh", "-c", cmdStr)
	cmd.Dir = projectRoot
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

// ListCommands 列出所有自定义命令
func (c *GocarConfig) ListCommands() map[string]string {
	return c.Commands
}
