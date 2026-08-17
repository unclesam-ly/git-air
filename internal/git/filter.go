package git

import (
	"path/filepath"
	"regexp"
	"strings"
)

// 默认需要忽略的文件名与后缀列表
var ignoredPatterns = []string{
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

// isIgnoredFile 判断文件路径是否应该在 Code Review 中被忽略
func isIgnoredFile(filePath string) bool {
	cleanPath := filepath.ToSlash(filePath)
	baseName := filepath.Base(cleanPath)

	for _, pattern := range ignoredPatterns {
		// 精确匹配文件名（如 go.sum）
		if pattern == baseName {
			return true
		}

		// 通配符匹配（如 *.pb.go）
		if matched, _ := filepath.Match(pattern, baseName); matched {
			return true
		}
	}

	return false
}

var diffHeaderRegex = regexp.MustCompile(`(?m)^diff --git a/(.*?) b/(.*?)$`)

// FilterDiff 过滤掉 diff 中无意义的锁文件和自动生成文件的区块
func FilterDiff(rawDiff string) string {
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
			if isIgnoredFile(filePath) {
				continue // 命中忽略规则，跳过这个文件的 diff
			}
		}

		filteredChunks = append(filteredChunks, fullChunk)
	}

	return strings.Join(filteredChunks, "")
}
