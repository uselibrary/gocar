package cli

import (
	"errors"
	"fmt"
	"os"

	"gocar/internal/config"
	"gocar/internal/project"
)

// Version 版本号
var Version = "0.2.0"

// ErrCommandNotFound 命令未找到错误
var ErrCommandNotFound = errors.New("command not found")

// protectedCommands 保护命令列表，这些命令不能被自定义命令覆盖
// new: 创建项目时还没有配置文件
// init: 生成配置文件本身，不能被覆盖
var protectedCommands = map[string]bool{
	"new":  true,
	"init": true,
}

// App CLI 应用
type App struct {
	commands map[string]Command
}

// Command 命令接口
type Command interface {
	Run(args []string) error
	Help() string
}

// NewApp 创建 CLI 应用
func NewApp() *App {
	app := &App{
		commands: make(map[string]Command),
	}

	// 注册命令
	app.commands["new"] = &NewCommand{}
	app.commands["build"] = &BuildCommand{}
	app.commands["run"] = &RunCommand{}
	app.commands["clean"] = &CleanCommand{}
	app.commands["add"] = &AddCommand{}
	app.commands["update"] = &UpdateCommand{}
	app.commands["tidy"] = &TidyCommand{}
	app.commands["init"] = &InitCommand{}

	return app
}

// isProtectedCommand 检查命令是否为保护命令
func isProtectedCommand(cmdName string) bool {
	return protectedCommands[cmdName]
}

// Run 运行应用
func (a *App) Run(args []string) error {
	if len(args) < 2 {
		printHelp()
		return nil
	}

	cmdName := args[1]

	// 处理特殊命令
	switch cmdName {
	case "help", "-h", "--help":
		printHelp()
		return nil
	case "version", "-v", "--version":
		fmt.Printf("gocar %s\n", Version)
		return nil
	}

	// 执行命令
	cmd, ok := a.commands[cmdName]
	if !ok {
		// 尝试执行自定义命令
		if err := a.tryRunCustomCommand(cmdName, args[2:]); err == nil {
			return nil
		}
		fmt.Printf("Unknown command: %s\n", cmdName)
		printHelp()
		return fmt.Errorf("unknown command: %s", cmdName)
	}

	// 对于非保护命令，检查是否有同名的自定义命令覆盖
	if !isProtectedCommand(cmdName) {
		err := a.tryRunCustomCommand(cmdName, args[2:])
		if err == nil {
			return nil
		}
		// 只有当命令不存在时才回退到内置命令
		// 如果命令存在但执行失败，则直接返回错误
		if !errors.Is(err, ErrCommandNotFound) {
			return err
		}
	}

	return cmd.Run(args[2:])
}

// tryRunCustomCommand 尝试执行自定义命令
// 返回 ErrCommandNotFound 表示命令不存在，其他错误表示命令执行失败
func (a *App) tryRunCustomCommand(cmdName string, args []string) error {
	// 检测项目
	projectRoot, _, _, err := project.DetectProject()
	if err != nil {
		return ErrCommandNotFound
	}

	// 加载配置
	cfg, err := config.Load(projectRoot)
	if err != nil {
		return ErrCommandNotFound
	}

	// 检查是否有这个自定义命令
	if _, ok := cfg.GetCommand(cmdName); !ok {
		return ErrCommandNotFound
	}

	// 执行自定义命令
	fmt.Printf("Running custom command: %s\n\n", cmdName)
	if err := cfg.RunCustomCommand(projectRoot, cmdName, args); err != nil {
		fmt.Printf("Command failed: %v\n", err)
		os.Exit(1)
	}

	return nil
}

// printHelp 打印帮助信息
func printHelp() {
	help := `gocar - A cargo-like tool for Go projects

USAGE:
    gocar <COMMAND> [OPTIONS]

COMMANDS:
    new <name> [--mode simple|project]     Create a new Go project
    init                                   Initialize .gocar.toml in current project
    build [--release]                      Build the project
    run [args...]                          Run the project
    clean                                  Clean build artifacts
    add <package>...                       Add dependencies to go.mod
    update [package]...                    Update dependencies
    tidy                                   Tidy up go.mod and go.sum
    help                                   Print this help message
    version                                Print version info

CUSTOM COMMANDS:
    Define custom commands in .gocar.toml [commands] section.
    Custom commands can override built-in commands (except: new, init).
    Example: gocar vet, gocar fmt, gocar test

EXAMPLES:
    gocar new myapp                        Create a simple project
    gocar new myapp --mode project         Create a project-mode project
    gocar init                             Create .gocar.toml config file
    gocar build                            Build in debug mode
    gocar build --release                  Build in release mode
    gocar run                              Build and run the project
    gocar add github.com/gin-gonic/gin     Add a dependency
    gocar vet                              Run custom vet command
`
	fmt.Print(help)
}
