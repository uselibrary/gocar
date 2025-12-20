package cli

import (
	"fmt"
	"os"
	"strings"

	"gocar/internal/project"
)

// NewCommand new 命令
type NewCommand struct{}

// Run 执行 new 命令
func (c *NewCommand) Run(args []string) error {
	if len(args) < 1 {
		fmt.Println("Error: Missing project name")
		fmt.Println("Usage: gocar new <name> [--mode simple|project]")
		os.Exit(1)
	}

	// Check for help
	if args[0] == "help" || args[0] == "--help" || args[0] == "-h" {
		fmt.Print(c.Help())
		return nil
	}

	appName := args[0]

	// Validate project name
	if err := project.ValidateProjectName(appName); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
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
				fmt.Println("Error: --mode requires a value")
				os.Exit(1)
			}
		default:
			if strings.HasPrefix(args[i], "-") {
				fmt.Printf("Error: Unknown option '%s'\n", args[i])
				fmt.Println("Run 'gocar new --help' for usage.")
				os.Exit(1)
			}
		}
	}

	// 检查是否是有效模式
	if mode != "simple" && mode != "project" {
		fmt.Printf("Error: Unknown mode '%s'\n", mode)
		fmt.Println("\nAvailable modes: simple, project")
		os.Exit(1)
	}

	fmt.Printf("Creating new %s project: %s\n", mode, appName)

	creator := project.NewCreator(appName, mode)
	if err := creator.Create(); err != nil {
		fmt.Printf("Error creating project: %v\n", err)
		os.Exit(1)
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
