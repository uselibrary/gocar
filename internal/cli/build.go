package cli

import (
	"fmt"

	"gocar/internal/build"
	"gocar/internal/config"
	"gocar/internal/project"
)

// BuildCommand build 命令
type BuildCommand struct{}

// Run 执行 build 命令
func (c *BuildCommand) Run(args []string) error {
	buildConfig := build.NewConfig()
	target := ""

	// Parse arguments
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "help", "--help", "-h":
			fmt.Print(c.Help())
			return nil
		case "--release":
			buildConfig.Release = true
		case "--with-cgo":
			buildConfig.WithCGO = true
		case "--target":
			if i+1 < len(args) {
				target = args[i+1]
				i++ // skip next arg
			} else {
				return fmt.Errorf("--target requires a value")
			}
		default:
			return fmt.Errorf("unknown option '%s' (run 'gocar build --help' for usage)", arg)
		}
	}

	// Get project info
	projectRoot, appName, projectMode, err := project.DetectProject()
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	// Load config
	cfg, err := config.Load(projectRoot)
	if err != nil {
		fmt.Printf("Warning: %v\n", err)
		cfg = config.DefaultConfig()
	}

	// Apply config overrides
	if cfg.Project.Mode != "" {
		projectMode = cfg.Project.Mode
	}
	appName = cfg.GetProjectName(appName)

	// Parse target if specified
	if target != "" {
		targetOS, targetArch, err := build.ParseTarget(target)
		if err != nil {
			return fmt.Errorf("%v; expected format: <os>/<arch> (example: linux/amd64)", err)
		}
		buildConfig.SetTarget(targetOS, targetArch)
	}

	// Create builder
	builder := build.NewBuilder(projectRoot, appName, projectMode, buildConfig, cfg)

	// Print build info
	builder.PrintBuildInfo()

	// Execute build
	if err := builder.Build(); err != nil {
		return err
	}

	return nil
}

// Help 返回帮助信息
func (c *BuildCommand) Help() string {
	return `gocar build - Build the project

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
    gocar build --release --target linux/arm64   Cross-compile for Linux ARM (release)
    gocar build --with-cgo                       Build with CGO enabled
    gocar build --release --with-cgo             Build in release mode with CGO enabled

COMMON TARGETS:
    linux/amd64     Linux AMD 64-bit
    linux/arm64     Linux ARM 64-bit
    darwin/amd64    macOS Intel
    darwin/arm64    macOS Apple Silicon
    windows/amd64   Windows 64-bit
`
}
