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
	Profile  ProfilesConfig    `toml:"profile"`
	Commands map[string]string `toml:"commands"`
}

// ProjectConfig 项目配置
type ProjectConfig struct {
	Mode    string `toml:"mode"`    // simple | project
	Name    string `toml:"name"`    // 项目名称，为空时使用目录名
	Version string `toml:"version"` // 项目版本号，构建时自动注入到 main.version
}

// ProfilesConfig 构建配置档案
type ProfilesConfig struct {
	Debug   ProfileConfig `toml:"debug"`
	Release ProfileConfig `toml:"release"`
}

// ProfileConfig 单个构建档案配置
type ProfileConfig struct {
	Ldflags    string `toml:"ldflags"`     // ldflags 参数
	Gcflags    string `toml:"gcflags"`     // 编译器参数
	Trimpath   *bool  `toml:"trimpath"`    // 是否移除路径信息
	CgoEnabled *bool  `toml:"cgo_enabled"` // 是否启用 CGO
	Race       bool   `toml:"race"`        // 是否启用竞态检测
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
	trueVal := true
	falseVal := false
	return &GocarConfig{
		Project: ProjectConfig{
			Mode:    "",
			Name:    "",
			Version: "",
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
		Profile: ProfilesConfig{
			Debug: ProfileConfig{
				Ldflags:    "",
				Gcflags:    "",
				Trimpath:   &falseVal,
				CgoEnabled: nil, // nil 表示跟随系统默认
				Race:       false,
			},
			Release: ProfileConfig{
				Ldflags:    "-s -w",
				Gcflags:    "",
				Trimpath:   &trueVal,
				CgoEnabled: &falseVal,
				Race:       false,
			},
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
		if projectName != "" {
			entry = "cmd/" + projectName
		} else {
			entry = "cmd/server"
		}
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

# 项目版本号
# version = "1.0.0"

# 构建配置
[build]
# 构建入口路径 (相对于项目根目录)
# simple 模式默认为 ".", project 模式默认为 "cmd/<appName>"（即项目名）
entry = "%s"

# 输出目录
output = "bin"

# 额外的 ldflags，会追加到 profile 的 ldflags 之后
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

# Debug 构建配置
# 使用: gocar build (默认)
[profile.debug]
# ldflags = ""              # Debug 默认无 ldflags
# gcflags = "all=-N -l"     # 禁用优化，方便调试
# trimpath = false          # 保留路径信息
# cgo_enabled = true        # 跟随系统默认
# race = false              # 竞态检测 (会显著降低性能)

# Release 构建配置
# 使用: gocar build --release
[profile.release]
ldflags = "-s -w"           # 裁剪符号表和调试信息
# gcflags = ""              # 编译器参数
trimpath = true             # 移除编译路径信息
cgo_enabled = false         # 禁用 CGO 以生成静态二进制
# race = false              # 竞态检测

# 自定义命令
# 格式: 命令名 = "要执行的 shell 命令"
# 使用: gocar <命令名>
# 命令会在项目根目录下执行
#
# 自定义命令可以覆盖以下内置命令: build, run, clean, add, update, tidy
# 保护命令 (new, init) 不可被覆盖
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

# 覆盖内置命令示例 (取消注释以启用):
# build = "make build"
# run = "docker-compose up"
# clean = "make clean && rm -rf dist/"
`, projectMode, projectName, entry)
}

// Load 从指定目录加载配置
func Load(projectRoot string) (*GocarConfig, error) {
	configPath := filepath.Join(projectRoot, ConfigFileName)

	// 使用内置默认配置作为基础
	baseConfig := DefaultConfig()

	// 如果项目配置文件不存在，返回默认配置
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return baseConfig, nil
	}

	// 加载项目配置
	projectConfig := &GocarConfig{}
	if _, err := toml.DecodeFile(configPath, projectConfig); err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", ConfigFileName, err)
	}

	// 合并: 项目配置覆盖默认配置
	finalConfig := mergeProjectConfig(baseConfig, projectConfig)

	return finalConfig, nil
}

// mergeProjectConfig 将项目配置合并到基础配置（项目配置优先）
func mergeProjectConfig(base *GocarConfig, project *GocarConfig) *GocarConfig {
	// Project 配置
	if project.Project.Mode != "" {
		base.Project.Mode = project.Project.Mode
	}
	if project.Project.Name != "" {
		base.Project.Name = project.Project.Name
	}
	if project.Project.Version != "" {
		base.Project.Version = project.Project.Version
	}

	// Build 配置
	if project.Build.Entry != "" {
		base.Build.Entry = project.Build.Entry
	}
	if project.Build.Output != "" {
		base.Build.Output = project.Build.Output
	}
	if project.Build.Ldflags != "" {
		base.Build.Ldflags = project.Build.Ldflags
	}
	if len(project.Build.Tags) > 0 {
		base.Build.Tags = project.Build.Tags
	}
	if len(project.Build.ExtraEnv) > 0 {
		base.Build.ExtraEnv = project.Build.ExtraEnv
	}

	// Run 配置
	if project.Run.Entry != "" {
		base.Run.Entry = project.Run.Entry
	}
	if len(project.Run.Args) > 0 {
		base.Run.Args = project.Run.Args
	}

	// Profile 配置
	// Debug
	if project.Profile.Debug.Ldflags != "" {
		base.Profile.Debug.Ldflags = project.Profile.Debug.Ldflags
	}
	if project.Profile.Debug.Gcflags != "" {
		base.Profile.Debug.Gcflags = project.Profile.Debug.Gcflags
	}
	if project.Profile.Debug.Trimpath != nil {
		base.Profile.Debug.Trimpath = project.Profile.Debug.Trimpath
	}
	if project.Profile.Debug.CgoEnabled != nil {
		base.Profile.Debug.CgoEnabled = project.Profile.Debug.CgoEnabled
	}
	if project.Profile.Debug.Race {
		base.Profile.Debug.Race = project.Profile.Debug.Race
	}

	// Release
	if project.Profile.Release.Ldflags != "" {
		base.Profile.Release.Ldflags = project.Profile.Release.Ldflags
	}
	if project.Profile.Release.Gcflags != "" {
		base.Profile.Release.Gcflags = project.Profile.Release.Gcflags
	}
	if project.Profile.Release.Trimpath != nil {
		base.Profile.Release.Trimpath = project.Profile.Release.Trimpath
	}
	if project.Profile.Release.CgoEnabled != nil {
		base.Profile.Release.CgoEnabled = project.Profile.Release.CgoEnabled
	}
	if project.Profile.Release.Race {
		base.Profile.Release.Race = project.Profile.Release.Race
	}

	// Commands - 项目命令覆盖全局命令
	for name, cmd := range project.Commands {
		base.Commands[name] = cmd
	}

	return base
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
		if c.Project.Name != "" {
			return "cmd/" + c.Project.Name
		}
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

// GetProfile 获取指定模式的构建配置
func (c *GocarConfig) GetProfile(release bool) *ProfileConfig {
	if release {
		return &c.Profile.Release
	}
	return &c.Profile.Debug
}
