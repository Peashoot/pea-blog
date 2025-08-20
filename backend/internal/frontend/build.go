package frontend

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"pea-blog-backend/internal/config"
	"pea-blog-backend/pkg/logger"
)

type BuildService struct {
	config *config.FrontendConfig
	logger *logger.Logger
}

func NewBuildService(cfg *config.FrontendConfig, log *logger.Logger) *BuildService {
	return &BuildService{
		config: cfg,
		logger: log,
	}
}

// CheckAndBuild 检查前端是否需要构建并执行构建
func (b *BuildService) CheckAndBuild() error {
	if !b.config.AutoBuild {
		b.logger.Info("Frontend auto-build is disabled")
		return nil
	}

	// 检查前端源码目录是否存在
	if !b.existsDir(b.config.SourcePath) {
		b.logger.Warn("Frontend source path not found: " + b.config.SourcePath)
		return nil
	}

	// 检查是否需要构建
	needsBuild, err := b.needsBuild()
	if err != nil {
		return fmt.Errorf("failed to check if build is needed: %w", err)
	}

	if !needsBuild {
		b.logger.Info("Frontend build is up to date")
		return nil
	}

	b.logger.Info("Frontend needs rebuild, starting build process...")
	return b.buildFrontend()
}

// needsBuild 检查是否需要重新构建前端
func (b *BuildService) needsBuild() (bool, error) {
	distPath := b.config.DistPath
	sourcePath := b.config.SourcePath

	// 如果dist目录不存在，需要构建
	if !b.existsDir(distPath) {
		b.logger.Info("Frontend dist directory not found, build required")
		return true, nil
	}

	// 检查关键文件是否存在
	indexPath := filepath.Join(distPath, "index.html")
	if !b.existsFile(indexPath) {
		b.logger.Info("Frontend index.html not found, build required")
		return true, nil
	}

	// 获取dist目录的最后修改时间
	distInfo, err := os.Stat(distPath)
	if err != nil {
		return true, err
	}

	// 检查源码文件是否有更新
	srcModTime, err := b.getLatestModTime(sourcePath)
	if err != nil {
		b.logger.Warn("Failed to check source modification time: " + err.Error())
		return false, nil // 如果无法检查，假设不需要构建
	}

	// 如果源码比构建文件新，需要重新构建
	needsBuild := srcModTime.After(distInfo.ModTime())
	if needsBuild {
		b.logger.Info("Source files are newer than build, rebuild required")
	}

	return needsBuild, nil
}

// buildFrontend 执行前端构建
func (b *BuildService) buildFrontend() error {
	b.logger.Info("Starting frontend build...")

	// 解析构建命令
	cmdParts := strings.Fields(b.config.BuildCommand)
	if len(cmdParts) == 0 {
		return fmt.Errorf("empty build command")
	}

	command := cmdParts[0]
	args := cmdParts[1:]

	// 创建命令
	cmd := exec.Command(command, args...)
	cmd.Dir = b.config.SourcePath

	// 设置输出
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 添加环境变量
	cmd.Env = append(os.Environ(), "NODE_ENV=production")

	b.logger.Info(fmt.Sprintf("Executing: %s %s", command, strings.Join(args, " ")))

	// 执行命令
	start := time.Now()
	err := cmd.Run()
	duration := time.Since(start)

	if err != nil {
		b.logger.Error("Frontend build failed", err)
		return fmt.Errorf("frontend build failed: %w", err)
	}

	b.logger.Info(fmt.Sprintf("Frontend build completed successfully in %v", duration))
	return nil
}

// getLatestModTime 获取目录下所有相关文件的最新修改时间
func (b *BuildService) getLatestModTime(dir string) (time.Time, error) {
	var latestTime time.Time

	// 需要检查的文件模式
	patterns := []string{
		filepath.Join(dir, "src", "**", "*.vue"),
		filepath.Join(dir, "src", "**", "*.ts"),
		filepath.Join(dir, "src", "**", "*.js"),
		filepath.Join(dir, "src", "**", "*.css"),
		filepath.Join(dir, "package.json"),
		filepath.Join(dir, "vite.config.ts"),
		filepath.Join(dir, "index.html"),
	}

	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			continue
		}

		for _, match := range matches {
			info, err := os.Stat(match)
			if err != nil {
				continue
			}

			if info.ModTime().After(latestTime) {
				latestTime = info.ModTime()
			}
		}
	}

	// 简化方法：检查src目录的修改时间
	srcDir := filepath.Join(dir, "src")
	if b.existsDir(srcDir) {
		err := filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil // 继续遍历
			}

			// 只检查相关文件类型
			ext := filepath.Ext(path)
			if ext == ".vue" || ext == ".ts" || ext == ".js" || ext == ".css" {
				if info.ModTime().After(latestTime) {
					latestTime = info.ModTime()
				}
			}
			return nil
		})

		if err != nil {
			return latestTime, err
		}
	}

	// 检查package.json等配置文件
	configFiles := []string{
		filepath.Join(dir, "package.json"),
		filepath.Join(dir, "vite.config.ts"),
		filepath.Join(dir, "index.html"),
	}

	for _, file := range configFiles {
		if info, err := os.Stat(file); err == nil {
			if info.ModTime().After(latestTime) {
				latestTime = info.ModTime()
			}
		}
	}

	return latestTime, nil
}

// existsDir 检查目录是否存在
func (b *BuildService) existsDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// existsFile 检查文件是否存在
func (b *BuildService) existsFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

// ForceBuild 强制重新构建前端
func (b *BuildService) ForceBuild() error {
	if !b.config.AutoBuild {
		return fmt.Errorf("frontend auto-build is disabled")
	}
	
	b.logger.Info("Force rebuilding frontend...")
	return b.buildFrontend()
}