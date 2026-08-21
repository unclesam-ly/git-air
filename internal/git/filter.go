package git

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// 默认需要忽略的文件名与后缀列表
var defaultIgnoredPatterns = []string{
	// 依赖与锁文件
	"go.sum",
	"package-lock.json",
	"yarn.lock",
	"pnpm-lock.yaml",
	"composer.lock",
	"Cargo.lock",
	"Gemfile.lock",
	// 自动生成的代码
	"*.pb.go",
	"*.gen.go",
	"*_generated.go",
	"*.g.dart",
	// 压缩/静态构建产物
	"*.min.js",
	"*.min.css",
	"*.map",
	"*.bundle.js",
	// 图片与多媒体
	"*.png",
	"*.jpg",
	"*.jpeg",
	"*.gif",
	"*.svg",
	"*.ico",
	"*.pdf",
}

// LoadAirIgnore 从仓库根目录加载 .airignore 文件中的自定义忽略规则
func LoadAirIgnore(repoRoot string) []string {
	if repoRoot == "" {
		return nil
	}
	ignoreFile := filepath.Join(repoRoot, ".airignore")
	file, err := os.Open(ignoreFile)
	if err != nil {
		return nil
	}
	defer file.Close()

	var patterns []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		trimmed := strings.TrimSpace(scanner.Text())
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		patterns = append(patterns, trimmed)
	}
	return patterns
}

// isIgnoredFile 判断文件路径是否命中忽略规则（内置规则 + 用户自定义规则）
func isIgnoredFile(filePath string, customPatterns []string) bool {
	cleanPath := filepath.ToSlash(strings.TrimSpace(filePath))
	baseName := filepath.Base(cleanPath)

	// 1. 先匹配用户自定义规则 (.airignore)
	for _, pattern := range customPatterns {
		if matchPattern(pattern, cleanPath, baseName) {
			return true
		}
	}

	// 2. 再匹配内置默认规则
	for _, pattern := range defaultIgnoredPatterns {
		if matchPattern(pattern, cleanPath, baseName) {
			return true
		}
	}

	return false
}

// matchPattern 单条规则智能匹配（支持目录/、通配符*、精确匹配）
func matchPattern(pattern, cleanPath, baseName string) bool {
	pattern = filepath.ToSlash(strings.TrimSpace(pattern))
	if pattern == "" {
		return false
	}

	// 目录规则：如 "docs/" 或 "vendor/"
	if strings.HasSuffix(pattern, "/") {
		dirPrefix := strings.TrimSuffix(pattern, "/")
		if strings.HasPrefix(cleanPath, dirPrefix+"/") || cleanPath == dirPrefix {
			return true
		}
	}

	// 精确匹配文件名或相对路径
	if pattern == baseName || pattern == cleanPath {
		return true
	}

	// 通配符匹配文件名（如 *.pb.go）
	if matched, _ := filepath.Match(pattern, baseName); matched {
		return true
	}

	// 通配符匹配全路径（如 mock/*.go）
	if matched, _ := filepath.Match(pattern, cleanPath); matched {
		return true
	}

	return false
}

// FilterDiff 过滤掉 diff 中命中忽略规则（如锁文件、生成的代码、.airignore 排除项）的文件区块
func FilterDiff(rawDiff string, customPatterns ...string) string {
	if strings.TrimSpace(rawDiff) == "" {
		return ""
	}

	// 按照 diff --git 切分各个文件的差异区块
	chunks := strings.Split(rawDiff, "diff --git ")
	var filteredChunks []string

	for _, chunk := range chunks {
		if strings.TrimSpace(chunk) == "" {
			continue
		}

		// 还原切分前缀
		fullChunk := "diff --git " + chunk
		lines := strings.Split(chunk, "\n")
		if len(lines) == 0 {
			continue
		}

		// 解析第一行中的文件路径，如：a/internal/ent/schema.go b/internal/ent/schema.go
		firstLine := lines[0]
		parts := strings.Fields(firstLine)
		if len(parts) >= 2 {
			filePath := strings.TrimPrefix(parts[1], "b/")
			if isIgnoredFile(filePath, customPatterns) {
				continue // 命中忽略规则，跳过这个文件的 diff
			}
		}

		filteredChunks = append(filteredChunks, fullChunk)
	}

	return strings.Join(filteredChunks, "")
}
