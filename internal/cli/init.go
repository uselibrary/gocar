package cli

import (
	"fmt"
	"os"

	"gocar/internal/config"
	"gocar/internal/project"
)

// InitCommand init 命令
type InitCommand struct{}

// Run 执行 init 命令
func (c *InitCommand) Run(args []string) error {
	// 检查是否请求帮助
	for _, arg := range args {
		if arg == "help" || arg == "--help" || arg == "-h" {
			fmt.Print(c.Help())
			return nil
		}
	}

	// 检测项目
	projectRoot, appName, projectMode, err := project.DetectProject()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		fmt.Println("Please run this command in a Go project directory (where go.mod exists)")
		os.Exit(1)
	}

	// 检查配置文件是否已存在
	if config.Exists(projectRoot) {
		fmt.Printf("%s already exists in this project\n", config.ConfigFileName)
		fmt.Println("Use a text editor to modify it if needed")
		return nil
	}

	// 创建配置文件
	if err := config.Save(projectRoot, appName, projectMode); err != nil {
		fmt.Printf("Error creating %s: %v\n", config.ConfigFileName, err)
		os.Exit(1)
	}

	fmt.Printf("Created %s in %s\n", config.ConfigFileName, projectRoot)
	fmt.Println("\nYou can now customize:")
	fmt.Println("  - [project] section: project mode and name")
	fmt.Println("  - [build] section: build entry path and options")
	fmt.Println("  - [run] section: run entry path and default args")
	fmt.Println("  - [commands] section: custom commands like vet, fmt, test")

	return nil
}

// Help 返回帮助信息
func (c *InitCommand) Help() string {
	return `gocar init - Initialize .gocar.toml configuration file

USAGE:
    gocar init

DESCRIPTION:
    Creates a .gocar.toml configuration file in the current project root.
    The config file allows you to:
    
    - Override the default build entry path (e.g., cmd/myapp instead of cmd/server)
    - Define custom commands (e.g., gocar vet, gocar lint)
    - Set default build tags, ldflags, and environment variables
    - Configure default run arguments

EXAMPLES:
    gocar init                     Create .gocar.toml in current project
`
}
