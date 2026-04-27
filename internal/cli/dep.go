package cli

import (
	"fmt"

	"gocar/internal/project"
	"gocar/internal/util"
)

// AddCommand add 命令
type AddCommand struct{}

// Run 执行 add 命令
func (c *AddCommand) Run(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("missing package name (usage: gocar add <package>...)")
	}

	// Check for help
	if args[0] == "help" || args[0] == "--help" || args[0] == "-h" {
		fmt.Print(c.Help())
		return nil
	}

	// Check if we're in a Go module
	projectRoot, appName, _, err := project.DetectProject()
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	fmt.Printf("Adding dependencies to '%s'...\n", appName)

	// Add each package
	for _, pkg := range args {
		fmt.Printf("  Adding %s...\n", pkg)
		if err := util.RunCommand(projectRoot, "go", "get", pkg); err != nil {
			return fmt.Errorf("error adding %s: %w", pkg, err)
		}
	}

	// Run go mod tidy to clean up
	fmt.Println("Tidying go.mod...")
	if err := util.RunCommand(projectRoot, "go", "mod", "tidy"); err != nil {
		fmt.Printf("Warning: Failed to tidy go.mod: %v\n", err)
	}

	fmt.Println("Successfully added dependencies")
	return nil
}

// Help 返回帮助信息
func (c *AddCommand) Help() string {
	return `gocar add - Add dependencies to the project

USAGE:
    gocar add <package>...

EXAMPLES:
    gocar add github.com/gin-gonic/gin
    gocar add github.com/gin-gonic/gin github.com/spf13/cobra
`
}

// UpdateCommand update 命令
type UpdateCommand struct{}

// Run 执行 update 命令
func (c *UpdateCommand) Run(args []string) error {
	// Check for help
	if len(args) > 0 && (args[0] == "help" || args[0] == "--help" || args[0] == "-h") {
		fmt.Print(c.Help())
		return nil
	}

	// Check if we're in a Go module
	projectRoot, appName, _, err := project.DetectProject()
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	if len(args) == 0 {
		// Update all dependencies
		fmt.Printf("Updating all dependencies for '%s'...\n", appName)
		if err := util.RunCommand(projectRoot, "go", "get", "-u", "./..."); err != nil {
			return fmt.Errorf("error updating dependencies: %w", err)
		}
	} else {
		// Update specific packages
		fmt.Printf("Updating specified dependencies for '%s'...\n", appName)
		for _, pkg := range args {
			fmt.Printf("  Updating %s...\n", pkg)
			if err := util.RunCommand(projectRoot, "go", "get", "-u", pkg); err != nil {
				return fmt.Errorf("error updating %s: %w", pkg, err)
			}
		}
	}

	// Run go mod tidy to clean up
	fmt.Println("Tidying go.mod...")
	if err := util.RunCommand(projectRoot, "go", "mod", "tidy"); err != nil {
		fmt.Printf("Warning: Failed to tidy go.mod: %v\n", err)
	}

	fmt.Println("Successfully updated dependencies")
	return nil
}

// Help 返回帮助信息
func (c *UpdateCommand) Help() string {
	return `gocar update - Update dependencies

USAGE:
    gocar update [package]...

EXAMPLES:
    gocar update                           Update all dependencies
    gocar update github.com/gin-gonic/gin  Update specific package
`
}

// TidyCommand tidy 命令
type TidyCommand struct{}

// Run 执行 tidy 命令
func (c *TidyCommand) Run(args []string) error {
	// Check if we're in a Go module
	projectRoot, appName, _, err := project.DetectProject()
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	fmt.Printf("Tidying go.mod for '%s'...\n", appName)

	if err := util.RunCommand(projectRoot, "go", "mod", "tidy"); err != nil {
		return fmt.Errorf("error tidying go.mod: %w", err)
	}

	fmt.Println("Successfully tidied go.mod")
	return nil
}

// Help 返回帮助信息
func (c *TidyCommand) Help() string {
	return `gocar tidy - Tidy up go.mod and go.sum

USAGE:
    gocar tidy

DESCRIPTION:
    Add missing dependencies and remove unused ones.
`
}
