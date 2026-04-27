package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"gocar/internal/config"
	"gocar/internal/project"
)

// CleanCommand clean 命令
type CleanCommand struct{}

// Run 执行 clean 命令
func (c *CleanCommand) Run(args []string) error {
	projectRoot, appName, _, err := project.DetectProject()
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	cfg, err := config.Load(projectRoot)
	if err != nil {
		return fmt.Errorf("failed to load %s: %w", config.ConfigFileName, err)
	}

	buildOutputDir, err := cfg.ResolveBuildOutputDir(projectRoot)
	if err != nil {
		return err
	}

	// Remove configured build output directory contents
	entries, err := os.ReadDir(buildOutputDir)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Nothing to clean.")
			return nil
		}
		return fmt.Errorf("error reading build output directory: %w", err)
	}

	if len(entries) == 0 {
		fmt.Println("Nothing to clean.")
		return nil
	}

	for _, entry := range entries {
		path := filepath.Join(buildOutputDir, entry.Name())
		if err := os.RemoveAll(path); err != nil {
			fmt.Printf("Error removing %s: %v\n", path, err)
		}
	}

	fmt.Printf("Cleaned build artifacts for '%s' in %s\n", appName, buildOutputDir)
	return nil
}

// Help 返回帮助信息
func (c *CleanCommand) Help() string {
	return `gocar clean - Clean build artifacts

USAGE:
    gocar clean

DESCRIPTION:
	Remove all build artifacts from configured build output directory.
`
}
