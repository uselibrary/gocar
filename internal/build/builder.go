package build

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"gocar/internal/config"
)

// Builder 构建器
type Builder struct {
	config      *Config
	projectRoot string
	appName     string
	projectMode string
	gocarConfig *config.GocarConfig
}

// NewBuilder 创建构建器
func NewBuilder(projectRoot, appName, projectMode string, config *Config, gocarConfig *config.GocarConfig) *Builder {
	return &Builder{
		config:      config,
		projectRoot: projectRoot,
		appName:     appName,
		projectMode: projectMode,
		gocarConfig: gocarConfig,
	}
}

// Build 执行构建
func (b *Builder) Build() error {
	outputPath := b.GetOutputPath()

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

	fmt.Printf("Build successful: %s\n", b.GetRelativeOutputPath())
	return nil
}

// GetOutputPath 获取完整输出路径
func (b *Builder) GetOutputPath() string {
	outputPath := filepath.Join(b.projectRoot, b.GetRelativeOutputPath())
	return outputPath
}

// GetRelativeOutputPath 获取相对输出路径
func (b *Builder) GetRelativeOutputPath() string {
	targetDir := fmt.Sprintf("%s-%s", b.config.TargetOS, b.config.TargetArch)
	outputDir := filepath.Join("bin", b.config.BuildMode(), targetDir)
	outputPath := filepath.Join(outputDir, b.appName)

	if b.config.TargetOS == "windows" {
		outputPath += ".exe"
	}

	return outputPath
}

// buildCommand 构建 go build 命令
func (b *Builder) buildCommand(outputPath string) *exec.Cmd {
	args := []string{"build"}

	// 构建 ldflags
	ldflags := ""
	if b.config.Release {
		ldflags = "-s -w"
	}
	// 追加配置文件中的 ldflags
	if b.gocarConfig != nil && b.gocarConfig.Build.Ldflags != "" {
		if ldflags != "" {
			ldflags += " " + b.gocarConfig.Build.Ldflags
		} else {
			ldflags = b.gocarConfig.Build.Ldflags
		}
	}
	if ldflags != "" {
		args = append(args, "-ldflags="+ldflags)
	}

	if b.config.Release {
		args = append(args, "-trimpath")
	}

	// 添加构建标签
	if b.gocarConfig != nil && len(b.gocarConfig.Build.Tags) > 0 {
		tags := ""
		for i, tag := range b.gocarConfig.Build.Tags {
			if i > 0 {
				tags += ","
			}
			tags += tag
		}
		args = append(args, "-tags="+tags)
	}

	args = append(args, "-o", outputPath)

	// 从配置获取构建入口
	var entry string
	if b.gocarConfig != nil {
		entry = b.gocarConfig.GetBuildEntry(b.projectMode)
	} else if b.projectMode == "project" {
		entry = "./cmd/server"
	} else {
		entry = "."
	}

	// 确保路径格式正确
	if entry != "." && !filepath.IsAbs(entry) && entry[0] != '.' {
		entry = "./" + entry
	}
	args = append(args, entry)

	cmd := exec.Command("go", args...)
	cmd.Dir = b.projectRoot
	cmd.Env = b.buildEnv()

	return cmd
}

// buildEnv 构建环境变量
func (b *Builder) buildEnv() []string {
	env := os.Environ()

	env = append(env, fmt.Sprintf("GOOS=%s", b.config.TargetOS))
	env = append(env, fmt.Sprintf("GOARCH=%s", b.config.TargetArch))

	if b.config.WithCGO {
		env = append(env, "CGO_ENABLED=1")
	} else if b.config.Release {
		env = append(env, "CGO_ENABLED=0")
	}

	// 添加配置文件中的额外环境变量
	if b.gocarConfig != nil && len(b.gocarConfig.Build.ExtraEnv) > 0 {
		env = append(env, b.gocarConfig.Build.ExtraEnv...)
	}

	return env
}

// PrintBuildInfo 打印构建信息
func (b *Builder) PrintBuildInfo() {
	mode := "debug"
	if b.config.Release {
		mode = "release"
	}

	if !b.config.IsCurrentPlatform() {
		fmt.Printf("Building in %s mode for %s/%s", mode, b.config.TargetOS, b.config.TargetArch)
	} else {
		fmt.Printf("Building in %s mode", mode)
	}

	if b.config.WithCGO {
		fmt.Print(" with CGO enabled")
	}
	fmt.Println("...")
}
