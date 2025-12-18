package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"gocar/internal/config"
	"gocar/internal/project"
)

// RunCommand run 命令
type RunCommand struct{}

// Run 执行 run 命令
func (c *RunCommand) Run(args []string) error {
	// Get project info
	projectRoot, appName, projectMode, err := project.DetectProject()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
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

	// Get entry from config
	sourcePath := cfg.GetRunEntry(projectMode)
	if sourcePath != "." && !filepath.IsAbs(sourcePath) && len(sourcePath) > 0 && sourcePath[0] != '.' {
		sourcePath = "./" + sourcePath
	}

	fmt.Printf("Running %s...\n\n", appName)

	runArgs := []string{"run", sourcePath}

	// Add default args from config
	if len(cfg.Run.Args) > 0 {
		runArgs = append(runArgs, cfg.Run.Args...)
	}

	// Add command line args
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

	return nil
}

// Help 返回帮助信息
func (c *RunCommand) Help() string {
	return `gocar run - Run the project

USAGE:
    gocar run [args...]

EXAMPLES:
    gocar run                Run the project
    gocar run --help         Pass --help to the application
`
}
