package project

import (
	"fmt"
	"os"
	"path/filepath"

	"gocar/internal/config"
	"gocar/internal/util"
)

// Creator 项目创建器
type Creator struct {
	Name string
	Mode string
}

// NewCreator 创建项目创建器
func NewCreator(name, mode string) *Creator {
	return &Creator{
		Name: name,
		Mode: mode,
	}
}

// Create 创建项目
func (c *Creator) Create() error {
	// Check if directory already exists
	if _, err := os.Stat(c.Name); !os.IsNotExist(err) {
		return fmt.Errorf("directory '%s' already exists", c.Name)
	}

	var err error
	if c.Mode == "simple" {
		err = c.createSimpleProject()
	} else {
		err = c.createProjectMode()
	}

	if err != nil {
		return err
	}

	// Initialize git
	if err := util.InitGit(c.Name); err != nil {
		fmt.Printf("Warning: Failed to initialize git: %v\n", err)
	}

	return nil
}

// createSimpleProject 创建简单项目
func (c *Creator) createSimpleProject() error {
	// Create directories
	dirs := []string{
		c.Name,
		filepath.Join(c.Name, "bin"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Create go.mod
	if err := util.RunCommandSilent(c.Name, "go", "mod", "init", c.Name); err != nil {
		return fmt.Errorf("failed to initialize go.mod: %w", err)
	}

	// Create main.go
	if err := util.WriteFile(filepath.Join(c.Name, "main.go"), SimpleMainTemplate(c.Name)); err != nil {
		return err
	}

	// Create README.md
	if err := util.WriteFile(filepath.Join(c.Name, "README.md"), SimpleReadmeTemplate(c.Name)); err != nil {
		return err
	}

	// Create .gitignore
	if err := util.WriteFile(filepath.Join(c.Name, ".gitignore"), GitignoreTemplate(c.Name)); err != nil {
		return err
	}

	// Create .gocar.toml
	if err := config.Save(c.Name, c.Name, "simple"); err != nil {
		return err
	}

	return nil
}

// createProjectMode 创建项目模式
func (c *Creator) createProjectMode() error {
	// Create directories
	dirs := []string{
		c.Name,
		filepath.Join(c.Name, "cmd", "server"),
		filepath.Join(c.Name, "internal"),
		filepath.Join(c.Name, "pkg"),
		filepath.Join(c.Name, "test"),
		filepath.Join(c.Name, "bin"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Create go.mod
	if err := util.RunCommandSilent(c.Name, "go", "mod", "init", c.Name); err != nil {
		return fmt.Errorf("failed to initialize go.mod: %w", err)
	}

	// Create cmd/server/main.go
	if err := util.WriteFile(filepath.Join(c.Name, "cmd", "server", "main.go"), ProjectMainTemplate(c.Name)); err != nil {
		return err
	}

	// Create .gitkeep files for empty directories
	emptyDirs := []string{
		filepath.Join(c.Name, "internal", ".gitkeep"),
		filepath.Join(c.Name, "pkg", ".gitkeep"),
		filepath.Join(c.Name, "test", ".gitkeep"),
	}
	for _, f := range emptyDirs {
		if err := util.WriteFile(f, ""); err != nil {
			return err
		}
	}

	// Create README.md
	if err := util.WriteFile(filepath.Join(c.Name, "README.md"), ProjectReadmeTemplate(c.Name)); err != nil {
		return err
	}

	// Create .gitignore
	if err := util.WriteFile(filepath.Join(c.Name, ".gitignore"), GitignoreTemplate(c.Name)); err != nil {
		return err
	}

	// Create .gocar.toml
	if err := config.Save(c.Name, c.Name, "project"); err != nil {
		return err
	}

	return nil
}
