package project

import (
	"fmt"
	"os"
	"path/filepath"
)

// Info 项目信息
type Info struct {
	Root string // 项目根目录
	Name string // 项目名称
	Mode string // 项目模式: "simple" or "project"
}

// Detector 项目检测器
type Detector struct{}

// NewDetector 创建项目检测器
func NewDetector() *Detector {
	return &Detector{}
}

// Detect 检测项目信息
func (d *Detector) Detect() (*Info, error) {
	// 查找项目根目录
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current directory: %w", err)
	}

	root := cwd
	for {
		if _, err := os.Stat(filepath.Join(root, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(root)
		if parent == root {
			return nil, fmt.Errorf("not in a Go module (go.mod not found)")
		}
		root = parent
	}

	// 检测项目模式
	mode := d.detectMode(root)
	if mode == "" {
		return nil, fmt.Errorf("cannot detect project mode: no main.go found and cmd/<appName> or cmd/*/main.go don't exist")
	}

	return &Info{
		Root: root,
		Name: filepath.Base(root),
		Mode: mode,
	}, nil
}

// detectMode 检测项目模式
func (d *Detector) detectMode(root string) string {
	// Check for project mode:
	// 1. legacy: cmd/<appName> directory exists
	// 2. any cmd/*/main.go exists (common project layout)
	cmdServerDir := filepath.Join(root, "cmd", "server")
	if stat, err := os.Stat(cmdServerDir); err == nil && stat.IsDir() {
		return "project"
	}

	// check for any cmd/*/main.go
	cmdGlob := filepath.Join(root, "cmd", "*", "main.go")
	matches, err := filepath.Glob(cmdGlob)
	if err == nil && len(matches) > 0 {
		return "project"
	}

	// Simple mode: main.go in root
	if _, err := os.Stat(filepath.Join(root, "main.go")); err == nil {
		return "simple"
	}

	return ""
}

// DetectProject 便捷函数：检测当前项目
func DetectProject() (projectRoot, appName, projectMode string, err error) {
	detector := NewDetector()
	info, err := detector.Detect()
	if err != nil {
		return "", "", "", err
	}
	return info.Root, info.Name, info.Mode, nil
}
