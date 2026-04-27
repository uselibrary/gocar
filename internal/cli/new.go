package cli

import (
	"fmt"
	"strings"

	"gocar/internal/project"
)

// NewCommand new 命令
type NewCommand struct{}

// Run 执行 new 命令
func (c *NewCommand) Run(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("missing project name (usage: gocar new <name> [--mode simple|project])")
	}

	// Check for help
	if args[0] == "help" || args[0] == "--help" || args[0] == "-h" {
		fmt.Print(c.Help())
		return nil
	}

	appName := args[0]

	// Validate project name
	if err := project.ValidateProjectName(appName); err != nil {
		return err
	}

	mode := "simple" // default mode

	// Parse --mode flag
	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "--mode":
			if i+1 < len(args) {
				mode = args[i+1]
				i++ // skip next arg
			} else {
				return fmt.Errorf("--mode requires a value")
			}
		default:
			if strings.HasPrefix(args[i], "-") {
				return fmt.Errorf("unknown option '%s' (run 'gocar new --help' for usage)", args[i])
			}
		}
	}

	// 检查是否是有效模式
	if mode != "simple" && mode != "project" {
		return fmt.Errorf("unknown mode '%s' (available: simple, project)", mode)
	}

	fmt.Printf("Creating new %s project: %s\n", mode, appName)

	creator := project.NewCreator(appName, mode)
	if err := creator.Create(); err != nil {
		return fmt.Errorf("error creating project: %w", err)
	}

	fmt.Printf("\nSuccessfully created project '%s'\n", appName)
	fmt.Printf("\nTo get started:\n")
	fmt.Printf("    cd %s\n", appName)
	fmt.Printf("    gocar build\n")
	fmt.Printf("    gocar run\n")

	return nil
}

// Help 返回帮助信息
func (c *NewCommand) Help() string {
	helpText := `gocar new - Create a new Go project

USAGE:
    gocar new <name> [--mode simple|project]

OPTIONS:
    --mode <mode>    Project mode
                     Available: 'simple' (default), 'project'

EXAMPLES:
    gocar new myapp                   Create a simple project
    gocar new myapp --mode project    Create a project-mode project
`
	return helpText
}
